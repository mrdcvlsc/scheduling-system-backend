package geneticalgorithm

import (
	"log"
)

type WeekSchedule [SCHOOL_DAYS_PER_WEEK]DaySchedule

func (sections_weekly_schedules *WeekSchedule) GetDaySchedule(day_idx int) *DaySchedule {
	if day_idx < 0 || day_idx >= SCHOOL_DAYS_PER_WEEK {
		log.Fatalf(
			"GetDaySchedule(day_idx = %d | min:max = 0:%d): error index out of bounds",
			day_idx, SCHOOL_DAYS_PER_WEEK,
		)
	}

	return &(*sections_weekly_schedules)[day_idx]
}
