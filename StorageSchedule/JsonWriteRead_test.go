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

		sched_gen_err_cnt := 0

		fmt.Println("Reading Resources")

		persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

		////////////////////////////////////////////////////////////////////////////////////////

		curriculums, err_curriculums := persistence.ReaderService.GetAllCurriculum()

		if err_curriculums != nil {
			t.Fatal(err_curriculums)
		}

		////////////////////////////////////////////////////////////////////////////////////////

		encoding_resource, encoding_resource_err := GeneticAlgorithm.ReadDefaultEncodingResource(&persistence)

		if encoding_resource_err != nil {
			t.Fatal(encoding_resource_err)
		}

		////////////////////////////////////////////////////////////////////////////////////////

		fmt.Print("Entering Loop\n")

		for i := 0; i < 40; i++ {
			empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, GeneticAlgorithm.TERM_1ST_SEMESTER)

			university_schedule, _, err := GeneticAlgorithm.EncodeIndividualGenome(
				empty_university_schedule,
				curriculums,
				encoding_resource, nil,
				GeneticAlgorithm.TERM_1ST_SEMESTER, 0,
			)

			if err != nil {
				fmt.Println(err)

				sched_gen_err_cnt++
				continue
			}

			if sched_gen_err_cnt > 28 {
				fmt.Print("Too many errors\n")
				t.Fatal(err)

			}

			populations_of_schedules = append(populations_of_schedules, university_schedule)
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

		vertical_validation_err := first_university_schedule.VerticalValidation(&resource_persistence)

		if first_university_schedule.IsEmpty() {
			t.Fatal("loaded university schedules are empty")
		}

		for e := range vertical_validation_err {
			t.Fatal(e)
		}

		horizontal_validation_err := first_university_schedule.HorizontalValidation(&resource_persistence, GeneticAlgorithm.TERM_1ST_SEMESTER)

		for e := range horizontal_validation_err {
			t.Fatal(e)
		}

		sched_persistence := StorageSchedule.Persistence{WriterService: &StorageSchedule.JsonWriter{}}

		wrote_sched = first_university_schedule
		write_err := sched_persistence.WriterService.SaveSchedules(first_university_schedule, GeneticAlgorithm.TERM_1ST_SEMESTER)

		if write_err != nil {
			t.Fatal(write_err)
		}
	}

	{
		schedule_persistence := StorageSchedule.Persistence{ReaderService: &StorageSchedule.JsonReader{}}
		resource_persistence := StorageResources.Persistence{ReaderService: &StorageResources.JsonReader{}}

		load_university_schedules, load_err := schedule_persistence.ReaderService.LoadSchedules(GeneticAlgorithm.TERM_1ST_SEMESTER)

		if load_err != nil {
			t.Fatal(load_err)
		}

		vertical_validation_err := load_university_schedules.VerticalValidation(&resource_persistence)

		if load_university_schedules.IsEmpty() {
			t.Fatal("loaded university schedules are empty")
		}

		for e := range vertical_validation_err {
			t.Fatal(e)
		}

		if load_err == nil {
			horizontal_validation_err := load_university_schedules.HorizontalValidation(&resource_persistence, GeneticAlgorithm.TERM_1ST_SEMESTER)

			for e := range horizontal_validation_err {
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
