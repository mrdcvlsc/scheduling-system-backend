package scheduledsa

import (
	"fmt"
)

// The type that represent all of the weekly schedules of each classes / sections
// in the whole university.
//
// this is just an array of `ScheduleWeek` types.
type SchedulesInUniversity []ScheduleWeek

func NewUniversitySchedules(num_of_sections uint) SchedulesInUniversity {
	return make(SchedulesInUniversity, num_of_sections)
}

func (university_schedules *SchedulesInUniversity) GetSectionWeeklySchedule(section_idx int) *ScheduleWeek {
	total_university_sections := len(*university_schedules)

	if section_idx < 0 || section_idx >= total_university_sections {
		panic(fmt.Sprintf(
			"GetSectionSchedule(section_idx = %d | min:max = 0:%d): error index out of bounds",
			section_idx, total_university_sections,
		))
	}

	return &(*university_schedules)[section_idx]
}
