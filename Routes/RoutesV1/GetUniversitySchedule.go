package RoutesV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

// GET:  /v1/university_schedule?semester=0  or  /v1/university_schedule?semester=1
func GetUniversitySchedule(ctx *gin.Context) {

	selected_semester, is_valid_semester_para := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_para {
		return
	}

	university_schedules, has_obtained := ObtainUniversitySchedule(ctx, selected_semester)

	if !has_obtained {
		return
	}

	serialized_schedule := Schedule.SerializeUniversitySchedule(university_schedules)

	ctx.Data(http.StatusOK, "application/octet-stream", serialized_schedule)
}
