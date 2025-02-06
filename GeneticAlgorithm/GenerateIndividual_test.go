package GeneticAlgorithm_test

import (
	"fmt"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
)

func TestEstimateResourceAvailabilityFirstSem(t *testing.T) {
	persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	err := GeneticAlgorithm.EstimateResourceAvailability(&persistence, GeneticAlgorithm.TERM_1ST_SEMESTER, 0)

	for _, e := range err {
		t.Error(e)
		fmt.Println()
	}
}

func TestEstimateResourceAvailabilitySecondSem(t *testing.T) {
	persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	err := GeneticAlgorithm.EstimateResourceAvailability(&persistence, GeneticAlgorithm.TERM_2ND_SEMESTER, 0)

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
	persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	t.Logf("Semester : %d\n\n", target_semester)

	total_test_iterations := 512
	allowed_generation_errors := 0.8 // 80% error rate allowed.

	generation_error_list := make([]error, 0, 8)
	validation_error_list := make([]error, 0, 8)

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_curriculums := persistence.ReaderService.GetAllCurriculum()

	if err_curriculums != nil {
		t.Fatal(err_curriculums)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_room_type_to_rooms, err_dept_id_to_room_type_to_rooms := GeneticAlgorithm.GenerateMapDeptIdToRoomTypeToRooms(&persistence)

	if err_dept_id_to_room_type_to_rooms != nil {
		t.Fatal(err_dept_id_to_room_type_to_rooms)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_instructors, err_dept_id_to_instructors := GeneticAlgorithm.GenerateMapDeptIdToInstructors(&persistence)

	if err_dept_id_to_instructors != nil {
		t.Fatal(err_dept_id_to_instructors)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		t.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	for i := 0; i < total_test_iterations; i++ {
		if (i == 0) || (((i + 1) % 32) == 0) {
			fmt.Printf("Generating schedules (%d)..................................\n", (i + 1))
		}

		university_schedules, err := GeneticAlgorithm.EncodeIndividualGenome(
			curriculums,
			dept_id_to_department,
			dept_id_to_instructors,
			dept_id_to_room_type_to_rooms,
			target_semester, 0,
		)

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

		t.Logf("Schedules Generated : %d", len(university_schedules))

		err_vertical_validations := university_schedules.VerticalValidation(&persistence)

		for _, e := range err_vertical_validations {
			t.Error(e)
			validation_error_list = append(validation_error_list, e)
		}

		if err == nil {
			err_horizontal_validations := university_schedules.HorizontalValidation(&persistence, target_semester)

			for _, e := range err_horizontal_validations {
				// schedules_persistence := StorageSchedule.Persistence{WriterService: &StorageSchedule.JsonWriter{}}
				// schedules_persistence.WriterService.SaveSchedules(university_schedules, target_semester)
				t.Fatal(e)
			}
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
	persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_curriculums := persistence.ReaderService.GetAllCurriculum()

	if err_curriculums != nil {
		b.Fatal(err_curriculums)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_room_type_to_rooms, err_dept_id_to_room_type_to_rooms := GeneticAlgorithm.GenerateMapDeptIdToRoomTypeToRooms(&persistence)

	if err_dept_id_to_room_type_to_rooms != nil {
		b.Fatal(err_dept_id_to_room_type_to_rooms)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_instructors, err_dept_id_to_instructors := GeneticAlgorithm.GenerateMapDeptIdToInstructors(&persistence)

	if err_dept_id_to_instructors != nil {
		b.Fatal(err_dept_id_to_instructors)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		b.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	for i := 0; i < b.N; i++ {
		GeneticAlgorithm.EncodeIndividualGenome(
			curriculums,
			dept_id_to_department,
			dept_id_to_instructors,
			dept_id_to_room_type_to_rooms,
			GeneticAlgorithm.TERM_2ND_SEMESTER, 0,
		)
	}
}
