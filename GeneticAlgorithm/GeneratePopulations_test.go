package geneticalgorithm_test

import (
	"fmt"
	"testing"

	GA "github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
)

func TestNewPopulationEnoughRooms(t *testing.T) {

	fmt.Print("\n################ First Semester ################ \n\n")
	GA.EstimateResourceAvailability(GA.TERM_1ST_SEMESTER, 0)

	// fmt.Print("\n################ Second Semester ################ \n\n")
	// GA.EstimateResourceAvailability(1, 0)
}

func TestNewPopulationFirstSem(t *testing.T) {

	println("1st Semester Information:")
	uni_time_table1, _ := GA.NewIndividual(GA.TERM_1ST_SEMESTER, 0)
	uni1errs := uni_time_table1.Validate()

	for _, err := range uni1errs {
		t.Error(err)
	}
}

// func TestNewPopulationSecondSem(t *testing.T) {

// 	println("2nd Semester Information:")
// 	uni_time_table2, _ := GA.NewIndividual(GA.TERM_2ND_SEMESTER, 0)

// 	uni2errs := uni_time_table2.Validate()

// 	for _, err := range uni2errs {
// 		t.Error(err)
// 	}
// }

func BenchmarkNewPopulationSem1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GA.NewIndividual(GA.TERM_1ST_SEMESTER, 0)
	}
}

// func BenchmarkNewPopulationSem2(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		GA.NewIndividual(GA.TERM_2ND_SEMESTER, 0)
// 	}
// }
