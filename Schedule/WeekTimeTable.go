package Schedule

import (
	"fmt"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
)

// type for weekly schedule of a class / section.
//
// this is just an array of `Day` types.
type WeekTimeTable [Const.N_WEEKLY_SCHOOL_DAYS]DayTimeTable

func (week *WeekTimeTable) GetDayTimeTable(day_idx int) *DayTimeTable {
	if day_idx < 0 || day_idx >= Const.N_WEEKLY_SCHOOL_DAYS {
		panic(fmt.Sprintf(
			"GetDaySchedule(day_idx = %d | min:max = 0:%d): error index out of bounds",
			day_idx, (Const.N_WEEKLY_SCHOOL_DAYS - 1),
		))
	}

	return &week[day_idx]
}
