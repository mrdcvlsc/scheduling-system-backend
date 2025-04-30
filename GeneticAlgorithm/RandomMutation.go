package GeneticAlgorithm

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

const MAX_TIME_SLOT_NUDGE int = 3
const SUBJECT_TIME_SLOT_NUDGE_PROBABILITY int = 90         // %
const SUBJECT_TIME_SLOT_AND_DAY_NUDGE_PROBABILITY int = 50 // %
const SUBJECT_DAY_SWAP_PROBABILITY int = 90                // %
const DAY_SWAP_PERCENT_PROBABILITY int = 5                 // %
const SECTION_WEEK_CLEAR_PERCENT_PROBABILITY int = 2       // %
const SUBJECT_ERASURE_PROBABILITY int = 7                  // %

func ApplyClearDepartmentSchedule(sched Schedule.UniTimeTables, all_curriculums []Curriculum.Curriculum, department_id uint16, selected_semester int) {
	IterateSectionsWeekSchedule(sched, all_curriculums, selected_semester, nil, nil, func(indecies IterIndices, values IterValues) IterReturnType {

		if values.Curriculum.DepartmentID == department_id {
			values.Sched[indecies.Usi] = Schedule.WeekTimeTable{}
		}

		return IterProceed
	})
}

func ApplyRandomDaySwapTimeSlots(
	sched Schedule.UniTimeTables, all_curriculums []Curriculum.Curriculum,
	department_id uint16, selected_semester int,
	rooms []Rooms.Room,
	instructor_id_to_instructor map[uint16]*Instructors.Instructor,
) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	day_swap_attempts := 0
	day_swap_success := 0

	for day := range Const.N_WEEKLY_SCHOOL_DAYS {

		if rng.Int31n(100) >= int32(DAY_SWAP_PERCENT_PROBABILITY) {
			continue
		}

		day_swap := rng.Intn(int(Const.N_WEEKLY_SCHOOL_DAYS))

		if day_swap == day {
			continue
		}

		IterateSectionsWeekSchedule(sched, all_curriculums, selected_semester, nil, nil, func(indecies IterIndices, values IterValues) IterReturnType {
			if values.Curriculum.DepartmentID == department_id {

				mfit := MeasureWeekTimeTableBasicFitness(*values.WeekSched)

				if mfit > 10 {
					return IterProceed
				}

				day_swap_attempts++

				usi := indecies.Usi

				is_instructor_available_day := true
				is_instructor_available_day_swap := true

				for time_slot := range Const.N_DAILY_TIME_SLOTS {
					sched[usi][day][time_slot], sched[usi][day_swap][time_slot] = sched[usi][day_swap][time_slot], sched[usi][day][time_slot]

					instructor_id_a := sched[usi][day][time_slot].GetInstructorID()
					if instructor_id_a != 0 {
						if !instructor_id_to_instructor[instructor_id_a].Time.GetAvailability(day, time_slot) {
							is_instructor_available_day = false
						}
					}

					instructor_id_b := sched[usi][day_swap][time_slot].GetInstructorID()
					if instructor_id_b != 0 {
						if !instructor_id_to_instructor[instructor_id_b].Time.GetAvailability(day_swap, time_slot) {
							is_instructor_available_day_swap = false
						}
					}
				}

				err_day_a := sched.VerticalRangedValidation(rooms, day, 1, 0, Const.N_DAILY_TIME_SLOTS)
				err_day_b := sched.VerticalRangedValidation(rooms, day_swap, 1, 0, Const.N_DAILY_TIME_SLOTS)

				// if there are vertical errors, undo the mutation

				if !((len(err_day_a) == 0) && (len(err_day_b) == 0) && is_instructor_available_day && is_instructor_available_day_swap) {
					for time_slot := range Const.N_DAILY_TIME_SLOTS {
						sched[usi][day][time_slot], sched[usi][day_swap][time_slot] = sched[usi][day_swap][time_slot], sched[usi][day][time_slot]
					}
				} else {
					day_swap_success++
				}
			}

			return IterProceed
		})
	}

	if os.Getenv("LOG_MODE") != "verbose" {
		log.Printf(
			"Random Mutation : [section-swap-days] %d attempts and, %d successful day swaps. (%d/%d)\n",
			day_swap_attempts, day_swap_success,
			day_swap_attempts, day_swap_success,
		)
	}
}

func ApplyRandomSubjectDaySwap(
	sched Schedule.UniTimeTables,
	rooms []Rooms.Room,
	all_curriculums []Curriculum.Curriculum,
	department_id uint16, selected_semester int,
	instructor_id_to_instructor map[uint16]*Instructors.Instructor,
) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	successful_subject_day_swaps := 0
	total_tried_day_swaps := 0
	total_lec_and_lab_subjects := 0

	IterateSectionsWeekSchedule(sched, all_curriculums, selected_semester, nil, nil, func(indicies IterIndices, values IterValues) IterReturnType {

		curriculum := values.Curriculum

		usi := indicies.Usi

		if curriculum.DepartmentID == department_id {

			mfit := MeasureWeekTimeTableBasicFitness(*values.WeekSched)

			if mfit > 10 {
				return IterProceed
			}

			subjects_json := sched[usi].GetWeekSubjectsJSON()

			if len(subjects_json) == 0 {
				return IterProceed
			}

			rng_n := len(subjects_json)

			if rng_n <= 0 {
				return IterProceed
			}

			subject_count_to_try_day_swap := rng.Intn(rng_n) + 1

			total_lec_and_lab_subjects += len(subjects_json)
			total_tried_day_swaps += subject_count_to_try_day_swap

			rng.Shuffle(len(subjects_json), func(i, j int) {
				subjects_json[i], subjects_json[j] = subjects_json[j], subjects_json[i]
			})

			for shuffled_idx := range subject_count_to_try_day_swap {

				if rng.Int31n(100) >= int32(SUBJECT_DAY_SWAP_PROBABILITY) {
					continue
				}

				rand_subject := subjects_json[shuffled_idx]
				day_swap := rng.Intn(Const.N_WEEKLY_SCHOOL_DAYS)

				is_free_time_slot := true
				for i := range rand_subject.TimeSlotSize {
					swap_slot := sched[usi][day_swap].GetTimeSlot(rand_subject.StartingTimeSlot + i)

					is_time_slot_available := swap_slot.GetSubjectID() == 0
					is_instructor_available := instructor_id_to_instructor[rand_subject.InstructorID].Time.GetAvailability(day_swap, rand_subject.StartingTimeSlot+i)

					if !(is_time_slot_available && is_instructor_available) {
						is_free_time_slot = false
						break
					}
				}

				if !is_free_time_slot {
					continue
				}

				for i := range rand_subject.TimeSlotSize {
					old_slot := sched[usi][rand_subject.Day].GetTimeSlot(rand_subject.StartingTimeSlot + i)
					old_slot.Set(0, 0, 0)

					swap_slot := sched[usi][day_swap].GetTimeSlot(rand_subject.StartingTimeSlot + i)
					swap_slot.Set(rand_subject.SubjectID, rand_subject.InstructorID, rand_subject.RoomID)
				}

				err_overlaps := sched.VerticalRangedValidation(rooms,
					day_swap, 1,
					rand_subject.StartingTimeSlot, rand_subject.TimeSlotSize,
				)

				if len(err_overlaps) > 0 {
					for i := range rand_subject.TimeSlotSize {
						old_slot := sched[usi][rand_subject.Day].GetTimeSlot(rand_subject.StartingTimeSlot + i)
						old_slot.Set(rand_subject.SubjectID, rand_subject.InstructorID, rand_subject.RoomID)

						swap_slot := sched[usi][day_swap].GetTimeSlot(rand_subject.StartingTimeSlot + i)
						swap_slot.Set(0, 0, 0)
					}
				} else {
					successful_subject_day_swaps++
				}
			}
		}

		return IterProceed
	})

	if os.Getenv("LOG_MODE") != "verbose" {
		log.Printf("Random Mutation : [subject-day-swaps], from %d lec & lab subjects, there are %d/%d successful day swaps\n", total_lec_and_lab_subjects, successful_subject_day_swaps, total_tried_day_swaps)
	}
}

func ApplyRandomSubjectTimeSlotNudge(
	sched Schedule.UniTimeTables,
	rooms []Rooms.Room,
	all_curriculums []Curriculum.Curriculum,
	department_id uint16, selected_semester int,
	instructor_id_to_instructor map[uint16]*Instructors.Instructor,
) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	successful_subject_time_slot_nudge := 0
	total_tried_time_slot_nudge := 0
	total_lec_and_lab_subjects := 0

	IterateSectionsWeekSchedule(sched, all_curriculums, selected_semester, nil, nil, func(indicies IterIndices, values IterValues) IterReturnType {

		curriculum := values.Curriculum

		usi := indicies.Usi

		if curriculum.DepartmentID == department_id {

			mfit := MeasureWeekTimeTableBasicFitness(*values.WeekSched)

			if mfit > 10 {
				return IterProceed
			}

			subjects_json := sched[usi].GetWeekSubjectsJSON()
			total_lec_and_lab_subjects += len(subjects_json)

			if len(subjects_json) == 0 {
				return IterProceed
			}

			rng_n := len(subjects_json)

			if rng_n <= 0 {
				return IterProceed
			}

			subject_count_to_try_time_slot_nudge := rng.Intn(rng_n) + 1

			rng.Shuffle(len(subjects_json), func(i, j int) {
				subjects_json[i], subjects_json[j] = subjects_json[j], subjects_json[i]
			})

			for shuffled_idx := 0; shuffled_idx < subject_count_to_try_time_slot_nudge; shuffled_idx++ {

				if rng.Int31n(100) >= int32(SUBJECT_TIME_SLOT_NUDGE_PROBABILITY) {
					continue
				}

				total_tried_time_slot_nudge++

				rnd_subject := subjects_json[shuffled_idx]
				nudge_value := Utils.RandomInRange(-MAX_TIME_SLOT_NUDGE, MAX_TIME_SLOT_NUDGE)

				is_nudge_start_idx_lt_min := (rnd_subject.StartingTimeSlot + nudge_value) < 0
				is_nudge_start_idx_gt_max := (rnd_subject.StartingTimeSlot + nudge_value) >= Const.N_DAILY_TIME_SLOTS
				is_nudge_start_idx_valid := !is_nudge_start_idx_lt_min && !is_nudge_start_idx_gt_max

				is_nudge_end_idx_lt_min := (rnd_subject.StartingTimeSlot + nudge_value + rnd_subject.TimeSlotSize - 1) < 0
				is_nudge_end_idx_gt_max := (rnd_subject.StartingTimeSlot + nudge_value + rnd_subject.TimeSlotSize - 1) >= Const.N_DAILY_TIME_SLOTS
				is_nudge_end_idx_valid := !is_nudge_end_idx_lt_min && !is_nudge_end_idx_gt_max

				if !is_nudge_start_idx_valid || !is_nudge_end_idx_valid {
					continue
				}

				is_free_time_slot := true
				for i := 0; i < rnd_subject.TimeSlotSize; i++ {
					nudge_slot := sched[usi][rnd_subject.Day].GetTimeSlot(rnd_subject.StartingTimeSlot + nudge_value + i)

					is_empty_slot := nudge_slot.GetSubjectID() == 0

					is_same_subject_and_room := (nudge_slot.GetSubjectID() == rnd_subject.SubjectID) &&
						(nudge_slot.GetInstructorID() == rnd_subject.InstructorID) &&
						(nudge_slot.GetRoomID() == rnd_subject.RoomID)

					is_instructor_available := instructor_id_to_instructor[rnd_subject.InstructorID].Time.GetAvailability(rnd_subject.Day, rnd_subject.StartingTimeSlot+nudge_value+i)

					if !((is_empty_slot || is_same_subject_and_room) && is_instructor_available) {
						is_free_time_slot = false
						break
					}
				}

				if !is_free_time_slot {
					continue
				}

				for i := 0; i < rnd_subject.TimeSlotSize; i++ {
					old_slot := sched[usi][rnd_subject.Day].GetTimeSlot(rnd_subject.StartingTimeSlot + i)
					old_slot.Set(0, 0, 0)
				}

				for i := 0; i < rnd_subject.TimeSlotSize; i++ {
					nudge_slot := sched[usi][rnd_subject.Day].GetTimeSlot(rnd_subject.StartingTimeSlot + nudge_value + i)
					nudge_slot.Set(rnd_subject.SubjectID, rnd_subject.InstructorID, rnd_subject.RoomID)
				}

				err_overlaps := sched.VerticalRangedValidation(rooms,
					rnd_subject.Day, 1,
					rnd_subject.StartingTimeSlot+nudge_value, rnd_subject.TimeSlotSize,
				)

				if len(err_overlaps) > 0 {
					for i := 0; i < rnd_subject.TimeSlotSize; i++ {
						nudge_slot := sched[usi][rnd_subject.Day].GetTimeSlot(rnd_subject.StartingTimeSlot + nudge_value + i)
						nudge_slot.Set(0, 0, 0)
					}

					for i := 0; i < rnd_subject.TimeSlotSize; i++ {
						old_slot := sched[usi][rnd_subject.Day].GetTimeSlot(rnd_subject.StartingTimeSlot + i)
						old_slot.Set(rnd_subject.SubjectID, rnd_subject.InstructorID, rnd_subject.RoomID)
					}
				} else {
					successful_subject_time_slot_nudge++
				}
			}
		}

		return IterProceed
	})

	if os.Getenv("LOG_MODE") != "verbose" {
		log.Printf(
			"Random Mutation : [time-slot-nudge] from %d lec and lab subjects, there are %d/%d successful subjects nudge on different time slot\n",
			total_lec_and_lab_subjects, successful_subject_time_slot_nudge, total_tried_time_slot_nudge,
		)
	}
}

// TODO: debug there is an error here (currently this mutation function is not being used)
func ApplyRandomSubjectErasure(
	sched Schedule.UniTimeTables,
	all_curriculums []Curriculum.Curriculum,
	department_id uint16, selected_semester int,
) {

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	successful_subject_erased := 0
	total_tried_subject_erased := 0
	total_lec_and_lab_subjects := 0

	IterateSectionsWeekSchedule(sched, all_curriculums, selected_semester, nil, nil, func(indicies IterIndices, values IterValues) IterReturnType {

		curriculum := values.Curriculum

		usi := indicies.Usi

		if curriculum.DepartmentID == department_id {

			mfit := MeasureWeekTimeTableBasicFitness(*values.WeekSched)

			if mfit > 10 {
				return IterProceed
			}

			subjects_json := sched[usi].GetWeekSubjectsJSON()

			if len(subjects_json)/2 == 0 {
				return IterProceed
			}

			rng_n := len(subjects_json)

			if rng_n <= 0 {
				return IterProceed
			}

			subject_count_to_try_erase := rng.Intn(rng_n) + 1

			total_lec_and_lab_subjects += len(subjects_json)

			rng.Shuffle(len(subjects_json), func(i, j int) {
				subjects_json[i], subjects_json[j] = subjects_json[j], subjects_json[i]
			})

			for shuffled_idx := range subject_count_to_try_erase {
				for _, subj_json := range subjects_json {
					if subj_json.SubjectID == subjects_json[shuffled_idx].SubjectID {
						total_tried_subject_erased++

						if rng.Int31n(100) >= int32(SUBJECT_ERASURE_PROBABILITY) {
							continue
						}

						for i := range subj_json.TimeSlotSize {
							clear_slot := sched[usi][subj_json.Day].GetTimeSlot(subj_json.StartingTimeSlot + i)
							clear_slot.Set(0, 0, 0)
						}

						successful_subject_erased++
					}
				}
			}
		}

		return IterProceed
	})

	if os.Getenv("LOG_MODE") != "verbose" {
		log.Printf("Random Mutation : [subject-clear] from %d lec and lab subjects, there are %d/%d successful subjects cleared\n",
			total_lec_and_lab_subjects, successful_subject_erased, total_tried_subject_erased,
		)
	}
}

func ApplyRandomSubjectTimeSlotAndDayNudge(
	sched Schedule.UniTimeTables,
	rooms []Rooms.Room,
	all_curriculums []Curriculum.Curriculum,
	department_id uint16, selected_semester int,
	instructor_id_to_instructor map[uint16]*Instructors.Instructor,
) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	successful_subject_nudge := 0
	total_tried_nudge := 0
	total_lec_and_lab_subjects := 0

	IterateSectionsWeekSchedule(sched, all_curriculums, selected_semester, nil, nil, func(indicies IterIndices, values IterValues) IterReturnType {

		curriculum := values.Curriculum

		usi := indicies.Usi

		if curriculum.DepartmentID == department_id {

			mfit := MeasureWeekTimeTableBasicFitness(*values.WeekSched)

			if mfit > 10 {
				return IterProceed
			}

			subjects_json := sched[usi].GetWeekSubjectsJSON()
			total_lec_and_lab_subjects += len(subjects_json)

			if len(subjects_json) == 0 {
				return IterProceed
			}

			rng_n := len(subjects_json)

			if rng_n <= 0 {
				return IterProceed
			}

			subject_count_to_try_nudge := rng.Intn(rng_n) + 1

			rng.Shuffle(len(subjects_json), func(i, j int) {
				subjects_json[i], subjects_json[j] = subjects_json[j], subjects_json[i]
			})

			for shuffled_idx := 0; shuffled_idx < subject_count_to_try_nudge; shuffled_idx++ {

				if rng.Int31n(100) >= int32(SUBJECT_TIME_SLOT_AND_DAY_NUDGE_PROBABILITY) {
					continue
				}

				total_tried_nudge++

				rnd_subject := subjects_json[shuffled_idx]
				day_swap := rng.Intn(Const.N_WEEKLY_SCHOOL_DAYS)
				nudge_value := Utils.RandomInRange(-MAX_TIME_SLOT_NUDGE, MAX_TIME_SLOT_NUDGE)

				is_nudge_start_idx_lt_min := (rnd_subject.StartingTimeSlot + nudge_value) < 0
				is_nudge_start_idx_gt_max := (rnd_subject.StartingTimeSlot + nudge_value) >= Const.N_DAILY_TIME_SLOTS
				is_nudge_start_idx_valid := !is_nudge_start_idx_lt_min && !is_nudge_start_idx_gt_max

				is_nudge_end_idx_lt_min := (rnd_subject.StartingTimeSlot + nudge_value + rnd_subject.TimeSlotSize - 1) < 0
				is_nudge_end_idx_gt_max := (rnd_subject.StartingTimeSlot + nudge_value + rnd_subject.TimeSlotSize - 1) >= Const.N_DAILY_TIME_SLOTS
				is_nudge_end_idx_valid := !is_nudge_end_idx_lt_min && !is_nudge_end_idx_gt_max

				if !is_nudge_start_idx_valid || !is_nudge_end_idx_valid {
					continue
				}

				is_free_time_slot := true
				for i := 0; i < rnd_subject.TimeSlotSize; i++ {
					nudge_slot := sched[usi][day_swap].GetTimeSlot(rnd_subject.StartingTimeSlot + nudge_value + i)

					is_empty_slot := nudge_slot.GetSubjectID() == 0

					is_same_subject_and_room := (nudge_slot.GetSubjectID() == rnd_subject.SubjectID) &&
						(nudge_slot.GetInstructorID() == rnd_subject.InstructorID) &&
						(nudge_slot.GetRoomID() == rnd_subject.RoomID)

					is_instructor_available := instructor_id_to_instructor[rnd_subject.InstructorID].Time.GetAvailability(day_swap, rnd_subject.StartingTimeSlot+nudge_value+i)

					if !((is_empty_slot || is_same_subject_and_room) && is_instructor_available) {
						is_free_time_slot = false
						break
					}
				}

				if !is_free_time_slot {
					continue
				}

				for i := 0; i < rnd_subject.TimeSlotSize; i++ {
					old_slot := sched[usi][rnd_subject.Day].GetTimeSlot(rnd_subject.StartingTimeSlot + i)
					old_slot.Set(0, 0, 0)
				}

				for i := 0; i < rnd_subject.TimeSlotSize; i++ {
					nudge_slot := sched[usi][day_swap].GetTimeSlot(rnd_subject.StartingTimeSlot + nudge_value + i)
					nudge_slot.Set(rnd_subject.SubjectID, rnd_subject.InstructorID, rnd_subject.RoomID)
				}

				err_overlaps := sched.VerticalRangedValidation(rooms,
					day_swap, 1,
					rnd_subject.StartingTimeSlot+nudge_value, rnd_subject.TimeSlotSize,
				)

				if len(err_overlaps) > 0 {
					for i := 0; i < rnd_subject.TimeSlotSize; i++ {
						nudge_slot := sched[usi][day_swap].GetTimeSlot(rnd_subject.StartingTimeSlot + nudge_value + i)
						nudge_slot.Set(0, 0, 0)
					}

					for i := 0; i < rnd_subject.TimeSlotSize; i++ {
						old_slot := sched[usi][rnd_subject.Day].GetTimeSlot(rnd_subject.StartingTimeSlot + i)
						old_slot.Set(rnd_subject.SubjectID, rnd_subject.InstructorID, rnd_subject.RoomID)
					}
				} else {
					successful_subject_nudge++
				}
			}
		}

		return IterProceed
	})

	if os.Getenv("LOG_MODE") != "verbose" {
		log.Printf(
			"Random Mutation : [day-time-slot-nudge] from %d lec and lab subjects, there are %d/%d successful subjects nudge on different day & time slot\n",
			total_lec_and_lab_subjects, successful_subject_nudge, total_tried_nudge,
		)
	}
}
