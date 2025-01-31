package RoutesV1

import (
	"net/http"
	"strconv"

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

	// TODO: use this for authentication later on.

	_, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	// load university schedules

	university_schedules, read_err := RouteGlobals.SchedulePersistence.ReaderService.LoadSchedules(semester)

	if read_err != nil {
		ctx.String(http.StatusInternalServerError, "error reading schedule")
		return
	}

	if university_schedules.IsEmpty() {
		ctx.String(http.StatusNotFound, "schedule was not found, please generate one first")
		return
	}

	for _, validation_err := range university_schedules.Validate() {
		if validation_err != nil {
			ctx.String(http.StatusConflict, "server detected an invalid schedule")
			return
		}
	}

	// parse schedule_idx parameter

	param_schedule_idx := ctx.Query("schedule_idx")

	if param_schedule_idx == "" {
		ctx.String(http.StatusBadRequest, "mising 'schedule_idx' parameter or parameter value")
		return
	}

	schedule_idx, schedule_idx_atoi_err := strconv.Atoi(param_schedule_idx)

	if schedule_idx_atoi_err != nil {
		ctx.String(http.StatusBadRequest, "invalid 'schedule_idx' parameter value")
		return
	}

	if schedule_idx < 0 {
		ctx.String(http.StatusBadRequest, "university schedule underflow: class schedule does not exist")
		return
	}

	if schedule_idx >= len(university_schedules) {
		ctx.String(http.StatusBadRequest, "university schedule overflow: class schedule does not exist, university schedules might have been altered")
		return
	}

	// extract selected schedule

	serialized_schedule := Schedule.SerializeUniversitySchedule(university_schedules[schedule_idx:(schedule_idx + 1)])

	ctx.Data(http.StatusOK, "application/octet-stream", serialized_schedule)
}
