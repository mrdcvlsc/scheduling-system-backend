package GeneticAlgorithm_test

import (
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
)

func TestEmptyScheduleFitness(t *testing.T) {
	uni_sched := make(Schedule.UniTimeTables, 0)
	uni_sched = append(uni_sched, Schedule.WeekTimeTable{})

	fitness := GeneticAlgorithm.MeasureFitnessBasic(uni_sched)
	t.Logf("empty schedule fitness : %f", fitness)

	for range 50 {
		uni_sched = append(uni_sched, Schedule.WeekTimeTable{})
	}

	fitness_of_51_empty_schedules := GeneticAlgorithm.MeasureFitnessBasic(uni_sched)

	// TODO: complete this empty schedule fitness test

	t.Logf("50x empty schedule fitness : %f", fitness_of_51_empty_schedules)
}

func TestGeneratedScheduleFitness(t *testing.T) {
	persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}
	target_semester := 0

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_all_curriculums := persistence.ReaderService.ReadAllCurriculum()

	if err_all_curriculums != nil {
		t.Fatal(err_all_curriculums)
	}

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		t.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	default_encoding_resource, err_read_default_encoding_resource := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

	if err_read_default_encoding_resource != nil {
		t.Fatal(err_read_default_encoding_resource)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	var new_university_schedule *Schedule.UniTimeTables

	for i := 0; i < 50; i++ {
		empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, target_semester)

		university_schedules, encoding_resource, err := GeneticAlgorithm.EncodeIndividualGenome(
			empty_university_schedule,
			curriculums,
			dept_id_to_department, default_encoding_resource, nil,
			target_semester, 0,
		)

		if len(university_schedules) == 0 {
			t.Fatal("No university schedules generated")
		}

		if university_schedules == nil {
			t.Fatalf("returned a nil university schedule : loop iteration %d\n", i)
		}

		if university_schedules.IsEmpty() {
			t.Fatalf("returned an empty university schedule : loop iteration %d\n", i)
		}

		err_vertical_validations := university_schedules.VerticalValidation(&persistence)

		for _, e := range err_vertical_validations {
			t.Error(e)
		}

		if err == nil {
			err_horizontal_validations := university_schedules.HorizontalValidation(&persistence, nil, target_semester)

			for _, e := range err_horizontal_validations {
				t.Fatal(e)
			}
		}

		if encoding_resource != nil && err == nil {
			new_university_schedule = &university_schedules
			break
		}
	}

	fitness := GeneticAlgorithm.MeasureFitnessBasic(*new_university_schedule)

	// TODO: complete this generated schedule fitness test

	t.Logf("generated schedule fitness : %f", fitness)
}
