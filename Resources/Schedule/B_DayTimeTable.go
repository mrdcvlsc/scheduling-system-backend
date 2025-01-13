package Schedule

import (
	"fmt"
)

// represent a schedule of a class / section in a single day.
//
// this is just an array of `TimeSlot` types.
type DayTimeTable [N_DAILY_TIME_SLOTS]TimeSlot

func (day *DayTimeTable) Get(time_slot_idx int) *TimeSlot {
	if time_slot_idx < 0 || time_slot_idx >= N_DAILY_TIME_SLOTS {
		panic(fmt.Sprintf(
			"GetTimeSlot(time_slot_idx = %d | min:max = 0:%d): error index out of bounds",
			time_slot_idx, N_DAILY_TIME_SLOTS,
		))
	}

	return &(*day)[time_slot_idx]
}

func (day *DayTimeTable) Availability(time_slot_idx, hours int) bool {
	if time_slot_idx < 0 || time_slot_idx >= N_DAILY_TIME_SLOTS {
		panic(fmt.Sprintf(
			"GetTimeSlot(time_slot_idx = %d | min:max = 0:%d): error index out of bounds",
			time_slot_idx, N_DAILY_TIME_SLOTS,
		))
	}

	if (time_slot_idx + hours - 1) >= N_DAILY_TIME_SLOTS {
		return false
	}

	var counted_time_slots int
	for counted_time_slots = 0; counted_time_slots < hours; counted_time_slots++ {
		if day[time_slot_idx+counted_time_slots].GetSubjectID() != 0 {
			return false
		}
	}

	// sanity check, just return true later if there is no problem
	return counted_time_slots == hours
}
