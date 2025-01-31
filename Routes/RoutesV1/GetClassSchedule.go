package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

// GET : /v1/class_schedule?department_id=D&semester=S&schedule_idx=I
func GetClassSchedule(ctx *gin.Context) {

	// parse semester parameter

	semester, is_valid_semester_param := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	// parse department_id parameter

	// TODO: use department_id for authentication later on.

	_, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	// load university schedules

	university_schedules, has_obtained := ObtainUniversitySchedule(ctx, semester)

	if !has_obtained {
		return
	}

	// cache the found university schedule for the semester

	cache_err := RouteGlobals.SetCachedUniversitySchedule(semester, university_schedules)

	if cache_err != nil {
		log.Println(cache_err.Error())
	}

	// parse schedule_idx parameter

	schedule_idx, is_valid_idx := IsValidUniversityScheduleIndex(ctx, university_schedules)

	if !is_valid_idx {
		return
	}

	// extract selected schedule

	serialized_schedule := Schedule.SerializeUniversitySchedule(university_schedules[schedule_idx:(schedule_idx + 1)])

	ctx.Data(http.StatusOK, "application/octet-stream", serialized_schedule)
}
