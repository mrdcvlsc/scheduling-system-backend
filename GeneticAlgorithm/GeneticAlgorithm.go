package GeneticAlgorithm

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

const MAX_NEW_INDIVIDUAL_GENERATION_TRIALS int = 64

// TODO: implement the whole genetic algorithm function

func RunGeneticAlgorithm(
	rc_university_schedules Schedule.UniTimeTables,
	rc_curriculums []Curriculum.Curriculum,
	ro_dept_id_to_department map[uint16]Departments.Department,
	default_empty_encoding_resource, rc_encoding_resource *EncodingResource,
	ro_department_to_encode map[uint16]bool,
	selected_semester, population_size, generations int,
) (Schedule.UniTimeTables, *EncodingResource, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	saved_population := make([]Schedule.UniTimeTables, 0)
	saved_resources := make([]*EncodingResource, 0)

	saved_population = append(saved_population, rc_university_schedules)
	saved_resources = append(saved_resources, rc_encoding_resource)

	// NOTE: "base schedule" are the current initial or most fit individual in a generation
	// this should never be replaced or modify by the algorithm, the "base schedule"
	// will only change if a new fitter inidividual emerges as a new "base schedule".

	for g := range generations {
		population := make([]Schedule.UniTimeTables, 0, population_size+1)
		population_resources := make([]*EncodingResource, 0, population_size+1)

		population = append(population, saved_population...)
		population_resources = append(population_resources, saved_resources...)

		// add new individuals

		tries := 0

		// for len(population) < population_size/2 {
		for len(population) < population_size {
			empty_uni_sched := NewEmptyIndividual(rc_curriculums, selected_semester)

			new_uni_sched, encoding_resource, err_encode_new := EncodeIndividualGenome(
				empty_uni_sched,
				rc_curriculums,
				ro_dept_id_to_department,
				default_empty_encoding_resource,
				ro_department_to_encode,
				selected_semester, 0,
			)

			if err_encode_new != nil {
				tries++

				if tries >= MAX_NEW_INDIVIDUAL_GENERATION_TRIALS {
					return nil, nil, fmt.Errorf(
						"unable to generate new a individual during generation %d after %d tries",
						g, MAX_NEW_INDIVIDUAL_GENERATION_TRIALS,
					)
				}

				continue
			}

			population = append(population, new_uni_sched)
			population_resources = append(population_resources, encoding_resource)
		}

		// TODO: crossover some individuals and add them to the CURRENT population

		// TODO: apply random mutations to some of the CURRENT individuals in the population

		// TODO: apply tournament selection to be saved to the next generation

		new_population := population[1:]
		new_resources := population_resources[1:]

		rng.Shuffle(len(new_population), func(i, j int) {
			new_population[i], new_population[j] = new_population[j], new_population[i]
			new_resources[i], new_resources[j] = new_resources[j], new_resources[i]
		})

		saved_population = make([]Schedule.UniTimeTables, 0, population_size)
		saved_resources = make([]*EncodingResource, 0, population_size)

		saved_population = append(saved_population, population[0])
		saved_resources = append(saved_resources, population_resources[0])

		// TODO: sort individuals base on fitness
	}

	return nil, nil, nil
}
