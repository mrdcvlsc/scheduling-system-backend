package GeneticAlgorithm_test

import (
	"fmt"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
)

func TestEstimateResourceAvailabilityFirstSem(t *testing.T) {
	err := GeneticAlgorithm.EstimateResourceAvailability(GeneticAlgorithm.TERM_1ST_SEMESTER, 0)

	for _, e := range err {
		t.Error(e)
		fmt.Println()
	}
}

func TestEstimateResourceAvailabilitySecondSem(t *testing.T) {
	err := GeneticAlgorithm.EstimateResourceAvailability(GeneticAlgorithm.TERM_2ND_SEMESTER, 0)

	for _, e := range err {
		t.Error(e)
		fmt.Println()
	}
}

func TestNewPopulationFirstSem(t *testing.T) {
	GeneratePopulations(t, GeneticAlgorithm.TERM_1ST_SEMESTER)
}

func TestNewPopulationSecondSem(t *testing.T) {
	GeneratePopulations(t, GeneticAlgorithm.TERM_2ND_SEMESTER)
}

func GeneratePopulations(t *testing.T, target_semester int) {
	t.Logf("Semester : %d\n\n", target_semester)

	total_test_iterations := 512
	allowed_generation_errors := 0.8 // 80% error rate allowed.

	generation_error_list := make([]error, 0, 8)
	validation_error_list := make([]error, 0, 8)

	for i := 0; i < total_test_iterations; i++ {
		if (i == 0) || (((i + 1) % 32) == 0) {
			fmt.Printf("Generating schedules (%d)..................................\n", (i + 1))
		}

		university_schedules, err := GeneticAlgorithm.NewIndividual(target_semester, 0)

		if len(university_schedules) == 0 {
			t.Fatal("No university schedules generated")
		}

		if err != nil {
			t.Log(err)
			generation_error_list = append(generation_error_list, err)
		}

		if university_schedules == nil {
			t.Fatalf("returned a nil university schedule : loop iteration %d\n", i)
		}

		if university_schedules.IsEmpty() {
			t.Fatalf("returned an empty university schedule : loop iteration %d\n", i)
		}

		err_validation := university_schedules.Validate()

		t.Logf("Schedules Generated : %d", len(university_schedules))

		for _, e := range err_validation {
			t.Error(e)
			validation_error_list = append(validation_error_list, e)
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

	t.Logf("There are a total of %d validation errors detected when generating university schedules", len(validation_error_list))
}

func BenchmarkNewPopulationFirstSem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GeneticAlgorithm.NewIndividual(GeneticAlgorithm.TERM_2ND_SEMESTER, 0)
	}
}
