package geneticalgorithm

import (
	"log"
)

type SectionsWeeklySchedules []WeekSchedule

func NewSectionsWeeklySchedules(num_of_sections uint) SectionsWeeklySchedules {
	return make(SectionsWeeklySchedules, num_of_sections)
}

func (university_schedules *SectionsWeeklySchedules) GetSectionWeeklySchedule(section_idx int) *WeekSchedule {
	total_university_sections := len(*university_schedules)

	if section_idx < 0 || section_idx >= total_university_sections {
		log.Fatalf(
			"GetSectionSchedule(section_idx = %d | min:max = 0:%d): error index out of bounds",
			section_idx, total_university_sections,
		)
	}

	return &(*university_schedules)[section_idx]
}
