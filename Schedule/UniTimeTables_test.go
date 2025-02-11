package Schedule_test

import (
	"os"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func Test_UniTimeTablesSerializationAndDeserialization(t *testing.T) {

	var wrote_sched Schedule.UniTimeTables

	storage_persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}
	resource_persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_curriculums := storage_persistence.ReaderService.GetAllCurriculum()

	if err_curriculums != nil {
		t.Fatal(err_curriculums)
	}

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&resource_persistence)

	if err_dept_id_to_department != nil {
		t.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	encoding_resource, encoding_resource_err := GeneticAlgorithm.ReadDefaultEncodingResource(&resource_persistence)

	if encoding_resource_err != nil {
		t.Fatal(encoding_resource_err)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	{
		empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, GeneticAlgorithm.TERM_1ST_SEMESTER)

		uni_sched, encoding_resource_output, err := GeneticAlgorithm.EncodeIndividualGenome(
			empty_university_schedule,
			curriculums, dept_id_to_department,
			encoding_resource, nil,
			GeneticAlgorithm.TERM_1ST_SEMESTER, 0,
		)

		if uni_sched == nil && err != nil {
			t.Fatalf("there was an error reading data: %v\n", err)
		}

		if uni_sched.IsEmpty() {
			t.Fatal("there was no schedule generated to be tested")
		}

		list_of_errors := make([]error, 0)

		list_of_errors = append(list_of_errors, uni_sched.VerticalValidation(&resource_persistence)...)

		if len(list_of_errors) > 0 {
			for _, e := range list_of_errors {
				t.Error(e)
			}
		}

		if err == nil {
			list_of_errors = append(list_of_errors, uni_sched.HorizontalValidation(&resource_persistence, nil, GeneticAlgorithm.TERM_1ST_SEMESTER)...)
		}

		if len(list_of_errors) > 0 {
			for _, e := range list_of_errors {
				t.Error(e)
			}
		}

		///////////////////////

		generated_encoding_resource, output_resources := GeneticAlgorithm.GenerateEncodingResourceFromUniTimeTable(
			uni_sched, curriculums, GeneticAlgorithm.TERM_1ST_SEMESTER, &resource_persistence,
		)

		if output_resources != nil {
			t.Fatal(output_resources)
		}

		if encoding_resource_output != nil {
			if !GeneticAlgorithm.IsEqualEncodingResource(generated_encoding_resource, encoding_resource_output) {
				t.Fatal("generated encoding resource from bare university schedule is not equal to the produced encoding resource of GA")
			}
		}

		/////////////////

		serialized_data := Schedule.SerializeUniversitySchedule(uni_sched)

		if err := Utils.SaveToBinFile("tmp-university.schedule", serialized_data); err != nil {
			t.Fatalf("failed to save serialized data to binary file: %v\n", err)
		}

		wrote_sched = uni_sched
	}

	{
		read_bytes, err := Utils.ReadFromBinFile("tmp-university.schedule")

		if err != nil {
			t.Fatalf("failed to read the serialized data from the binary file: %v\n", err)
		}

		deserialize_data := Schedule.DeserializeUniversitySchedule(read_bytes)

		for section_idx, original := range deserialize_data {
			for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
				for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
					if wrote_sched[section_idx][day][time_slot] != original[day][time_slot] {
						t.Errorf("Data missmatch at section index %d, time slot (day: %d, time_slot:%d)\n", section_idx, day, time_slot)
					}
				}
			}
		}
	}

	err := os.Remove("tmp-university.schedule")

	if err != nil {
		t.Fatalf("error deleting file: %v\n", err)
	}
}
