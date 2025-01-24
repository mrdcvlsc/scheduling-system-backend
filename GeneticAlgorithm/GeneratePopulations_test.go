package geneticalgorithm_test

import (
	"fmt"
	"testing"

	GA "github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
)

func TestEstimateResourceAvailabilityFirstSem(t *testing.T) {
	err := GA.EstimateResourceAvailability(GA.TERM_1ST_SEMESTER, 0)

	for _, e := range err {
		t.Error(e)
		fmt.Println()
	}
}

func TestEstimateResourceAvailabilitySecondSem(t *testing.T) {
	err := GA.EstimateResourceAvailability(GA.TERM_2ND_SEMESTER, 0)

	for _, e := range err {
		t.Error(e)
		fmt.Println()
	}
}

func TestNewPopulationFirstSem(t *testing.T) {
	GeneratePopulations(t, GA.TERM_1ST_SEMESTER)
}

func TestNewPopulationSecondSem(t *testing.T) {
	GeneratePopulations(t, GA.TERM_2ND_SEMESTER)
}

func GeneratePopulations(t *testing.T, target_semester int) {
	t.Logf("Semester : %d\n\n", target_semester)

	total_test_iterations := 1024
	allowed_generation_errors := 0.2

	generation_error_list := make([]error, 0, 8)
	validation_error_list := make([]error, 0, 8)

	for i := 0; i < total_test_iterations; i++ {

		if (i == 0) || (((i + 1) % 32) == 0) {
			fmt.Printf("Generating schedules (%d)...\n", (i + 1))
		}

		unit_scheds_first_sem, err := GA.NewIndividual(GA.TERM_1ST_SEMESTER, 0)

		if err != nil {
			t.Log(err)
			generation_error_list = append(generation_error_list, err)
		}

		err_validation := unit_scheds_first_sem.Validate()

		if err_validation != nil {
			for e := range err_validation {
				t.Log(e)
			}
			validation_error_list = append(validation_error_list, err_validation...)
		}
	}

	failed_individuals := float64(len(generation_error_list))
	successful_individuals := total_test_iterations - int(failed_individuals)

	if failed_individuals > (float64(total_test_iterations) * allowed_generation_errors) {
		t.Errorf(
			"total of %d fails (%.2f%%) and %d success (%.2f%%) out of the %d populations generated which is above the maximum error treashold of %.2f%%",
			int(failed_individuals), failed_individuals/float64(total_test_iterations)*100.0,
			successful_individuals, float64(successful_individuals)/float64(total_test_iterations)*100.0,
			total_test_iterations, allowed_generation_errors*100.0,
		)
	} else {
		fmt.Printf(
			"total of %d fails (%.2f%%) and %d success (%.2f%%) out of the %d populations generated which is below the maximum error treashold of %.2f%%\n",
			int(failed_individuals), failed_individuals/float64(total_test_iterations)*100.0,
			successful_individuals, float64(successful_individuals)/float64(total_test_iterations)*100.0,
			total_test_iterations, allowed_generation_errors*100.0,
		)
	}

	if len(validation_error_list) > 0 {
		t.Errorf("There are a total of %d validation errors detected when generating university schedules", len(validation_error_list))
	}
}

func BenchmarkNewPopulationFirstSem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GA.NewIndividual(GA.TERM_2ND_SEMESTER, 0)
	}
}
