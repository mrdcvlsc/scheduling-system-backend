package GeneticAlgorithm

import (
	"log"
	"math/rand"
	"time"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

const MAX_TIME_SLOT_NUDGE int = 4
const SUBJECT_TIME_SLOT_NUDGE_PROBABILITY int = 20   // %
const SUBJECT_DAY_SWAP_PROBABILITY int = 20          // %
const DAY_SWAP_PERCENT_PROBABILITY int = 1           // %
const SECTION_WEEK_CLEAR_PERCENT_PROBABILITY int = 2 // %
const SUBJECT_ERASURE_PROBABILITY int = 3            // %

func ApplyRandomDaySwapTimeSlots(sched Schedule.UniTimeTables) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	for usi := range len(sched) {
		for day := range Const.N_WEEKLY_SCHOOL_DAYS {

			if rng.Int31n(100) >= int32(DAY_SWAP_PERCENT_PROBABILITY) {
				continue
			}

			log.Print("random day swap mutation")

			day_swap := rng.Intn(int(Const.N_WEEKLY_SCHOOL_DAYS))

			if day_swap == day {
				continue
			}

			for time_slot := range Const.N_DAILY_TIME_SLOTS {
				sched[usi][day][time_slot], sched[usi][day_swap][time_slot] = sched[usi][day_swap][time_slot], sched[usi][day][time_slot]
			}
		}
	}
}

func ApplyRandomSectionWeekClear(sched Schedule.UniTimeTables) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	for usi := range len(sched) {
		if rng.Int31n(100) >= int32(SECTION_WEEK_CLEAR_PERCENT_PROBABILITY) {
			continue
		}

		sched[usi] = Schedule.WeekTimeTable{}
	}
}

func ApplyRandomSubjectDaySwap(sched Schedule.UniTimeTables, resource_persistence *StorageResources.Persistence) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	successful_subject_day_swaps := 0
	total_tried_day_swaps := 0
	total_lec_and_lab_subjects := 0

	for usi := range len(sched) {

		subjects_json := sched[usi].GetWeekSubjectsJSON()
		subject_count_to_try_day_swap := rng.Intn(len(subjects_json)) + 1

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

				if swap_slot.GetSubjectID() != 0 {
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

			err_overlaps := sched.VerticalRangedValidation(resource_persistence,
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

	log.Printf("Random Mutation : from %d lec and lab subjects, there are %d/%d successful day swaps\n", total_lec_and_lab_subjects, successful_subject_day_swaps, total_tried_day_swaps)
}

func ApplyRandomSubjectTimeSlotNudge(sched Schedule.UniTimeTables, resource_persistence *StorageResources.Persistence) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	successful_subject_time_slot_nudge := 0
	total_tried_time_slot_nudge := 0
	total_lec_and_lab_subjects := 0

	for usi := range len(sched) {

		subjects_json := sched[usi].GetWeekSubjectsJSON()
		subject_count_to_try_time_slot_nudge := rng.Intn(len(subjects_json)) + 1

		total_lec_and_lab_subjects += len(subjects_json)
		total_tried_time_slot_nudge += subject_count_to_try_time_slot_nudge

		rng.Shuffle(len(subjects_json), func(i, j int) {
			subjects_json[i], subjects_json[j] = subjects_json[j], subjects_json[i]
		})

		for shuffled_idx := range subject_count_to_try_time_slot_nudge {

			if rng.Int31n(100) >= int32(SUBJECT_TIME_SLOT_NUDGE_PROBABILITY) {
				continue
			}

			rand_subject := subjects_json[shuffled_idx]
			nudge_value := Utils.RandomInRange(-MAX_TIME_SLOT_NUDGE, MAX_TIME_SLOT_NUDGE)

			if rand_subject.StartingTimeSlot+nudge_value < 0 || (rand_subject.StartingTimeSlot+nudge_value+rand_subject.TimeSlotSize) > Const.N_DAILY_TIME_SLOTS {
				continue
			}

			is_free_time_slot := true
			for i := range rand_subject.TimeSlotSize {
				nudge_slot := sched[usi][rand_subject.Day].GetTimeSlot(rand_subject.StartingTimeSlot + nudge_value + i)

				if !(nudge_slot.GetSubjectID() == 0 || (nudge_slot.GetSubjectID() == rand_subject.SubjectID && nudge_slot.GetRoomID() == rand_subject.RoomID)) {
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

				nudge_slot := sched[usi][rand_subject.Day].GetTimeSlot(rand_subject.StartingTimeSlot + nudge_value + i)
				nudge_slot.Set(rand_subject.SubjectID, rand_subject.InstructorID, rand_subject.RoomID)
			}

			err_overlaps := sched.VerticalRangedValidation(resource_persistence,
				rand_subject.Day, 1,
				rand_subject.StartingTimeSlot+nudge_value, rand_subject.TimeSlotSize,
			)

			if len(err_overlaps) > 0 {
				for i := range rand_subject.TimeSlotSize {
					nudge_slot := sched[usi][rand_subject.Day].GetTimeSlot(rand_subject.StartingTimeSlot + nudge_value + i)
					nudge_slot.Set(0, 0, 0)

					old_slot := sched[usi][rand_subject.Day].GetTimeSlot(rand_subject.StartingTimeSlot + i)
					old_slot.Set(rand_subject.SubjectID, rand_subject.InstructorID, rand_subject.RoomID)
				}
			} else {
				successful_subject_time_slot_nudge++
			}
		}

	}

	log.Printf("Random Mutation : from %d lec and lab subjects, there are %d/%d successful subjects nudge on different time slot\n", total_lec_and_lab_subjects, successful_subject_time_slot_nudge, total_tried_time_slot_nudge)
}

func ApplyRandomSubjectErasure(sched Schedule.UniTimeTables, resource_persistence *StorageResources.Persistence) {

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	successful_subject_erased := 0
	total_tried_subject_erased := 0
	total_lec_and_lab_subjects := 0

	for usi := range len(sched) {

		subjects_json := sched[usi].GetWeekSubjectsJSON()
		subject_count_to_try_erase := rng.Intn(len(subjects_json)/2) + 1 // 50% of the subjects to try to erase

		total_lec_and_lab_subjects += len(subjects_json)
		total_tried_subject_erased += subject_count_to_try_erase

		rng.Shuffle(len(subjects_json), func(i, j int) {
			subjects_json[i], subjects_json[j] = subjects_json[j], subjects_json[i]
		})

		for shuffled_idx := range subject_count_to_try_erase {

			if rng.Int31n(100) >= int32(SUBJECT_ERASURE_PROBABILITY) {
				continue
			}

			rand_subject := subjects_json[shuffled_idx]

			for i := range rand_subject.TimeSlotSize {
				clear_slot := sched[usi][rand_subject.Day].GetTimeSlot(rand_subject.StartingTimeSlot + i)
				clear_slot.Set(0, 0, 0)
			}
		}

	}

	log.Printf("Random Mutation : from %d lec and lab subjects, there are %d/%d successful subjects cleared\n", total_lec_and_lab_subjects, successful_subject_erased, total_tried_subject_erased)
}
