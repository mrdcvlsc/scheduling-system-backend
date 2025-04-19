package GeneticAlgorithm

import (
	"math"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

const PREFERED_MAX_CLASS_HOUR_PER_DAY float64 = 7.0

/*
* Output Range:

	(0, 1]

* Behavior:

	When actual_hours = target_hours, return = 1
	As |actual_hours - target_hours| --> +Inf, return --> 0
*/
func reciprocal_distance(actual_hours, target_hours float64) float64 {
	return 1.0 / (1.0 + math.Abs(actual_hours-target_hours))

}

func MeasureWeekTimeTableBasicFitness(week_sched Schedule.WeekTimeTable) float64 {
	week_sched_fitness := 0.0

	for day := range Const.N_WEEKLY_SCHOOL_DAYS {

		has_class_after_5pm := false
		has_time_for_lunch := false
		day_total_hours := 0.0
		continuous_hours := 0.0

		for time_slot := range Const.N_DAILY_TIME_SLOTS {
			if time_slot >= 20 && continuous_hours > 0 {
				has_class_after_5pm = true
			}

			if week_sched[day][time_slot].GetSubjectID() > 0 {
				day_total_hours += (1.0 / float64(Const.N_HOUR_TIME_SLOTS))
				continuous_hours += (1.0 / float64(Const.N_HOUR_TIME_SLOTS))
			} else {
				if continuous_hours > 4.5 {
					week_sched_fitness -= 0.75
				}

				continuous_hours = 0.0
			}

			if time_slot >= 8 && time_slot <= 12 && week_sched[day][time_slot].GetSubjectID() == 0 {
				has_time_for_lunch = true
			}
		}

		if has_time_for_lunch {
			week_sched_fitness += 7.0
		} else {
			week_sched_fitness -= 2.0
		}

		if has_class_after_5pm {
			week_sched_fitness -= 0.5
		}

		if day_total_hours >= PREFERED_MAX_CLASS_HOUR_PER_DAY {
			fitness_punishment := day_total_hours - PREFERED_MAX_CLASS_HOUR_PER_DAY
			week_sched_fitness -= fitness_punishment * 0.75
		} else if day_total_hours == 0 {
			week_sched_fitness = 0
			continue
		} else {
			week_sched_fitness += 3.5
		}

		// no class during saturday.
		if day_total_hours > 4 && day == Const.N_WEEKLY_SCHOOL_DAYS-1 {
			week_sched_fitness -= 0.75
		}

	}

	return week_sched_fitness
}

// A basic fitness function
func MeasureCompleteUniSchedBasicFitness(complete_uni_sched Schedule.UniTimeTables, all_curriculums []Curriculum.Curriculum, department_to_measure map[uint16]bool, selected_semester int) float64 {
	if complete_uni_sched.IsEmpty() {
		return 0.0
	}

	accumulated_fitness := 0.0
	total_fitness_measurements := 0

	IterateSectionsWeekSchedule(complete_uni_sched, all_curriculums, selected_semester, nil, nil, func(indicies IterIndices, values IterValues) IterReturnType {

		if len(department_to_measure) > 0 {
			is_to_measure, has_key := department_to_measure[values.Curriculum.DepartmentID]

			if !(has_key && is_to_measure) {
				return IterProceed
			}
		}

		accumulated_fitness += MeasureWeekTimeTableBasicFitness(*values.WeekSched)
		total_fitness_measurements++

		return IterProceed
	})

	return accumulated_fitness / float64(total_fitness_measurements)
}

func MeasureFitnessPrefHeatMapComparison(uni_sched Schedule.DayTimeTable) float64 {
	// TODO: implement preference heat map comparison based fitness function

	return 1.0
}
