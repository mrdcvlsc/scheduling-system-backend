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
	allowed_failure_rate := 0.8 // 80% error rate allowed.

	err_list_generation := make([]error, 0, 8)
	err_list_validation := make([]error, 0, 8)

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
			err_list_generation = append(err_list_generation, err)
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
			err_list_validation = append(err_list_validation, e)
		}

		///////////////////////

		generated_encoding_resource, err_gen_encode_resource := GeneticAlgorithm.GenerateEncodingResourceFromUniTimeTable(
			university_schedules, curriculums, target_semester, &persistence,
		)

		if err_gen_encode_resource != nil {
			t.Fatal(err_gen_encode_resource)
		}

		if encoding_resource != nil {
			if !GeneticAlgorithm.IsEqualEncodingResource(generated_encoding_resource, encoding_resource) {
				schedules_persistence := StorageSchedule.Persistence{SaveService: &StorageSchedule.JsonWriter{}}
				schedules_persistence.SaveService.SaveSchedules(university_schedules, target_semester)
				t.Fatal("generated encoding resource from bare university schedule is not equal to the produced encoding resource of GA")
			}
		}

		/////////////////

		if err == nil {
			err_horizontal_validations := university_schedules.HorizontalValidation(&persistence, nil, target_semester)

			for _, e := range err_horizontal_validations {
				t.Fatal(e)
			}
		}
	}

	failed_individuals := float64(len(err_list_generation))
	successful_individuals := total_test_iterations - int(failed_individuals)

	if failed_individuals > (float64(total_test_iterations) * allowed_failure_rate) {
		t.Errorf(
			"total of %d fails (%.2f%%) and %d success (%.2f%%) out of the %d populations generated which is above the maximum error treashold of %.2f%%",
			int(failed_individuals), failed_individuals/float64(total_test_iterations)*100.0,
			successful_individuals, float64(successful_individuals)/float64(total_test_iterations)*100.0,
			total_test_iterations, allowed_failure_rate*100.0,
		)
	} else {
		fmt.Printf(
			"total of %d fails (%.2f%%) and %d success (%.2f%%) out of the %d populations generated which is below the maximum error treashold of %.2f%%\n",
			int(failed_individuals), failed_individuals/float64(total_test_iterations)*100.0,
			successful_individuals, float64(successful_individuals)/float64(total_test_iterations)*100.0,
			total_test_iterations, allowed_failure_rate*100.0,
		)
	}

	t.Logf("There are a total of %d validation errors detected when generating university schedules", len(err_list_validation))
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

	total_test_iterations := 256
	allowed_failure_rate := 0.8 // 80% error rate allowed.

	err_list_generation := make([]error, 0, 8)
	err_list_validation := make([]error, 0, 8)

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

	encoding_resource, err_read_default_encoding_resource := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

	if err_read_default_encoding_resource != nil {
		t.Fatal(err_read_default_encoding_resource)
	}

	////////////////////////////////////////////////////////////////////////////////////////

new_population_loop:
	for i := range total_test_iterations {
		if (i == 0) || (((i + 1) % 16) == 0) {
			fmt.Printf("Generating university schedules [per-department] (%d)..................................\n", (i + 1))
		}

		total_sections_calculated := Curriculum.GetTotalNumberOfSections(curriculums, target_semester)
		empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, target_semester)

		if total_sections_calculated != len(empty_university_schedule) {
			t.Fatalf(
				"the calculated total sections for the semester index %d is %d, but the generated empty schedules only contains %d which is a mismatch",
				target_semester, total_sections_calculated, len(empty_university_schedule),
			)
		}

		all_departments, err_read_all_departments := persistence.ReaderService.ReadAllDepartments()

		if err_read_all_departments != nil {
			t.Fatal(err_read_all_departments)
		}

		is_department_id_to_has_curriculum := make(map[uint16]bool)

		for _, curriculum := range curriculums {
			is_department_id_to_has_curriculum[curriculum.DepartmentID] = true
		}

		track_resources := encoding_resource
		track_schedules := empty_university_schedule

		if !track_schedules.IsEmpty() {
			t.Fatal("returned not empty university schedule : before loop")
		}

		all_departments = all_departments[1:]

		for department_idx, department := range all_departments {

			// setup specific department for schedule generation

			if has_curriculum := is_department_id_to_has_curriculum[department.DepartmentID]; !has_curriculum {
				continue // skip departments that don't have curriculums yet
			}

			department_to_encode := make(map[uint16]bool)
			department_to_encode[department.DepartmentID] = true

			if track_schedules.IsEmpty() {
				t.Fatalf("returned an empty university schedule : loop iteration %d\n", i)
			}

			var err_copy_resource error
			var err_not_enough_resource error

			retries := 0
			max_retries := 7

			// start generating specific department schedule

			for {
				err_copy_resource = nil
				err_not_enough_resource = nil

				output_schedules, output_resources, err_encode_individual_genome := GeneticAlgorithm.EncodeIndividualGenome(
					track_schedules,
					curriculums, dept_id_to_department,
					track_resources, department_to_encode,
					target_semester, 0,
				)

				if output_schedules == nil && output_resources == nil && err_encode_individual_genome != nil {
					err_copy_resource = err_encode_individual_genome
				} else if output_schedules != nil && output_resources == nil && err_encode_individual_genome != nil {
					err_not_enough_resource = err_encode_individual_genome
				}

				if err_copy_resource != nil {
					t.Fatal(err_copy_resource)
				}

				if err_not_enough_resource != nil {
					retries++

					if retries > max_retries {
						t.Logf("failed to generate individual schedule number %d, after %d tries, for %s error %s\n", i, retries, department.Code, err_not_enough_resource.Error())
						err_list_generation = append(err_list_generation, err_not_enough_resource)
						continue new_population_loop
					}

					t.Logf("retry (%d : %s) - %s\n", retries, department.Code, err_not_enough_resource.Error())
					continue
				}

				track_schedules = output_schedules
				track_resources = output_resources
				break // department schedule generated - end retry loop
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
				err_list_validation = append(err_list_validation, e)
			}

			///////////////////////

			// check generated schedule's encoding resource by generating from the previous generated uni time table and comparing
			// it to the resulting encoding resource from the same previous generated department schedule uni time table

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

				// test horizontal validation for the whole university schedule - there should be an error

				fmt.Printf("Generated schedules for all departments, the department %s\n", department.Name)

				err_intentional_horizontal_validation := track_schedules.HorizontalValidation(&persistence, nil, target_semester)

				if err_intentional_horizontal_validation == nil {
					t.Fatal("there should be a missing subject error here since the university schedule is not complete yet")
				}

				// test horizontal validation for the department specific schedule - there should be NO error

				department_to_validate := make(map[uint16]bool)
				department_to_validate[department.DepartmentID] = true

				err_department_horizontal_validations := track_schedules.HorizontalValidation(&persistence, department_to_validate, target_semester)

				for _, e := range err_department_horizontal_validations {
					t.Fatal(e)
				}
			} else {

				// at the very last department, do test horizontal validation for the whole university schedule - there should be NO error

				fmt.Printf("Generated schedules for all departments, the last department schedules generated is %s\n", department.Name)

				if track_schedules.IsEmpty() {
					t.Fatalf("returned an empty university schedule : loop iteration %d\n", i)
				}

				err_horizontal_validations := track_schedules.HorizontalValidation(&persistence, nil, target_semester)

				for _, e := range err_horizontal_validations {
					t.Fatal(e)
				}
			}
		}
	}

	failed_individuals := float64(len(err_list_generation))
	successful_individuals := total_test_iterations - int(failed_individuals)

	if failed_individuals > (float64(total_test_iterations) * allowed_failure_rate) {
		t.Errorf(
			"total of %d fails (%.2f%%) and %d success (%.2f%%) out of the %d populations generated which is above the maximum error treashold of %.2f%%",
			int(failed_individuals), failed_individuals/float64(total_test_iterations)*100.0,
			successful_individuals, float64(successful_individuals)/float64(total_test_iterations)*100.0,
			total_test_iterations, allowed_failure_rate*100.0,
		)
	} else {
		fmt.Printf(
			"total of %d fails (%.2f%%) and %d success (%.2f%%) out of the %d populations generated which is below the maximum error treashold of %.2f%%\n",
			int(failed_individuals), failed_individuals/float64(total_test_iterations)*100.0,
			successful_individuals, float64(successful_individuals)/float64(total_test_iterations)*100.0,
			total_test_iterations, allowed_failure_rate*100.0,
		)
	}

	t.Logf("There are a total of %d validation errors detected when generating university schedules", len(err_list_validation))
}

////////////////////////////////
// BENCHMARK
/////////////////////////////////

func BenchmarkNewPopulationFirstSem(b *testing.B) {
	persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_all_curriculums := persistence.ReaderService.ReadAllCurriculum()

	if err_all_curriculums != nil {
		b.Fatal(err_all_curriculums)
	}

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		b.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	encoding_resource, err_read_default_encoding_resource := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

	if err_read_default_encoding_resource != nil {
		b.Fatal(err_read_default_encoding_resource)
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

	curriculums, err_all_curriculums := persistence.ReaderService.ReadAllCurriculum()

	if err_all_curriculums != nil {
		b.Fatal(err_all_curriculums)
	}

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&persistence)

	if err_dept_id_to_department != nil {
		b.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	encoding_resource, err_read_default_encoding_resource := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

	if err_read_default_encoding_resource != nil {
		b.Fatal(err_read_default_encoding_resource)
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
