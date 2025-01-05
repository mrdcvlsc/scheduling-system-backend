package geneticalgorithm

import (
	"fmt"
)

type WeekSchedule [N_WEEKLY_SCHOOL_DAYS]DaySchedule

func (sections_weekly_schedules *WeekSchedule) GetDaySchedule(day_idx int) *DaySchedule {
	if day_idx < 0 || day_idx >= N_WEEKLY_SCHOOL_DAYS {
		panic(fmt.Sprintf(
			"GetDaySchedule(day_idx = %d | min:max = 0:%d): error index out of bounds",
			day_idx, N_WEEKLY_SCHOOL_DAYS,
		))
	}

	return &(*sections_weekly_schedules)[day_idx]
}
