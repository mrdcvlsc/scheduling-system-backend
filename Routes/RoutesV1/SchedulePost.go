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

/*
POST:

	"/generate_schedule?semester[0-1]"
*/
func GenerateSchedule(ctx *gin.Context) {

	semester, is_valid_semester_idx := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_idx {
		return
	}

	if RouteGlobals.IsGeneratingSchedule.Load() {
		log.Print("GenerateSchedule: already generating schedule")
		ctx.String(http.StatusAccepted, "already generating schedule")
		return
	}

	log.Print("GenerateSchedule: generating schedule")
	ctx.String(http.StatusAccepted, "generating schedule")

	go generate_schedule(semester)
}

// TODO: this is just for testing

func generate_schedule(semester int) {
	RouteGlobals.IsGeneratingSchedule.Store(true)
	defer RouteGlobals.IsGeneratingSchedule.Store(false)

	var generate_university_schedule *Schedule.UniTimeTables

	maximum_trials := 10

	log.Println("generate_schedule: generating schedule...")

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_curriculums != nil {
		log.Fatal("generate_schedule:", err_curriculums)
	}

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(RouteGlobals.ResourcesPersistence)

	if err_dept_id_to_department != nil {
		log.Fatal(err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	encoding_resource, err_read_default_encoding_resource := GeneticAlgorithm.ReadDefaultEncodingResource(RouteGlobals.ResourcesPersistence)

	if err_read_default_encoding_resource != nil {
		log.Fatal(err_read_default_encoding_resource)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	for i := range maximum_trials {
		empty_university_schedule := GeneticAlgorithm.NewEmptyIndividual(curriculums, semester)

		university_schedule, _, err := GeneticAlgorithm.EncodeIndividualGenome(
			empty_university_schedule,
			curriculums, dept_id_to_department,
			encoding_resource, nil,
			semester, 0,
		)

		if (err != nil) && (i == (maximum_trials - 1)) {
			log.Print("generate_schedule:", err.Error())
		}

		if err != nil {
			continue
		}

		generate_university_schedule = &university_schedule
		break
	}

	log.Println("generate_schedule: done generating...")

	if generate_university_schedule == nil {
		log.Print("generate_schedule:", fmt.Sprintf("unable to generate schedules after %d tries", maximum_trials))
		return
	}

	if generate_university_schedule.IsEmpty() {
		log.Print("generate_schedule:", "the generated schedule was empty")
	}

	log.Println("generate_schedule: validating schedule")

	for _, e := range generate_university_schedule.VerticalValidation(RouteGlobals.ResourcesPersistence) {
		log.Print("generate_schedule:", e.Error())
	}

	for _, e := range generate_university_schedule.HorizontalValidation(RouteGlobals.ResourcesPersistence, nil, semester) {
		log.Print("generate_schedule:", e.Error())
	}

	log.Println("generate_schedule: saving schedule")

	err_save_schedules := RouteGlobals.SchedulePersistence.SaveService.SaveSchedules(*generate_university_schedule, semester)

	if err_save_schedules != nil {
		log.Print("generate_schedule:", err_save_schedules.Error())
	}

	log.Println("generate_schedule: caching schedule")

	err_set_cache := RouteGlobals.SetCachedUniversitySchedule(semester, *generate_university_schedule)

	if err_set_cache != nil {
		log.Print("generate_schedule:", err_set_cache.Error())
	}

	log.Println("generate_schedule: ended, schedule was generated")
}
