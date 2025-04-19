package GeneticAlgorithm

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

const MAX_NEW_INDIVIDUAL_GENERATION_TRIALS int = 64
const MAX_RE_ENCODE_TRYS int = 64

type SchedAndResources struct {
	UniSched  Schedule.UniTimeTables
	Resources *EncodingResource
}

// TODO: implement the whole genetic algorithm function

func RunGeneticAlgorithm(
	base_uni_sched Schedule.UniTimeTables,
	curriculums []Curriculum.Curriculum,
	dept_id_to_department map[uint16]Departments.Department,
	default_empty_encoding_resource, base_encoding_resource *EncodingResource,
	department_to_encode map[uint16]bool,
	selected_semester, population_size, generations int,
	resource_persistence *StorageResources.Persistence,
) (Schedule.UniTimeTables, *EncodingResource, error) {

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	department_id := uint16(0)

	if len(department_to_encode) != 1 {
		panic("department to encode should only have one value for now")
	}

	for k := range department_to_encode {
		department_id = k
	}

	////////////////////////////////////////////////////////////////////////////////////////
	//             PUT THE BASE SCHEDULE AT THE TOP OF GENESIS POPULATION
	////////////////////////////////////////////////////////////////////////////////////////

	new_base_sched, new_base_sched_resource, err_encoding_new_base_sched := EncodeIndividualGenome(
		base_uni_sched,
		curriculums,
		dept_id_to_department,
		base_encoding_resource,
		department_to_encode,
		selected_semester, 0,
	)

	if err_encoding_new_base_sched != nil {
		return new_base_sched, new_base_sched_resource, err_encoding_new_base_sched
	}

	log.Printf("university schedule fitness : %f", MeasureCompleteUniSchedBasicFitness(base_uni_sched, curriculums, nil, selected_semester))
	log.Printf("department schedule fitness : %f", MeasureCompleteUniSchedBasicFitness(base_uni_sched, curriculums, department_to_encode, selected_semester))

	genesis_population := make([]SchedAndResources, 0)

	genesis_population = append(genesis_population, SchedAndResources{
		UniSched:  new_base_sched,
		Resources: new_base_sched_resource,
	})

	////////////////////////////////////////////////////////////////////////////////////////
	//               POPULATE THE GENESIS POPULATION WITH RANDOM INDIVIDUALS
	////////////////////////////////////////////////////////////////////////////////////////

	genesis_generation_tries := 0

	log.Print("generating genesis population")

	for len(genesis_population) < population_size {

		copy_uni_sched := make(Schedule.UniTimeTables, len(base_uni_sched))

		copied_week_time_table := copy(copy_uni_sched, base_uni_sched)

		if copied_week_time_table != len(base_uni_sched) {
			return nil, nil, fmt.Errorf("slice elements copied %d, internal university schedule copy operation failed in generate new individual function", copied_week_time_table)
		}

		ApplyClearDepartmentSchedule(copy_uni_sched, curriculums, department_id, selected_semester)

		copy_encoding_resource, err_gen_copy_encoding_resource := GenerateEncodingResourceFromUniTimeTable(copy_uni_sched, curriculums, selected_semester, resource_persistence)

		if err_gen_copy_encoding_resource != nil {
			return nil, nil, fmt.Errorf("unable to generate encoding resource from individual during genesis generation")
		}

		initial_sched, initial_encoding_resource, err_encode_initial := EncodeIndividualGenome(
			copy_uni_sched,
			curriculums,
			dept_id_to_department,
			copy_encoding_resource,
			department_to_encode,
			selected_semester, 0,
		)

		if err_encode_initial != nil {
			genesis_generation_tries++

			if genesis_generation_tries >= MAX_NEW_INDIVIDUAL_GENERATION_TRIALS {
				if initial_sched == nil {
					return nil, nil, fmt.Errorf(
						"unable to generate new a individual for the genesis population after %d tries : %s",
						MAX_NEW_INDIVIDUAL_GENERATION_TRIALS, err_encode_initial.Error(),
					)
				} else {
					return initial_sched, nil, fmt.Errorf(
						"unable to generate new a individual for the genesis population after %d tries : %s",
						MAX_NEW_INDIVIDUAL_GENERATION_TRIALS, err_encode_initial.Error(),
					)
				}
			}

			continue
		} else {
			genesis_generation_tries = 0
		}

		genesis_population = append(genesis_population, SchedAndResources{
			UniSched:  initial_sched,
			Resources: initial_encoding_resource,
		})
	}

	////////////////////////////////////////////////////////////////////////////////////////
	//                   START GENETIC ALGORITHM SCHEDULE GENERATION
	////////////////////////////////////////////////////////////////////////////////////////

	// NOTE: "base schedule" are the current initial or most fit individual in a generation
	// this should never be replaced or modify by the algorithm, the "base schedule"
	// will only change if a new fitter inidividual emerges as a new "base schedule".

	for g := range generations {

		////////////////////////////////////////////////////////////////////////////////////////
		//                     CREATE POPULATION FOR A NEW GENERATION
		////////////////////////////////////////////////////////////////////////////////////////

		log.Printf("running genetic algorithm generation %d", g)

		population := make([]SchedAndResources, 0, population_size+1)

		////////////////////////////////////////////////////////////////////////////////////////
		//	   take the best individual from the previous generation or genesis population
		////////////////////////////////////////////////////////////////////////////////////////

		population = append(population, genesis_population[0])

		////////////////////////////////////////////////////////////////////////////////////////
		//				                TOURNAMENT SELECTION
		////////////////////////////////////////////////////////////////////////////////////////

		log.Printf("apply tournament selection to the population")

		non_elite_population := genesis_population[1:]

		rng.Shuffle(len(non_elite_population), func(i, j int) {
			non_elite_population[i], non_elite_population[j] = non_elite_population[j], non_elite_population[i]
		})

		for i := 0; i < len(non_elite_population)-2; i += 2 {
			A := MeasureCompleteUniSchedBasicFitness(non_elite_population[i].UniSched, curriculums, department_to_encode, selected_semester)
			B := MeasureCompleteUniSchedBasicFitness(non_elite_population[i+1].UniSched, curriculums, department_to_encode, selected_semester)

			if A > B {
				population = append(population, non_elite_population[i])
			} else {
				population = append(population, non_elite_population[i+1])
			}
		}

		remaining_missing_population := population_size - len(population)

		log.Printf("remaining missing population after tournament selection : %d", remaining_missing_population)

		////////////////////////////////////////////////////////////////////////////////////////
		//				                     CROSSOVER
		////////////////////////////////////////////////////////////////////////////////////////

		log.Printf("populate the population with new offspring from parents")

		crossover_tries := 0

		population_size_before_crossover := len(population)

		for len(population) < population_size {

			parent1_idx := rng.Intn(population_size_before_crossover)
			parent2_idx := rng.Intn(population_size_before_crossover)

			if parent1_idx == parent2_idx {
				continue
			}

			parent1 := population[parent1_idx]
			parent2 := population[parent2_idx]

			offspring, err_crossover := Crossover(
				parent1.UniSched, parent2.UniSched,
				curriculums, selected_semester,
				dept_id_to_department, department_to_encode,
				resource_persistence,
			)

			if err_crossover != nil {
				crossover_tries++

				if crossover_tries >= MAX_NEW_INDIVIDUAL_GENERATION_TRIALS {
					log.Printf("failed to crossover some parents %d and %d after %d re-tries", parent1_idx, parent2_idx, crossover_tries)

					if offspring.UniSched == nil {
						return nil, nil, fmt.Errorf(
							"unable to generate new a individual during generation %d after %d tries : %s",
							g, MAX_NEW_INDIVIDUAL_GENERATION_TRIALS, err_crossover.Error(),
						)
					} else {
						return offspring.UniSched, nil, fmt.Errorf(
							"unable to generate new a individual during generation %d after %d tries : %s",
							g, MAX_NEW_INDIVIDUAL_GENERATION_TRIALS, err_crossover.Error(),
						)
					}
				}

				continue
			} else {
				log.Printf("successful crossover parents %d and %d after %d re-tries", parent1_idx, parent2_idx, crossover_tries)
				crossover_tries = 0
			}

			population = append(population, *offspring)
		}

		////////////////////////////////////////////////////////////////////////////////////////
		//				                   RANDOM MUTATIONS
		////////////////////////////////////////////////////////////////////////////////////////

		log.Print("applying random mutation to the population")

		for i := 1; i < len(population); i++ {

			// apply random mutations to some of the CURRENT individuals in the population

			ApplyRandomSubjectErasure(population[i].UniSched, resource_persistence, curriculums, department_id, selected_semester)

			// TODO: when code-base become stable remove vertical validation panic 1

			err_vv1 := population[i].UniSched.VerticalValidation(resource_persistence)

			if len(err_vv1) > 0 {
				panic("vertical validation error 1 : after subject erasure")
			}

			ApplyRandomDaySwapTimeSlots(population[i].UniSched, curriculums, department_id, selected_semester)

			// TODO: when code-base become stable remove vertical validation panic 2

			err_vv2 := population[i].UniSched.VerticalValidation(resource_persistence)

			if len(err_vv2) > 0 {
				for _, err := range err_vv2 {
					fmt.Printf("vertical validation error : %s\n", err.Error())
				}

				panic("vertical validation error 3 : after day swap time slots")
			}

			ApplyRandomSubjectDaySwap(population[i].UniSched, resource_persistence, curriculums, department_id, selected_semester)

			// TODO: when code-base become stable remove vertical validation panic 3

			err_vv3 := population[i].UniSched.VerticalValidation(resource_persistence)

			if len(err_vv3) > 0 {
				panic("vertical validation error 4 : after day swap")
			}

			ApplyRandomSubjectTimeSlotNudge(population[i].UniSched, resource_persistence, curriculums, department_id, selected_semester)

			// TODO: when code-base become stable remove vertical validation panic 4

			err_vv4 := population[i].UniSched.VerticalValidation(resource_persistence)

			if len(err_vv4) > 0 {
				panic("vertical validation error 5 : after time slot nudge")
			}

			ApplyRandomSubjectTimeSlotAndDayNudge(population[i].UniSched, resource_persistence, curriculums, department_id, selected_semester)

			// TODO: when code-base become stable remove vertical validation panic 5

			err_vv5 := population[i].UniSched.VerticalValidation(resource_persistence)

			if len(err_vv5) > 0 {
				panic("vertical validation error 5 : after time slot nudge")
			}

			// re-encode to fix missing geneome/individual's schedule

			re_encode_tries := 0

			for re_encode_tries < MAX_RE_ENCODE_TRYS {

				generated_encoding_resource, err_generate_encoding_resource := GenerateEncodingResourceFromUniTimeTable(
					population[i].UniSched, curriculums, selected_semester, resource_persistence,
				)

				if err_generate_encoding_resource != nil {
					return nil, nil, fmt.Errorf(
						"unable to generate encoding resource from individual during generation %d", g,
					)
				}

				re_encoded_individual, re_encoded_encoding_resource, err_re_encode_schedule := EncodeIndividualGenome(
					population[i].UniSched, curriculums, dept_id_to_department,
					generated_encoding_resource, department_to_encode,
					selected_semester, 0,
				)

				if err_re_encode_schedule != nil {
					re_encode_tries++

					if re_encode_tries >= MAX_RE_ENCODE_TRYS {
						if re_encoded_individual == nil {
							return nil, nil, fmt.Errorf(
								"unable to generate encoding resource from individual during generation %d after %d tries : %s",
								g, MAX_RE_ENCODE_TRYS, err_re_encode_schedule.Error(),
							)
						} else {
							return re_encoded_individual, nil, fmt.Errorf(
								"unable to generate encoding resource from individual during generation %d after %d tries : %s",
								g, MAX_RE_ENCODE_TRYS, err_re_encode_schedule.Error(),
							)
						}
					} else {
						ApplyRandomSubjectErasure(population[i].UniSched, resource_persistence, curriculums, department_id, selected_semester)
						ApplyRandomDaySwapTimeSlots(population[i].UniSched, curriculums, department_id, selected_semester)
						ApplyRandomSubjectDaySwap(population[i].UniSched, resource_persistence, curriculums, department_id, selected_semester)
						ApplyRandomSubjectTimeSlotNudge(population[i].UniSched, resource_persistence, curriculums, department_id, selected_semester)
						ApplyRandomSubjectTimeSlotAndDayNudge(population[i].UniSched, resource_persistence, curriculums, department_id, selected_semester)
					}

					continue
				} else {
					re_encode_tries = 0
				}

				////////////////////////////////////////////////////////////////////////////////////////
				//				          RANDOM MUTATIONS - SANITY CHECK FOR DEBUGGING
				////////////////////////////////////////////////////////////////////////////////////////

				if re_encoded_individual.IsEmpty() {
					panic(">>> re-encoded individual best individual is empty")
				}

				if len(re_encoded_individual.VerticalValidation(resource_persistence)) > 0 {
					panic(">>> re-encoded individual individual has vertical validation error")
				}

				err_hr := re_encoded_individual.HorizontalValidation(resource_persistence, department_to_encode, selected_semester)

				if len(err_hr) > 0 {
					fmt.Println(">>> department to encode:")
					Utils.PrettyPrint(department_to_encode)

					for _, err := range err_hr {
						fmt.Printf("horizontal validation error : %s\n", err.Error())
					}

					panic(">>> re-encoded individual has horizontal validation error")
				}

				if re_encoded_encoding_resource == nil {
					panic("this re-encoding resource is empty")
				}

				if len(re_encoded_encoding_resource.DeptIdToInstructors) <= 0 {
					panic("this re-encoding resource has an empty DeptIdToInstructors")
				}

				if len(re_encoded_encoding_resource.DeptIdToRoomtypeToRooms) <= 0 {
					panic("this re-encoding resource has an empty DeptIdToRoomtypeToRooms")
				}

				if len(re_encoded_encoding_resource.IsSchedIdxToSubIdToSkip) <= 0 {
					panic("this re-encoding resource has an empty IsSchedIdxToSubIdToSkip")
				}

				////////////////////////////////////////////////////////////////////////////////////////
				//				                   RANDOM MUTATIONS
				////////////////////////////////////////////////////////////////////////////////////////

				population[i].UniSched = re_encoded_individual
				population[i].Resources = re_encoded_encoding_resource

				break
			}
		}

		////////////////////////////////////////////////////////////////////////////////////////
		//				PREPARE PREPARE FINAL POPULATION FOR THE NEXT GENERATION
		////////////////////////////////////////////////////////////////////////////////////////

		sort.Slice(population, func(i, j int) bool {
			fitness_a := MeasureCompleteUniSchedBasicFitness(population[i].UniSched, curriculums, department_to_encode, selected_semester)
			fitness_b := MeasureCompleteUniSchedBasicFitness(population[j].UniSched, curriculums, department_to_encode, selected_semester)
			return fitness_a > fitness_b
		})

		// add back to the genesis population

		log.Print("adding back to the genesis population")

		genesis_population = make([]SchedAndResources, 0, population_size+1)
		genesis_population = append(genesis_population, population...)

		log.Printf("best individual fitness : %f", MeasureCompleteUniSchedBasicFitness(genesis_population[0].UniSched, curriculums, department_to_encode, selected_semester))
	}

	////////////////////////////////////////////////////////////////////////////////////////
	//            GENETIC ALGORITHM END : PICK THE BEST INDIVIDUAL SOLUTION
	////////////////////////////////////////////////////////////////////////////////////////

	log.Printf("final best individual fitness : %f", MeasureCompleteUniSchedBasicFitness(genesis_population[0].UniSched, curriculums, department_to_encode, selected_semester))

	if genesis_population[0].UniSched.IsEmpty() {
		panic("final best individual is empty")
	}

	if len(genesis_population[0].UniSched.VerticalValidation(resource_persistence)) > 0 {
		panic("final best individual has vertical validation error")
	}

	if len(genesis_population[0].UniSched.HorizontalValidation(resource_persistence, department_to_encode, selected_semester)) > 0 {
		fmt.Println(">>> department to encode:")
		Utils.PrettyPrint(department_to_encode)
		panic("final best individual has horizontal validation error")
	}

	return genesis_population[0].UniSched, genesis_population[0].Resources, nil
}
