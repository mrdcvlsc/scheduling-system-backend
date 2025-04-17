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
	base_uni_sched Schedule.UniTimeTables,
	curriculums []Curriculum.Curriculum,
	dept_id_to_department map[uint16]Departments.Department,
	default_empty_encoding_resource, base_encoding_resource *EncodingResource,
	department_to_encode map[uint16]bool,
	selected_semester, population_size, generations int,
) (Schedule.UniTimeTables, *EncodingResource, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	genesis_population := make([]Schedule.UniTimeTables, 0)
	genesis_resources := make([]*EncodingResource, 0)

	genesis_population = append(genesis_population, base_uni_sched)
	genesis_resources = append(genesis_resources, base_encoding_resource)

	genesis_generation_tries := 0
	for len(genesis_population) < population_size {
		empty_uni_sched := NewEmptyIndividual(curriculums, selected_semester)

		new_uni_sched, encoding_resource, err_encode_new := EncodeIndividualGenome(
			empty_uni_sched,
			curriculums,
			dept_id_to_department,
			default_empty_encoding_resource,
			department_to_encode,
			selected_semester, 0,
		)

		if err_encode_new != nil {
			genesis_generation_tries++

			if genesis_generation_tries >= MAX_NEW_INDIVIDUAL_GENERATION_TRIALS {
				return nil, nil, fmt.Errorf(
					"unable to generate new a individual for the genesis population after %d tries",
					MAX_NEW_INDIVIDUAL_GENERATION_TRIALS,
				)
			}

			continue
		} else {
			genesis_generation_tries = 0
		}

		genesis_population = append(genesis_population, new_uni_sched)
		genesis_resources = append(genesis_resources, encoding_resource)
	}

	// NOTE: "base schedule" are the current initial or most fit individual in a generation
	// this should never be replaced or modify by the algorithm, the "base schedule"
	// will only change if a new fitter inidividual emerges as a new "base schedule".

	for g := range generations {
		population := make([]Schedule.UniTimeTables, 0, population_size+1)
		population_resources := make([]*EncodingResource, 0, population_size+1)

		population = append(population, genesis_population...)
		population_resources = append(population_resources, genesis_resources...)

		// add new individuals

		tries := 0

		// for len(population) < population_size/2 {
		for len(population) < population_size {
			empty_uni_sched := NewEmptyIndividual(curriculums, selected_semester)

			new_uni_sched, encoding_resource, err_encode_new := EncodeIndividualGenome(
				empty_uni_sched,
				curriculums,
				dept_id_to_department,
				default_empty_encoding_resource,
				department_to_encode,
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

		new_population := population[1:]
		new_resources := population_resources[1:]

		// TODO: crossover some individuals and add them to the CURRENT population

		// TODO: apply random mutations to some of the CURRENT individuals in the population

		// TODO: after random mutations and crossover, some individuals might be incomplete so
		// we need to re-encode them again to complete their schedules.

		// TODO: apply tournament selection to be saved to the next generation

		rng.Shuffle(len(new_population), func(i, j int) {
			new_population[i], new_population[j] = new_population[j], new_population[i]
			new_resources[i], new_resources[j] = new_resources[j], new_resources[i]
		})

		genesis_population = make([]Schedule.UniTimeTables, 0, population_size)
		genesis_resources = make([]*EncodingResource, 0, population_size)

		genesis_population = append(genesis_population, population[0])
		genesis_resources = append(genesis_resources, population_resources[0])

		// TODO: sort individuals base on fitness
	}

	return nil, nil, nil
}
