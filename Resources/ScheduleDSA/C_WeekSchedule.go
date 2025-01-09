package scheduledsa

import (
	"fmt"
)

// type for weekly schedule of a class / section.
//
// this is just an array of `ScheduleDay` types.
type ScheduleWeek [N_WEEKLY_SCHOOL_DAYS]ScheduleDay

func (sections_weekly_schedules *ScheduleWeek) GetDaySchedule(day_idx int) *ScheduleDay {
	if day_idx < 0 || day_idx >= N_WEEKLY_SCHOOL_DAYS {
		panic(fmt.Sprintf(
			"GetDaySchedule(day_idx = %d | min:max = 0:%d): error index out of bounds",
			day_idx, N_WEEKLY_SCHOOL_DAYS,
		))
	}

	return &(*sections_weekly_schedules)[day_idx]
}
