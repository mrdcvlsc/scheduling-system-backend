package GeneticAlgorithm

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
)

const DAY_SWAP_PERCENT_PROBABILITY int = 10           // %
const SECTION_WEEK_CLEAR_PERCENT_PROBABILITY int = 10 // %

func ApplyRandomMutation(
	sched Schedule.UniTimeTables,
	curriculums []Curriculum.Curriculum,
	semester int,
	resource_persistence *StorageResources.Persistence,
) (Schedule.UniTimeTables, *EncodingResource, error) {

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	// TODO: random subject time slot nudge

	// TODO: random subject time slot day swap

	// day swaps

	for usi := range len(sched) {
		for day := range Const.N_WEEKLY_SCHOOL_DAYS {

			if rng.Int31n(100) >= int32(DAY_SWAP_PERCENT_PROBABILITY) {
				continue
			}

			log.Print("random day swap mutation")

			day_swap := rng.Intn(int(Const.N_WEEKLY_SCHOOL_DAYS))

			if day_swap == day {
				continue
			}

			for time_slot := range Const.N_DAILY_TIME_SLOTS {
				sched[usi][day][time_slot], sched[usi][day_swap][time_slot] = sched[usi][day_swap][time_slot], sched[usi][day][time_slot]
			}
		}
	}

	// random section clear

	for usi := range len(sched) {
		if rng.Int31n(100) >= int32(SECTION_WEEK_CLEAR_PERCENT_PROBABILITY) {
			continue
		}

		sched[usi] = Schedule.WeekTimeTable{}
	}

	// TODO: random subject clear
	// TODO: then generate/complete encoding

	// can be an argument
	generated_encoding_resource, err_gen_encoding_resource := GenerateEncodingResourceFromUniTimeTable(
		sched, curriculums, semester, resource_persistence,
	)
	if err_gen_encoding_resource != nil {
		return sched, generated_encoding_resource, err_gen_encoding_resource
	}

	// can be an argument
	dept_id_to_department, err_dept_id_to_department := GenerateMapDeptIdToDepartment(resource_persistence)
	if err_dept_id_to_department != nil {
		return sched, generated_encoding_resource, err_dept_id_to_department
	}

	university_schedules, encoding_resource, err := EncodeIndividualGenome(
		sched,
		curriculums,
		dept_id_to_department, generated_encoding_resource, nil,
		semester, 0,
	)

	if len(university_schedules) == 0 {
		return university_schedules, encoding_resource, errors.New("there are no university schedules")
	}

	if err != nil {
		return university_schedules, encoding_resource, err
	}

	return university_schedules, encoding_resource, nil
}
