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

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_room_type_to_rooms, err_dept_id_to_room_type_to_rooms := GeneticAlgorithm.GenerateMapDeptIdToRoomTypeToRooms(&storage_persistence)

	if err_dept_id_to_room_type_to_rooms != nil {
		t.Fatal(err_dept_id_to_room_type_to_rooms)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_instructors, err_dept_id_to_instructors := GeneticAlgorithm.GenerateMapDeptIdToInstructors(&storage_persistence)

	if err_dept_id_to_instructors != nil {
		t.Fatal(err_dept_id_to_instructors)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(&storage_persistence)

	if err_dept_id_to_department != nil {
		t.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	{
		uni_sched, err := GeneticAlgorithm.EncodeIndividualGenome(
			curriculums,
			dept_id_to_department,
			dept_id_to_instructors,
			dept_id_to_room_type_to_rooms,
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
			list_of_errors = append(list_of_errors, uni_sched.HorizontalValidation(&resource_persistence, GeneticAlgorithm.TERM_1ST_SEMESTER)...)
		}

		if len(list_of_errors) > 0 {
			for _, e := range list_of_errors {
				t.Error(e)
			}
		}

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
