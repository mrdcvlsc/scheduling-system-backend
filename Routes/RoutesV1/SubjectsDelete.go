package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
DELETE:

	"/subject_remove?subject_id=[N>0]"
*/
func DeleteSubject(ctx *gin.Context) {
	subject_id, is_valid_subject_id_param := IsValidSubjectID(ctx)

	if !is_valid_subject_id_param {
		return
	}

	{ // check if subject is assign in the first semester subjects
		university_schedules, _ := ObtainUniversityScheduleNoContext(nil, GeneticAlgorithm.TERM_1ST_SEMESTER)

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(GeneticAlgorithm.TERM_1ST_SEMESTER, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_subject_assigned(university_schedules, uint16(subject_id)) {
			ctx.String(http.StatusConflict, "can not delete a subject assigned to a schedule  in 1st semester, clear the schedules first")
			return
		}
	}

	{ // check if subject is assign in the second semester subjects
		university_schedules, _ := ObtainUniversityScheduleNoContext(nil, GeneticAlgorithm.TERM_2ND_SEMESTER)

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(GeneticAlgorithm.TERM_2ND_SEMESTER, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_subject_assigned(university_schedules, uint16(subject_id)) {
			ctx.String(http.StatusConflict, "can not delete a subject assigned to a schedule in 2nd semester, clear the schedules first")
			return
		}
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteSubject(uint16(subject_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the subject")
		return
	}

	ctx.String(http.StatusOK, "subject deleted successfully")
}

func is_subject_assigned(university_schedules Schedule.UniTimeTables, subject_id uint16) bool {
	for _, section_week_schedules := range university_schedules {
		for _, day_time_table := range section_week_schedules {
			for _, time_slot := range day_time_table {
				if subject_id == time_slot.GetSubjectID() {
					return true
				}
			}
		}
	}

	return false
}
