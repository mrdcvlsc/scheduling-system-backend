package geneticalgorithm_test

import (
	"testing"

	GA "github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
)

func TestNewPopulation(t *testing.T) {

	println("1st Semester Information:")
	GA.NewIndividual(GA.TERM_1ST_SEMESTER, 0)

	println("2nd Semester Information:")
	GA.NewIndividual(GA.TERM_2ND_SEMESTER, 0)

	// if available == false {
	// 	t.Errorf(
	// 		"GetAvailability(day = %d, time_slot = %d) : {loop test phase 4} should not be false again",
	// 		day, time_slot,
	// 	)
	// }
}

func BenchmarkNewPopulationSem1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GA.NewIndividual(GA.TERM_1ST_SEMESTER, 0)
	}
}

func BenchmarkNewPopulationSem2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GA.NewIndividual(GA.TERM_2ND_SEMESTER, 0)
	}
}
