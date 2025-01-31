package RoutesV1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

func GenerateSchedule(ctx *gin.Context) {

	semester, is_valid_semester_idx := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_idx {
		return
	}

	if RouteGlobals.IsGeneratingSchedule.Load() {
		ctx.String(http.StatusAccepted, "already generating schedule")
		return
	}

	ctx.String(http.StatusAccepted, "generating schedule")
	RouteGlobals.IsGeneratingSchedule.Store(true)

	go generate_schedule(semester)
}

// TODO: this is just for testing

func generate_schedule(semester int) {

	defer RouteGlobals.IsGeneratingSchedule.Store(false)

	log.Println("started...")
	// time.Sleep(time.Second * 60)
	// log.Println("wait done, now generating...")

	var generate_university_schedule *Schedule.UniTimeTables

	maximum_trials := 10

	for i := 0; i < maximum_trials; i++ {

		university_schedule, err := GeneticAlgorithm.NewIndividual(RouteGlobals.ResourcesPersistence, semester, 0)

		if (err != nil) && (i == (maximum_trials - 1)) {
			log.Print("generate_schedule:", err.Error())
		}

		if err != nil {
			continue
		}

		generate_university_schedule = &university_schedule
		break
	}

	if generate_university_schedule == nil {
		log.Print("generate_schedule:", fmt.Sprintf("unable to generate schedules after %d tries", maximum_trials))
		return
	}

	if generate_university_schedule.IsEmpty() {
		log.Print("generate_schedule:", "the generated schedule was empty")
	}

	err_validation := generate_university_schedule.Validate(RouteGlobals.ResourcesPersistence)

	for _, e := range err_validation {
		log.Print("generate_schedule:", e.Error())
	}

	write_err := RouteGlobals.SchedulePersistence.WriterService.SaveSchedules(*generate_university_schedule, semester)

	if write_err != nil {
		log.Print("generate_schedule:", write_err.Error())
	}

	cache_err := RouteGlobals.SetCachedUniversitySchedule(semester, *generate_university_schedule)

	if cache_err != nil {
		log.Print("generate_schedule:", cache_err.Error())
	}

	log.Println("ended, schedule was generated")
}
