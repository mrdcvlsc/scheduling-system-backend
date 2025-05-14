package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Auth"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
DELETE:

	"/instructor_remove?instructor_id=[N>0]"
*/
func DeleteInstructor(ctx *gin.Context) {

	if is_success := Auth.IsAuthSuccess(ctx); !is_success {
		return
	}

	instructor_id, is_valid_instructor_id_param := IsValidInstructorID(ctx)

	if !is_valid_instructor_id_param {
		return
	}

	{ // check if instructor is assign in the first semester subjects
		university_schedules, _ := ObtainUniversityScheduleNoContext(nil, GeneticAlgorithm.TERM_1ST_SEMESTER)

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(GeneticAlgorithm.TERM_1ST_SEMESTER, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_instructor_assigned(university_schedules, uint16(instructor_id)) {
			ctx.String(http.StatusConflict, "can not delete an instructor assigned to a schedule")
			return
		}
	}

	{ // check if instructor is assign in the second semester subjects
		university_schedules, _ := ObtainUniversityScheduleNoContext(nil, GeneticAlgorithm.TERM_2ND_SEMESTER)

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(GeneticAlgorithm.TERM_2ND_SEMESTER, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_instructor_assigned(university_schedules, uint16(instructor_id)) {
			ctx.String(http.StatusConflict, "can not delete an instructor assigned to a schedule")
			return
		}
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteInstructor(uint16(instructor_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the instructor")
		return
	}

	ctx.String(http.StatusOK, "instructor deleted successfully")
}

func is_instructor_assigned(university_schedules Schedule.UniTimeTables, instructor_id uint16) bool {
	for _, section_week_schedules := range university_schedules {
		for _, day_time_table := range section_week_schedules {
			for _, time_slot := range day_time_table {
				if instructor_id == time_slot.GetInstructorID() {
					return true
				}
			}
		}
	}

	return false
}
