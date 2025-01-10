package scheduledsa

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
