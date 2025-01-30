package RoutesV1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageSchedule"
)

// GET:    /university_schedule?sem=0    or    /university_schedule?sem=1
func GetUniversitySchedule(ctx *gin.Context) {
	param := ctx.Query("sem")

	if param != "0" && param != "1" {
		ctx.String(http.StatusBadRequest, "malformed query")
		return
	}

	semester, parse_err := strconv.Atoi(param)

	if parse_err != nil {
		ctx.String(http.StatusInternalServerError, "Opps! we are not able to parse that semester")
		return
	}

	persistence := StorageSchedule.Persistence{ReaderService: &StorageSchedule.JsonReader{}}

	university_schedules, read_err := persistence.ReaderService.LoadSchedules(semester)

	if read_err != nil {
		ctx.String(http.StatusInternalServerError, "Opps! error reading the schedule")
		return
	}

	if university_schedules.IsEmpty() {
		ctx.String(http.StatusNotFound, "a schedule was not found, please generate one first")
		return
	}

	for _, validation_err := range university_schedules.Validate() {
		if validation_err != nil {
			ctx.String(http.StatusConflict, "we detected an invalid schedule")
			return
		}
	}

	serialized_schedule := Schedule.SerializeUniversitySchedule(university_schedules)

	ctx.Data(http.StatusOK, "application/octet-stream", serialized_schedule)
}
