package scheduledsa

import (
	"fmt"
)

// type for weekly schedule of a class / section.
//
// this is just an array of `Day` types.
type WeekTimeTable [N_WEEKLY_SCHOOL_DAYS]DayTimeTable

func (week *WeekTimeTable) Get(day_idx int) *DayTimeTable {
	if day_idx < 0 || day_idx >= N_WEEKLY_SCHOOL_DAYS {
		panic(fmt.Sprintf(
			"GetDaySchedule(day_idx = %d | min:max = 0:%d): error index out of bounds",
			day_idx, N_WEEKLY_SCHOOL_DAYS,
		))
	}

	return &(*week)[day_idx]
}
