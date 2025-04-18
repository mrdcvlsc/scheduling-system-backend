package GeneticAlgorithm

import (
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

type IterReturnType int

const (
	IterProceed             IterReturnType = 0
	IterBreakCurriculumLoop IterReturnType = 1
	IterContinue            IterReturnType = 2
	IterBreak               IterReturnType = 3
)

type IterIndices struct {
	Usi        int // university time table schedule index
	Curriculum int // current curriculum index
	YearLevel  int // current year level index
	Semester   int // current semester index
	Section    int // current section index
}

type IterValues struct {
	Sched      Schedule.UniTimeTables  // whole university time table schedule
	WeekSched  *Schedule.WeekTimeTable // current week time table schedule
	Curriculum *Curriculum.Curriculum  // current curriculum
	YearLevel  *Curriculum.YearLevel   // current year level
	Semester   *Curriculum.Semester    // current semester
}

func IterateSectionsWeekSchedule(
	sched Schedule.UniTimeTables,
	curriculums []Curriculum.Curriculum,
	selected_semester int,

	fn_curriculum_block func(indicies IterIndices, values IterValues) IterReturnType,
	fn_semester_block func(indicies IterIndices, values IterValues) IterReturnType,

	fn_section_block func(
		indicies IterIndices, values IterValues,
	) IterReturnType,
) {
	usi := 0

curriculum_loop:
	for curriculum_idx, curriculum := range curriculums {

		if fn_curriculum_block != nil {
			irt_curriculum_block := fn_curriculum_block(
				IterIndices{
					Usi:        usi,
					Curriculum: curriculum_idx,
				},
				IterValues{
					Sched:      sched,
					Curriculum: &curriculum,
				},
			)

			switch irt_curriculum_block {
			case IterBreakCurriculumLoop:
				break curriculum_loop
			case IterContinue:
				continue
			case IterBreak:
				break
			}
		}

		for year_level_idx, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue // skip inactive year levels
			}

			for semester_idx, semester := range year_level.Semesters {

				if selected_semester != semester_idx {
					continue // skip not selected semesters
				}

				if len(semester.Subjects) == 0 {
					continue // skip semesters that don't have subjects
				}

				if fn_semester_block != nil {
					irt_semester_block := fn_semester_block(
						IterIndices{
							Usi:        usi,
							Curriculum: curriculum_idx,
							YearLevel:  year_level_idx,
							Semester:   semester_idx,
						},
						IterValues{
							Sched:      sched,
							Curriculum: &curriculum,
							YearLevel:  &year_level,
							Semester:   &semester,
						},
					)

					switch irt_semester_block {
					case IterBreakCurriculumLoop:
						break curriculum_loop
					case IterContinue:
						continue
					case IterBreak:
						break
					}
				}

				for section_idx := 0; section_idx < semester.Sections; section_idx++ {

					if fn_section_block != nil {

						var week_time_table *Schedule.WeekTimeTable

						if len(sched) > 0 {
							week_time_table = &sched[usi]
						} else {
							week_time_table = nil
						}

						irt_section_block := fn_section_block(
							IterIndices{
								Usi:        usi,
								Curriculum: curriculum_idx,
								YearLevel:  year_level_idx,
								Semester:   semester_idx,
								Section:    section_idx,
							},
							IterValues{
								Sched:      sched,
								WeekSched:  week_time_table,
								Curriculum: &curriculum,
								YearLevel:  &year_level,
								Semester:   &semester,
							},
						)

						switch irt_section_block {
						case IterBreakCurriculumLoop:
							break curriculum_loop
						case IterContinue:
							continue
						case IterBreak:
							break
						}
					}

					usi++
				}
			}
		}
	}
}
