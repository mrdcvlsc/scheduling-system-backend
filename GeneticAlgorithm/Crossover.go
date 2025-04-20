package GeneticAlgorithm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
)

func Crossover(
	parent1, parent2 Schedule.UniTimeTables,
	curriculums []Curriculum.Curriculum, selected_semester int,
	dept_id_to_department map[uint16]Departments.Department,
	department_to_encode map[uint16]bool,
	resource_persistence *StorageResources.Persistence,

) (*SchedAndResources, error) {

	if len(parent1) != len(parent2) {
		return nil, fmt.Errorf(
			"parents must have the same length, parent 1 length: %d, parent 2 length: %d",
			len(parent1), len(parent2),
		)
	}

	if os.Getenv("GIN_MODE") != "release" {
		err_p1_repaired_hv := parent1.HorizontalValidation(
			resource_persistence,
			department_to_encode,
			selected_semester,
		)

		if len(err_p1_repaired_hv) > 0 {
			for _, err := range err_p1_repaired_hv {
				log.Printf("Horizontal validation error P1: %s", err.Error())
			}

			panic("P1 : opps there is a horizontal validation error, which should not happen!")
		}

		err_p2_repaired_hv := parent2.HorizontalValidation(
			resource_persistence,
			department_to_encode,
			selected_semester,
		)

		if len(err_p2_repaired_hv) > 0 {
			for _, err := range err_p2_repaired_hv {
				log.Printf("Horizontal validation error P2: %s", err.Error())
			}

			panic("P2 : opps there is a horizontal validation error, which should not happen!")
		}
	}

	total_encoding_tries := 0
	successful_base_parent_encoded := 0
	successful_fallback_parent_encoded := 0
	failed_parents_encoding := 0

	offspring := make(Schedule.UniTimeTables, len(parent1))

	is_err_to_return := false
	var return_err error

	IterateSectionsWeekSchedule(nil, curriculums, selected_semester, nil, nil, func(indicies IterIndices, values IterValues) IterReturnType {
		parent_1_subjects := parent1[indicies.Usi].GetWeekSubjectsJSON()
		parent_2_subjects := parent2[indicies.Usi].GetWeekSubjectsJSON()

		if len(parent_1_subjects) != len(parent_2_subjects) {
			is_err_to_return = true
			return_err = fmt.Errorf(
				"parent subjects must have the same length, parent 1 subjects: %d, parent 2 subjects: %d",
				len(parent_1_subjects), len(parent_2_subjects),
			)

			return IterBreakCurriculumLoop
		}

		sort.Slice(parent_1_subjects, func(i, j int) bool {
			if parent_1_subjects[i].SubjectID == parent_1_subjects[j].SubjectID {
				return parent_1_subjects[i].TimeSlotSize < parent_1_subjects[j].TimeSlotSize
			}

			return parent_1_subjects[i].SubjectID < parent_1_subjects[j].SubjectID
		})

		sort.Slice(parent_2_subjects, func(i, j int) bool {
			if parent_2_subjects[i].SubjectID == parent_2_subjects[j].SubjectID {
				return parent_2_subjects[i].TimeSlotSize < parent_2_subjects[j].TimeSlotSize
			}

			return parent_2_subjects[i].SubjectID < parent_2_subjects[j].SubjectID
		})

		// sanity check
		for i := 0; i < max(len(parent_1_subjects), len(parent_2_subjects)); i++ {
			is_equal_subject_id := parent_1_subjects[i].SubjectID == parent_2_subjects[i].SubjectID
			is_equal_time_slot_size := parent_1_subjects[i].TimeSlotSize == parent_2_subjects[i].TimeSlotSize

			if !is_equal_subject_id {
				is_err_to_return = true
				return_err = errors.New("crossover unexpected error, parents have contain different subjects")

				return IterBreakCurriculumLoop
			}

			if !is_equal_time_slot_size {
				is_err_to_return = true
				return_err = errors.New("crossover unexpected error, parent subjects have different time slot sizes")

				return IterBreakCurriculumLoop
			}
		}

		for i := 0; i < len(parent_1_subjects); i++ {

			is_equal_subject_id := (parent_1_subjects[i].SubjectID == parent_2_subjects[i].SubjectID)
			is_equal_subject_time_slot_size := (parent_1_subjects[i].TimeSlotSize == parent_2_subjects[i].TimeSlotSize)

			if !is_equal_subject_id {
				is_err_to_return = true
				return_err = fmt.Errorf(
					"parents must have the same subject IDs, parent 1 subject ID: %d, parent 2 subject ID: %d",
					parent_1_subjects[i].SubjectID, parent_2_subjects[i].SubjectID,
				)

				return IterBreakCurriculumLoop
			}

			if !is_equal_subject_time_slot_size {
				is_err_to_return = true
				return_err = fmt.Errorf(
					"parents must have the same subject time slot size, parent 1 subject time slot size: %d, parent 2 subject time slot size: %d",
					parent_1_subjects[i].TimeSlotSize, parent_2_subjects[i].TimeSlotSize,
				)

				return IterBreakCurriculumLoop
			}

			total_encoding_tries++

			//////////////////////////////////////////////////////////////////////////////////////////////
			//                             ENCODE THE BASE PARENT SUBJECTS
			//////////////////////////////////////////////////////////////////////////////////////////////

			var base_parent_subjects []Schedule.TimeSlotSubjectJSON

			if i%2 == 0 {
				base_parent_subjects = parent_1_subjects
			} else {
				base_parent_subjects = parent_2_subjects
			}

			base_parent_result := inherit_trait_from_a_parent(
				i, indicies.Usi,
				offspring,
				base_parent_subjects,
				resource_persistence,
			)

			if base_parent_result.success {
				successful_base_parent_encoded++

				if base_parent_result.has_extended_subject {
					i++
				}

				continue // to next subject
			}

			//////////////////////////////////////////////////////////////////////////////////////////////
			//                            ENCODE THE FALLBACK PARENT SUBJECTS
			//////////////////////////////////////////////////////////////////////////////////////////////

			// if the base parent failed to encode, we will try to encode the other parent

			var fallback_parent_subjects []Schedule.TimeSlotSubjectJSON

			if i%2 == 0 {
				fallback_parent_subjects = parent_2_subjects
			} else {
				fallback_parent_subjects = parent_1_subjects
			}

			fallback_parent_result := inherit_trait_from_a_parent(
				i, indicies.Usi,
				offspring,
				fallback_parent_subjects,
				resource_persistence,
			)

			if fallback_parent_result.has_extended_subject {
				i++
			}

			if fallback_parent_result.success {
				successful_fallback_parent_encoded++
				continue // to next subject
			}

			// if both parents failed to encode, we will just try to re-encode to repair it later
			failed_parents_encoding++
		}

		return IterProceed
	})

	log.Printf(
		"Crossover: total encoding tries: %d, successful base parent encoded: %d, successful fallback parent encoded: %d, failed parents encoding: %d",
		total_encoding_tries, successful_base_parent_encoded, successful_fallback_parent_encoded, failed_parents_encoding,
	)

	if is_err_to_return {
		return nil, return_err
	}

	encoding_resource, err_gen_encoding_resource := GenerateEncodingResourceFromUniTimeTable(
		offspring, curriculums, selected_semester, resource_persistence,
	)

	if err_gen_encoding_resource != nil {
		return nil, err_gen_encoding_resource
	}

	if failed_parents_encoding > 0 {
		// re-encode the schedule - fillup missing time slots

		repaired_sched, repaired_encoding_resource, err_repair_encoding := EncodeIndividualGenome(
			offspring, curriculums,
			dept_id_to_department, encoding_resource,
			department_to_encode, selected_semester, 0,
		)

		if err_repair_encoding != nil {
			return nil, err_repair_encoding
		}

		if os.Getenv("GIN_MODE") != "release" {
			err_repaired_vv := repaired_sched.VerticalValidation(resource_persistence)

			if len(err_repaired_vv) > 0 {
				panic("why the fudge there is a vertical error here?")
			}

			err_repaired_hv := repaired_sched.HorizontalValidation(
				resource_persistence,
				department_to_encode,
				selected_semester,
			)

			if len(err_repaired_hv) > 0 {
				for _, err := range err_repaired_hv {
					log.Printf("Horizontal validation error: %s", err.Error())
				}

				panic("END: opps there is a horizontal validation, which should not happen!")
			}
		}

		return &SchedAndResources{
			UniSched:  repaired_sched,
			Resources: repaired_encoding_resource,
		}, nil
	}

	return &SchedAndResources{
		UniSched:  offspring,
		Resources: encoding_resource,
	}, nil
}

type inherit_trait_result struct {
	success              bool
	has_extended_subject bool
}

func inherit_trait_from_a_parent(
	i, usi int,
	offspring Schedule.UniTimeTables,
	json_subjects []Schedule.TimeSlotSubjectJSON,
	resource_persistence *StorageResources.Persistence,
) inherit_trait_result {
	subject := &json_subjects[i]

	has_extended_subject := false

	if i+1 < len(json_subjects) {
		if json_subjects[i+1].SubjectID == subject.SubjectID {
			has_extended_subject = true
		}
	}

	is_first_target_time_slot_free := true
	is_second_target_time_slot_free := true

	for j := 0; j < subject.TimeSlotSize; j++ {
		if offspring[usi][subject.Day][subject.StartingTimeSlot+j].GetSubjectID() != 0 {
			is_first_target_time_slot_free = false
			break
		}
	}

	if is_first_target_time_slot_free {
		for j := 0; j < subject.TimeSlotSize; j++ {
			offspring[usi][subject.Day][subject.StartingTimeSlot+j].SetSubjectID(subject.SubjectID)
			offspring[usi][subject.Day][subject.StartingTimeSlot+j].SetInstructorID(subject.InstructorID)
			offspring[usi][subject.Day][subject.StartingTimeSlot+j].SetRoomID(subject.RoomID)
		}
	}

	if has_extended_subject {
		subj_extend := &json_subjects[i+1]

		for j := 0; j < subj_extend.TimeSlotSize; j++ {
			if offspring[usi][subj_extend.Day][subj_extend.StartingTimeSlot+j].GetSubjectID() != 0 {
				is_second_target_time_slot_free = false
				break
			}
		}

		if is_second_target_time_slot_free {
			for j := 0; j < subj_extend.TimeSlotSize; j++ {
				offspring[usi][subj_extend.Day][subj_extend.StartingTimeSlot+j].SetSubjectID(subj_extend.SubjectID)
				offspring[usi][subj_extend.Day][subj_extend.StartingTimeSlot+j].SetInstructorID(subj_extend.InstructorID)
				offspring[usi][subj_extend.Day][subj_extend.StartingTimeSlot+j].SetRoomID(subj_extend.RoomID)
			}
		}
	}

	var err_vv []error
	var err_vv_extend []error

	if is_first_target_time_slot_free {
		err_vv = offspring.VerticalRangedValidation(
			resource_persistence,
			subject.Day, 1,
			subject.StartingTimeSlot, subject.TimeSlotSize,
		)
	}

	if has_extended_subject && is_second_target_time_slot_free {
		subj_extend := &json_subjects[i+1]

		err_vv_extend = offspring.VerticalRangedValidation(
			resource_persistence,
			subj_extend.Day, 1,
			subj_extend.StartingTimeSlot, subj_extend.TimeSlotSize,
		)
	}

	if has_extended_subject {
		if len(err_vv) == 0 && len(err_vv_extend) == 0 && is_first_target_time_slot_free && is_second_target_time_slot_free {
			return inherit_trait_result{
				success:              true,
				has_extended_subject: true,
			}
		}
	} else {
		if len(err_vv) == 0 && is_first_target_time_slot_free {
			return inherit_trait_result{
				success:              true,
				has_extended_subject: false,
			}
		}
	}

	if is_first_target_time_slot_free {
		for j := 0; j < subject.TimeSlotSize; j++ {
			offspring[usi][subject.Day][subject.StartingTimeSlot+j].SetSubjectID(0)
			offspring[usi][subject.Day][subject.StartingTimeSlot+j].SetInstructorID(0)
			offspring[usi][subject.Day][subject.StartingTimeSlot+j].SetRoomID(0)
		}
	}

	if has_extended_subject && is_second_target_time_slot_free {
		subj_extend := &json_subjects[i+1]

		for j := 0; j < subj_extend.TimeSlotSize; j++ {
			offspring[usi][subj_extend.Day][subj_extend.StartingTimeSlot+j].SetSubjectID(0)
			offspring[usi][subj_extend.Day][subj_extend.StartingTimeSlot+j].SetInstructorID(0)
			offspring[usi][subj_extend.Day][subj_extend.StartingTimeSlot+j].SetRoomID(0)
		}
	}

	return inherit_trait_result{
		success:              false,
		has_extended_subject: has_extended_subject,
	}
}
