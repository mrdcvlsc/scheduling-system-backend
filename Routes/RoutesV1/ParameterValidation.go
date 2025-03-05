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

// requires a 'page_size' parameter then validate it.
func IsValidPageSize(ctx *gin.Context) (int, bool) {
	parameter := ctx.Query("page_size")

	if parameter == "" {
		ctx.String(http.StatusBadRequest, "missing 'page_size' parameter or parameter value")
		return -1, false
	}

	page_size, err_atoi := strconv.Atoi(parameter)

	if err_atoi != nil {
		ctx.String(http.StatusBadRequest, "invalid 'page_size' parameter value")
		return -1, false
	}

	if page_size < 1 {
		ctx.String(http.StatusBadRequest, "invalid 'page_size' index value")
		return -1, false
	}

	return page_size, true
}

// requires a 'page' parameter then validate it.
func IsValidPage(ctx *gin.Context) (int, bool) {
	parameter := ctx.Query("page")

	if parameter == "" {
		ctx.String(http.StatusBadRequest, "missing 'page' parameter or parameter value")
		return -1, false
	}

	page, err_atoi := strconv.Atoi(parameter)

	if err_atoi != nil {
		ctx.String(http.StatusBadRequest, "invalid 'page' parameter value")
		return -1, false
	}

	if page < 0 {
		ctx.String(http.StatusBadRequest, "invalid 'page' index value")
		return -1, false
	}

	return page, true
}

func IsValidInstructorID(ctx *gin.Context) (int, bool) {
	param_instructor_id := ctx.Query("instructor_id")

	if param_instructor_id == "" {
		ctx.String(http.StatusBadRequest, "missing 'instructor_id' parameter or parameter value")
		return -1, false
	}

	instructor_id, err_atoi := strconv.Atoi(param_instructor_id)

	if err_atoi != nil {
		ctx.String(http.StatusBadRequest, "invalid 'instructor_id' parameter value")
		return -1, false
	}

	if instructor_id < 1 {
		ctx.String(http.StatusBadRequest, "invalid 'instructor_id' value")
		return -1, false
	}

	return instructor_id, true
}

func IsValidRoomID(ctx *gin.Context) (int, bool) {
	param_room_id := ctx.Query("room_id")

	if param_room_id == "" {
		ctx.String(http.StatusBadRequest, "missing 'room_id' parameter or parameter value")
		return -1, false
	}

	room_id, err_atoi := strconv.Atoi(param_room_id)

	if err_atoi != nil {
		ctx.String(http.StatusBadRequest, "invalid 'room_id' parameter value")
		return -1, false
	}

	if room_id < 1 {
		ctx.String(http.StatusBadRequest, "invalid 'room_id' value")
		return -1, false
	}

	return room_id, true
}
