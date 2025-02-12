package RoutesV1

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
retrieves university schedule from cache or persistence.

example usage inside a gin route:

	// setting departments_to_validate to nil will validate all departments.
	university_schedules, is_success := ObtainUniversitySchedule(ctx, nil, selected_semester)
	if !is_success {
		return
	}
*/
func ObtainUniversitySchedule(ctx *gin.Context, departments_to_validate map[uint16]bool, semester int) (Schedule.UniTimeTables, bool) {
	var university_schedules Schedule.UniTimeTables = nil

	cached_university_schedule, has_cache, cache_err := RouteGlobals.GetCachedUniversitySchedule(semester)

	if cache_err != nil {
		log.Println(cache_err.Error())
		ctx.String(http.StatusBadRequest, cache_err.Error())
		return nil, false
	}

	if has_cache {
		log.Println("retrieving university schedule from cache.")
		university_schedules = cached_university_schedule
	} else {
		log.Println("no cached detected loading from persistence")
		read_university_schedules, read_err := RouteGlobals.SchedulePersistence.SaveService.LoadSchedules(semester)

		if read_err != nil {
			log.Println(read_err)

			if errors.Is(read_err, os.ErrNotExist) {
				ctx.String(http.StatusNotFound, "that schedule is not created yet")
				return nil, false
			}

			ctx.String(http.StatusInternalServerError, "we failed to read that schedule")
			return nil, false
		}

		university_schedules = read_university_schedules
	}

	if university_schedules.IsEmpty() {
		log.Println("schedule is empty")
		ctx.String(http.StatusNotFound, "that schedule was empty, please fill it up or generate a schedule again")
		return nil, false
	}

	for _, validation_err := range university_schedules.VerticalValidation(RouteGlobals.ResourcesPersistence) {
		if validation_err != nil {
			log.Println("invalid schedule detected")
			ctx.String(http.StatusConflict, "server detected an invalid schedule with vertically overlapping data")
			return nil, false
		}
	}

	for _, validation_err := range university_schedules.HorizontalValidation(RouteGlobals.ResourcesPersistence, departments_to_validate, semester) {
		if validation_err != nil {
			log.Println("invalid schedule detected")
			ctx.String(http.StatusConflict, "server detected an invalid schedule with wrong horizontal data allocations")
			return nil, false
		}
	}

	log.Println("schedule found")
	return university_schedules, true
}
