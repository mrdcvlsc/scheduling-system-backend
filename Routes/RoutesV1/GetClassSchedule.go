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

	semester_param := ctx.Query("semester")

	if semester_param == "" {
		ctx.String(http.StatusBadRequest, "mising 'semester_param' parameter or parameter value")
		return
	}

	semester, semester_atoi_err := strconv.Atoi(semester_param)

	if semester_atoi_err != nil {
		ctx.String(http.StatusInternalServerError, "invalid 'semester' value")
		return
	}

	if semester != 0 && semester != 1 {
		ctx.String(http.StatusBadRequest, "invalid 'semester' index")
		return
	}

	// parse department_id parameter

	param_department_id := ctx.Query("department_id")

	if param_department_id == "" {
		ctx.String(http.StatusBadRequest, "mising 'department_id' parameter or parameter value")
		return
	}

	department_id, department_id_atoi_err := strconv.Atoi(param_department_id)

	if department_id_atoi_err != nil {
		ctx.String(http.StatusBadRequest, "invalid 'department_id' parameter value")
		return
	}

	if department_id <= 0 {
		ctx.String(http.StatusBadRequest, "invalid 'department_id' value")
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
