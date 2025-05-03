package GeneticAlgorithm

import (
	"fmt"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
validate assigned subjects to every section schedules in the whole university.

to validate whole university schedules, set department to encode to nil:

	department_to_validate = nil
*/
func HorizontalValidation(
	university_sched Schedule.UniTimeTables,
	curriculums []Curriculum.Curriculum,
	department_to_validate map[uint16]bool, selected_semester int,
) []error {

	errs_slice := make([]error, 0, 16)

	/////////////////////////////////////////////////////////////////////////////////
	//                            HORIZONTAL CHECKS
	/////////////////////////////////////////////////////////////////////////////////

	total_university_sections := Curriculum.GetTotalNumberOfSections(curriculums, selected_semester)

	if total_university_sections != len(university_sched) {
		errs_slice = append(errs_slice, fmt.Errorf(
			"read total university sections (%d) in persistence did not match the university schedule instance (%d)",
			total_university_sections, len(university_sched),
		))
	}

	IterateSectionsWeekSchedule(university_sched, curriculums, selected_semester, nil, nil, func(indicies IterIndices, values IterValues) IterReturnType {

		curriculum := values.Curriculum
		year_level := values.YearLevel
		semester := values.Semester

		section_idx := indicies.Section
		usi := indicies.Usi

		if department_to_validate != nil {
			is_to_validate := department_to_validate[curriculum.DepartmentID]

			if !is_to_validate {
				return IterProceed // skip section that does not need horizontal validation
			}
		}

		is_subject_id_to_is_in_curriculum := make(map[uint16]bool)
		stray_subject_ids := make([]uint16, 0)

		for _, section_subject := range values.Semester.Subjects {
			is_subject_id_to_is_in_curriculum[section_subject.ID] = true
		}

		subject_id_to_time_slot_count := make(map[uint16]int)

		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				subject_id := university_sched[usi][day][time_slot].GetSubjectID()

				if subject_id == 0 {
					continue
				}

				if !is_subject_id_to_is_in_curriculum[subject_id] {
					stray_subject_ids = append(stray_subject_ids, subject_id)
				}

				_, has_subjsubject_id := subject_id_to_time_slot_count[subject_id]

				if !has_subjsubject_id {
					subject_id_to_time_slot_count[subject_id] = 1
				} else {
					subject_id_to_time_slot_count[subject_id]++
				}
			}
		}

		if len(semester.Subjects) != len(subject_id_to_time_slot_count) {
			errs_slice = append(errs_slice, fmt.Errorf(
				"detected %d missing subject(s) in %s, %s, %s, section %s (usi:%d)",
				len(semester.Subjects)-len(subject_id_to_time_slot_count),
				curriculum.CurriculumCode,
				semester.Name,
				year_level.Name,
				Curriculum.SECTION[section_idx],
				usi,
			))
		}

		for _, subject := range semester.Subjects {
			_, has_subject_id := subject_id_to_time_slot_count[subject.ID]

			if !has_subject_id {
				errs_slice = append(errs_slice, fmt.Errorf(
					"the subject %s was not assigned to %s, %s, %s, section %s (usi:%d)",
					subject.Code, curriculum.CurriculumCode, year_level.Name, semester.Name, Curriculum.SECTION[section_idx], usi,
				))
			} else if ((subject.LecHours + subject.LabHours) * Const.N_HOUR_TIME_SLOTS) > uint8(subject_id_to_time_slot_count[subject.ID]) {
				errs_slice = append(errs_slice, fmt.Errorf(
					"the subject %s has missing time slot allocations, expecting %d, but only found %d in %s, %s, %s, section %s (usi:%d)",
					subject.Code,
					((subject.LecHours+subject.LabHours)*Const.N_HOUR_TIME_SLOTS), uint8(subject_id_to_time_slot_count[subject.ID]),
					curriculum.CurriculumCode, year_level.Name, semester.Name, Curriculum.SECTION[section_idx], usi,
				))
			} else if ((subject.LecHours + subject.LabHours) * Const.N_HOUR_TIME_SLOTS) < uint8(subject_id_to_time_slot_count[subject.ID]) {
				errs_slice = append(errs_slice, fmt.Errorf(
					"the subject %s has extra time slot allocations, expecting only %d, but found %d in %s, %s, %s, section %s (usi:%d)",
					subject.Code,
					((subject.LecHours+subject.LabHours)*Const.N_HOUR_TIME_SLOTS), uint8(subject_id_to_time_slot_count[subject.ID]),
					curriculum.CurriculumCode, year_level.Name, semester.Name, Curriculum.SECTION[section_idx], usi,
				))
			}
		}

		if len(stray_subject_ids) > 0 {
			errs_slice = append(errs_slice, fmt.Errorf(
				"the schedule in %s, %s, %s, section %s (usi:%d) has stray subject ids: [%v]",
				curriculum.CurriculumCode,
				year_level.Name,
				semester.Name,
				Curriculum.SECTION[section_idx],
				usi, stray_subject_ids,
			))
		}

		return IterProceed
	})

	return errs_slice
}
