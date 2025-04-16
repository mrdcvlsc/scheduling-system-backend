package StorageSchedule_test

import (
	"fmt"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageSchedule"
)

func Test_JsonReadWriteUniversitySchedules(t *testing.T) {

	var wrote_sched Schedule.UniTimeTables
	resource_persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

	{
		fmt.Print("Make Population Slice\n")

		populations_of_schedules := make([]Schedule.UniTimeTables, 0)

		fmt.Print("initialize error counter\n")

		fmt.Println("Reading Resources")

		persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

		////////////////////////////////////////////////////////////////////////////////////////

		curriculums, err_curriculums := persistence.ReaderService.ReadAllCurriculum()

		if err_curriculums != nil {
			t.Fatal(err_curriculums)
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

		fmt.Print("Entering Loop\n")

		tries := 128
		for range tries {

			empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, GeneticAlgorithm.TERM_1ST_SEMESTER)

			university_schedule, _, err := GeneticAlgorithm.EncodeIndividualGenome(
				empty_university_schedule,
				curriculums, dept_id_to_department,
				encoding_resource, nil,
				GeneticAlgorithm.TERM_1ST_SEMESTER, 0,
			)

			if err != nil {
				t.Logf("retrying due to error : %s", err.Error())
				continue
			}

			populations_of_schedules = append(populations_of_schedules, university_schedule)
			break
		}

		if len(populations_of_schedules) == 0 {
			t.Fatal("No university schedules generated")
		}

		first_university_schedule := populations_of_schedules[0]

		if first_university_schedule.IsEmpty() {
			t.Fatal("the first university schedule generated is empty")
		}

		err_validation := first_university_schedule.VerticalValidation(&resource_persistence)

		for _, e := range err_validation {
			t.Fatal(e)
		}

		err_vertical_validation := first_university_schedule.VerticalValidation(&resource_persistence)

		if first_university_schedule.IsEmpty() {
			t.Fatal("loaded university schedules are empty")
		}

		for e := range err_vertical_validation {
			t.Fatal(e)
		}

		err_horizontal_validation := first_university_schedule.HorizontalValidation(&resource_persistence, nil, GeneticAlgorithm.TERM_1ST_SEMESTER)

		for e := range err_horizontal_validation {
			t.Fatal(e)
		}

		sched_persistence := StorageSchedule.Persistence{SaveService: &StorageSchedule.JsonWriter{}}

		wrote_sched = first_university_schedule
		err_save_schedules := sched_persistence.SaveService.SaveSchedules(first_university_schedule, GeneticAlgorithm.TERM_1ST_SEMESTER)

		if err_save_schedules != nil {
			t.Fatal(err_save_schedules)
		}
	}

	{
		schedule_persistence := StorageSchedule.Persistence{LoadService: &StorageSchedule.JsonReader{}}
		resource_persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

		load_university_schedules, err_load_schedules := schedule_persistence.LoadService.LoadSchedules(GeneticAlgorithm.TERM_1ST_SEMESTER)

		if err_load_schedules != nil {
			t.Fatal(err_load_schedules)
		}

		err_vertical_validation := load_university_schedules.VerticalValidation(&resource_persistence)

		if load_university_schedules.IsEmpty() {
			t.Fatal("loaded university schedules are empty")
		}

		for e := range err_vertical_validation {
			t.Fatal(e)
		}

		if err_load_schedules == nil {
			err_horizontal_validation := load_university_schedules.HorizontalValidation(&resource_persistence, nil, GeneticAlgorithm.TERM_1ST_SEMESTER)

			for e := range err_horizontal_validation {
				t.Fatal(e)
			}
		}

		for i, week_time_tables := range load_university_schedules {
			for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
				for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
					if week_time_tables[day][time_slot].GetSubjectID() != wrote_sched[i][day][time_slot].GetSubjectID() {
						t.Errorf(
							"Loaded university schedule subject mismatch at section_idx: %d, day: %d, time_slot: %d",
							i, day, time_slot,
						)
					}

					if week_time_tables[day][time_slot].GetInstructorID() != wrote_sched[i][day][time_slot].GetInstructorID() {
						t.Errorf(
							"Loaded university schedule instructor mismatch at section_idx: %d, day: %d, time_slot: %d",
							i, day, time_slot,
						)
					}

					if week_time_tables[day][time_slot].GetRoomID() != wrote_sched[i][day][time_slot].GetRoomID() {
						t.Errorf(
							"Loaded university schedule room mismatch at section_idx: %d, day: %d, time_slot: %d",
							i, day, time_slot,
						)
					}
				}
			}
		}
	}
}
