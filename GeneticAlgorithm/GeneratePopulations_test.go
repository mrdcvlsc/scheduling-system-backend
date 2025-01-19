package geneticalgorithm_test

import (
	"fmt"
	"testing"

	GA "github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
)

func TestNewPopulation(t *testing.T) {

	println("1st Semester Information:")
	uni_time_table1, _ := GA.NewIndividual(GA.TERM_1ST_SEMESTER, 0)

	println("2nd Semester Information:")
	uni_time_table2, _ := GA.NewIndividual(GA.TERM_2ND_SEMESTER, 0)

	fmt.Println("=================================")

	fmt.Print("\nTime Slot section(0) BEFORE:\n")
	fmt.Printf("%+v", uni_time_table1.Get(0).Get(5).Get(23))
	fmt.Print("\n\n")

	fmt.Print("\nTime Slot section(1) BEFORE:\n")
	fmt.Printf("%+v", uni_time_table1.Get(1).Get(5).Get(23))
	fmt.Print("\n\n")

	fmt.Print("\nTime Slot section(2) BEFORE:\n")
	fmt.Printf("%+v", uni_time_table1.Get(2).Get(5).Get(23))
	fmt.Print("\n\n")

	fmt.Println("=================================")

	// if err := uni_time_table1.Get(0).Get(5).Get(23).SetInstructorID(33); err != nil {
	// 	t.Error(err)
	// }

	// if err := uni_time_table1.Get(7).Get(5).Get(23).SetInstructorID(33); err != nil {
	// 	t.Error(err)
	// }

	// if err := uni_time_table1.Get(2).Get(5).Get(23).SetInstructorID(33); err != nil {
	// 	fmt.Println(err.Error())
	// 	t.Error(err)
	// }

	fmt.Println("=================================")

	fmt.Print("\nTime Slot section(0) AFTER:\n")
	fmt.Printf("%+v", uni_time_table1.Get(0).Get(5).Get(23))
	fmt.Print("\n\n")

	fmt.Print("\nTime Slot section(1) AFTER:\n")
	fmt.Printf("%+v", uni_time_table1.Get(1).Get(5).Get(23))
	fmt.Print("\n\n")

	fmt.Print("\nTime Slot section(2) AFTER:\n")
	fmt.Printf("%+v", uni_time_table1.Get(2).Get(5).Get(23))
	fmt.Print("\n\n")

	uni1errs := uni_time_table1.Validate()
	uni2errs := uni_time_table2.Validate()

	fmt.Println("=================================")

	for _, err := range uni1errs {
		t.Error(err)
	}

	for _, err := range uni2errs {
		t.Error(err)
	}
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
