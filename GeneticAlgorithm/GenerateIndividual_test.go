package GeneticAlgorithm_test

import (
	"fmt"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageSchedule"
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

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		t.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	default_encoding_resource, encoding_resource_err := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

	if encoding_resource_err != nil {
		t.Fatal(encoding_resource_err)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	for i := 0; i < total_test_iterations; i++ {
		if (i == 0) || (((i + 1) % 32) == 0) {
			fmt.Printf("Generating schedules (%d)..................................\n", (i + 1))
		}

		total_sections_calculated := Curriculum.GetTotalNumberOfSections(curriculums, target_semester)
		empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, target_semester)

		t.Logf(
			"the calculated total sections for the semester index %d is %d, and the generated empty schedules contains %d | iter : %d",
			target_semester, total_sections_calculated, len(empty_university_schedule), i,
		)

		if total_sections_calculated != len(empty_university_schedule) {
			t.Fatalf(
				"the calculated total sections for the semester index %d is %d, but the generated empty schedules only contains %d which is a mismatch",
				target_semester, total_sections_calculated, len(empty_university_schedule),
			)
		}

		university_schedules, encoding_resource, err := GeneticAlgorithm.EncodeIndividualGenome(
			empty_university_schedule,
			curriculums,
			dept_id_to_department, default_encoding_resource, nil,
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

		// t.Logf("Schedules Generated : %d,    Calculated Total Number Of Sections : %d", len(university_schedules), total_sections_calculated)

		err_vertical_validations := university_schedules.VerticalValidation(&persistence)

		for _, e := range err_vertical_validations {
			t.Error(e)
			validation_error_list = append(validation_error_list, e)
		}

		///////////////////////

		generated_encoding_resource, gen_encode_resource_err := GeneticAlgorithm.GenerateEncodingResourceFromUniTimeTable(
			university_schedules, curriculums, target_semester, &persistence,
		)

		if gen_encode_resource_err != nil {
			t.Fatal(gen_encode_resource_err)
		}

		if encoding_resource != nil {
			if !GeneticAlgorithm.IsEqualEncodingResource(generated_encoding_resource, encoding_resource) {
				schedules_persistence := StorageSchedule.Persistence{WriterService: &StorageSchedule.JsonWriter{}}
				schedules_persistence.WriterService.SaveSchedules(university_schedules, target_semester)
				t.Fatal("generated encoding resource from bare university schedule is not equal to the produced encoding resource of GA")
			}
		}

		/////////////////

		if err == nil {
			err_horizontal_validations := university_schedules.HorizontalValidation(&persistence, target_semester)

			for _, e := range err_horizontal_validations {
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

//////////////////////////////////////
// WITH DEPARTMENT SPECIFIC ENCODING
/////////////////////////////////////

func TestNewPopulation1stSemWithDepartmentSelection(t *testing.T) {
	GeneratePopWithDepartmentSelection(t, GeneticAlgorithm.TERM_1ST_SEMESTER)
}

func TestNewPopulation2ndSemWithDepartmentSelection(t *testing.T) {
	GeneratePopWithDepartmentSelection(t, GeneticAlgorithm.TERM_2ND_SEMESTER)
}

func GeneratePopWithDepartmentSelection(t *testing.T, target_semester int) {
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

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		t.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	encoding_resource, encoding_resource_err := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

	if encoding_resource_err != nil {
		t.Fatal(encoding_resource_err)
	}

	////////////////////////////////////////////////////////////////////////////////////////

new_population_loop:
	for i := 0; i < total_test_iterations; i++ {
		if (i == 0) || (((i + 1) % 32) == 0) {
			fmt.Printf("Generating schedules (%d)..................................\n", (i + 1))
		}

		total_sections_calculated := Curriculum.GetTotalNumberOfSections(curriculums, target_semester)
		empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, target_semester)

		if total_sections_calculated != len(empty_university_schedule) {
			t.Fatalf(
				"the calculated total sections for the semester index %d is %d, but the generated empty schedules only contains %d which is a mismatch",
				target_semester, total_sections_calculated, len(empty_university_schedule),
			)
		}

		all_departments, all_departments_err := persistence.ReaderService.GetAllDepartments()

		if all_departments_err != nil {
			t.Fatal(all_departments_err)
		}

		is_department_id_to_has_curriculum := make(map[uint16]bool)

		for _, curriculum := range curriculums {
			is_department_id_to_has_curriculum[curriculum.DepartmentID] = true
		}

		track_resources := encoding_resource
		track_schedules := empty_university_schedule

		for department_idx, department := range all_departments {

			if has_curriculum := is_department_id_to_has_curriculum[department.DepartmentID]; !has_curriculum {
				continue // skip departments that don't have curriculums yet
			}

			department_to_encode := make(GeneticAlgorithm.DepartmentsToEncode)
			department_to_encode[department.DepartmentID] = true

			if department_idx <= 0 {
				if !track_schedules.IsEmpty() {
					t.Fatalf("returned a not empty university schedule : loop iteration %d\n", i)
				}
			} else {
				if track_schedules.IsEmpty() {
					t.Fatalf("returned an empty university schedule : loop iteration %d\n", i)
				}
			}

			retries := 0

			var resource_copy_err error
			var not_enough_resource_err error

			max_retries := 7

			for {
				resource_copy_err = nil
				not_enough_resource_err = nil

				output_schedules, output_resources, inner_gen_err := GeneticAlgorithm.EncodeIndividualGenome(
					track_schedules,
					curriculums, dept_id_to_department,
					track_resources, department_to_encode,
					target_semester, 0,
				)

				if output_schedules == nil && output_resources == nil && inner_gen_err != nil {
					resource_copy_err = inner_gen_err
				} else if output_schedules != nil && output_resources == nil && inner_gen_err != nil {
					not_enough_resource_err = inner_gen_err
				}

				if resource_copy_err != nil {
					t.Fatal(resource_copy_err)
				}

				if not_enough_resource_err == nil {
					track_schedules = output_schedules
					track_resources = output_resources
					break
				}

				if retries > max_retries {
					generation_error_list = append(generation_error_list, not_enough_resource_err)
					continue new_population_loop
				}

				t.Logf("retry (%d : %s) - %s\n", retries, department.Code, not_enough_resource_err.Error())
				retries++
			}

			if len(track_schedules) == 0 {
				t.Fatal("No university schedules generated")
			}

			if track_schedules == nil {
				t.Fatalf("returned a nil university schedule : loop iteration %d\n", i)
			}

			err_vertical_validations := track_schedules.VerticalValidation(&persistence)

			for _, e := range err_vertical_validations {
				t.Error(e)
				validation_error_list = append(validation_error_list, e)
			}

			///////////////////////

			generated_encoding_resource, output_resources := GeneticAlgorithm.GenerateEncodingResourceFromUniTimeTable(
				track_schedules, curriculums, target_semester, &persistence,
			)

			if output_resources != nil {
				t.Fatal(output_resources)
			}

			if track_resources != nil {
				if !GeneticAlgorithm.IsEqualEncodingResource(generated_encoding_resource, track_resources) {
					t.Fatal("generated encoding resource from bare university schedule is not equal to the produced encoding resource of GA")
				}
			}

			/////////////////

			if department_idx < len(all_departments)-1 {
				fmt.Printf("Generated schedules for all departments, the department %s\n", department.Name)

				err_horizontal_validations := track_schedules.HorizontalValidation(&persistence, target_semester)

				if err_horizontal_validations == nil {
					t.Fatal("there should be a missing subject error here since the university schedule is not complete yet")
				}

			} else {
				fmt.Printf("Generated schedules for all departments, the last department schedules generated is %s\n", department.Name)

				if track_schedules.IsEmpty() {
					t.Fatalf("returned an empty university schedule : loop iteration %d\n", i)
				}

				err_horizontal_validations := track_schedules.HorizontalValidation(&persistence, target_semester)

				for _, e := range err_horizontal_validations {
					t.Fatal(e)
				}
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

////////////////////////////////
// BENCHMARK
/////////////////////////////////

func BenchmarkNewPopulationFirstSem(b *testing.B) {
	persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_curriculums := persistence.ReaderService.GetAllCurriculum()

	if err_curriculums != nil {
		b.Fatal(err_curriculums)
	}

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		b.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	encoding_resource, encoding_resource_err := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

	if encoding_resource_err != nil {
		b.Fatal(encoding_resource_err)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	for i := 0; i < b.N; i++ {
		empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, GeneticAlgorithm.TERM_1ST_SEMESTER)

		GeneticAlgorithm.EncodeIndividualGenome(
			empty_university_schedule,
			curriculums, dept_id_to_department,
			encoding_resource, nil,
			GeneticAlgorithm.TERM_1ST_SEMESTER, 0,
		)
	}
}

func BenchmarkNewPopulationSecondSem(b *testing.B) {
	persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_curriculums := persistence.ReaderService.GetAllCurriculum()

	if err_curriculums != nil {
		b.Fatal(err_curriculums)
	}

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		b.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	encoding_resource, encoding_resource_err := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

	if encoding_resource_err != nil {
		b.Fatal(encoding_resource_err)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	for i := 0; i < b.N; i++ {
		empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, GeneticAlgorithm.TERM_2ND_SEMESTER)

		GeneticAlgorithm.EncodeIndividualGenome(
			empty_university_schedule,
			curriculums, dept_id_to_department,
			encoding_resource, nil,
			GeneticAlgorithm.TERM_2ND_SEMESTER, 0,
		)
	}
}
