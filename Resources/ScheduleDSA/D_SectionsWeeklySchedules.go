package scheduledsa

import (
	"fmt"
)

// The type that represent all of the weekly schedules of each classes / sections
// in the whole university.
//
// this is just an array of `ScheduleWeek` types.
type UniTimeTables []WeekTimeTable

func NewUniTimeTables(num_of_time_tables uint) UniTimeTables {
	return make(UniTimeTables, num_of_time_tables)
}

func (uni_sched *UniTimeTables) Get(class_section_idx int) *WeekTimeTable {
	total_university_sections := len(*uni_sched)

	if class_section_idx < 0 || class_section_idx >= total_university_sections {
		panic(fmt.Sprintf(
			"GetSectionSchedule(section_idx = %d | min:max = 0:%d): error index out of bounds",
			class_section_idx, total_university_sections,
		))
	}

	return &(*uni_sched)[class_section_idx]
}
