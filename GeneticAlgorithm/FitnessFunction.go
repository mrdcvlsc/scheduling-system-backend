package GeneticAlgorithm

import (
	"math"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

const PREFERED_CLASS_HOUR_PER_DAY float64 = 7.0

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

// A basic fitness function
func MeasureFitnessBasic(uni_sched Schedule.UniTimeTables) float64 {
	accumulated_fitness := 0.0

	for i := range uni_sched {
		section_sched_fitness := 0.0

		for day := range Const.N_WEEKLY_SCHOOL_DAYS {

			has_class_after_5pm := false
			total_hours := 0.0
			continuous_hours := 0.0

			for time_slot := range Const.N_DAILY_TIME_SLOTS {
				if time_slot >= 20 && continuous_hours > 0 {
					has_class_after_5pm = true
				}

				if uni_sched[i][day][time_slot].GetSubjectID() > 0 {
					total_hours += (1.0 / float64(Const.N_HOUR_TIME_SLOTS))
					continuous_hours += (1.0 / float64(Const.N_HOUR_TIME_SLOTS))
				} else {
					// check if there's a break after 3 to 5 hour class
					if continuous_hours >= 2.9 && continuous_hours <= 5.1 {
						section_sched_fitness += 1.0
					}

					if continuous_hours > 5.1 {
						section_sched_fitness -= 0.5
					}

					continuous_hours = 0.0
				}
			}

			if has_class_after_5pm {
				section_sched_fitness -= 0.5
			}

			if total_hours > 0.0 {
				section_sched_fitness += reciprocal_distance(total_hours, PREFERED_CLASS_HOUR_PER_DAY)
			} else {
				section_sched_fitness += 0.75
			}

			// no class during saturday.
			if total_hours <= 0.001 && day == Const.N_WEEKLY_SCHOOL_DAYS-1 {
				section_sched_fitness += 1
			}
		}

		accumulated_fitness += section_sched_fitness
	}

	return accumulated_fitness / float64(len(uni_sched))
}

func MeasureFitnessPrefHeatMapComparison(uni_sched Schedule.DayTimeTable) float64 {
	// TODO: implement preference heat map comparison based fitness function

	return 1.0
}
