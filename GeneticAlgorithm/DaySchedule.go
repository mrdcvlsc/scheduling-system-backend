package geneticalgorithm

import (
	"fmt"
)

type DaySchedule [DAILY_TIME_SLOTS]TimeSlot

func (day_schedule *DaySchedule) GetTimeSlot(time_slot_idx int) *TimeSlot {
	if time_slot_idx < 0 || time_slot_idx >= DAILY_TIME_SLOTS {
		panic(fmt.Sprintf(
			"GetTimeSlot(time_slot_idx = %d | min:max = 0:%d): error index out of bounds",
			time_slot_idx, DAILY_TIME_SLOTS,
		))
	}

	return &(*day_schedule)[time_slot_idx]
}
