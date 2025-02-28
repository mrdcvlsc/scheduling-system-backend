package RoutesV1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

// requires a 'semester' parameter then validate it.
func IsValidParameterSemesterIndex(ctx *gin.Context) (int, bool) {
	parameter := ctx.Query("semester")

	if parameter == "" {
		ctx.String(http.StatusBadRequest, "missing 'semester' parameter or parameter value")
		return -1, false
	}

	semester, err_atoi := strconv.Atoi(parameter)

	if err_atoi != nil {
		ctx.String(http.StatusBadRequest, "invalid 'semester' parameter value")
		return -1, false
	}

	if semester < 0 || semester >= 2 {
		ctx.String(http.StatusBadRequest, "invalid 'semester' index value")
		return -1, false
	}

	return semester, true
}

func IsValidParameterDepartmentID(ctx *gin.Context) (int, bool) {
	param_department_id := ctx.Query("department_id")

	if param_department_id == "" {
		ctx.String(http.StatusBadRequest, "missing 'department_id' parameter or parameter value")
		return -1, false
	}

	department_id, err_atoi := strconv.Atoi(param_department_id)

	if err_atoi != nil {
		ctx.String(http.StatusBadRequest, "invalid 'department_id' parameter value")
		return -1, false
	}

	if department_id < 0 {
		ctx.String(http.StatusBadRequest, "invalid 'department_id' value")
		return -1, false
	}

	return department_id, true
}

/*
usage inside a route:

	schedule_idx, is_valid_idx := IsValidUniversityScheduleIndex(ctx, university_schedules)
	if !is_valid_idx {
		return
	}
*/
func IsValidUniversityScheduleIndex(ctx *gin.Context, university_schedule Schedule.UniTimeTables) (int, bool) {
	param_schedule_idx := ctx.Query("schedule_idx")

	if param_schedule_idx == "" {
		ctx.String(http.StatusBadRequest, "mising 'schedule_idx' parameter or parameter value")
		return -1, false
	}

	schedule_idx, err_atoi := strconv.Atoi(param_schedule_idx)

	if err_atoi != nil {
		ctx.String(http.StatusBadRequest, "invalid 'schedule_idx' parameter value")
		return -1, false
	}

	if schedule_idx < 0 {
		ctx.String(http.StatusBadRequest, "university schedule underflow: class schedule does not exist")
		return -1, false
	}

	if schedule_idx >= len(university_schedule) {
		ctx.String(http.StatusBadRequest, "university schedule overflow: class schedule does not exist, university schedules might have been altered")
		return -1, false
	}

	return schedule_idx, true
}
