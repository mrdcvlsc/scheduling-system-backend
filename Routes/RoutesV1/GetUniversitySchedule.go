package RoutesV1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

// GET:  /v1/university_schedule?semester=0  or  /v1/university_schedule?semester=1
func GetUniversitySchedule(ctx *gin.Context) {

	param_semester := ctx.Query("semester")

	if param_semester == "" {
		ctx.String(http.StatusBadRequest, "mising 'semester' parameter or parameter value")
		return
	}

	selected_semester, semester_atoi_err := strconv.Atoi(param_semester)

	if semester_atoi_err != nil {
		ctx.String(http.StatusBadRequest, "invalid 'semester' parameter value")
		return
	}

	if selected_semester < 0 || selected_semester >= 2 {
		ctx.String(http.StatusBadRequest, "invalid 'semester' index value")
		return
	}

	university_schedules, read_err := RouteGlobals.SchedulePersistence.ReaderService.LoadSchedules(selected_semester)

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

	serialized_schedule := Schedule.SerializeUniversitySchedule(university_schedules)

	ctx.Data(http.StatusOK, "application/octet-stream", serialized_schedule)
}
