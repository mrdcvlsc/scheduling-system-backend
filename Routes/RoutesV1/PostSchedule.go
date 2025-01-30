package RoutesV1

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageSchedule"
)

func PostSchedule(ctx *gin.Context) {

	////////////////////////////////////////////////////////////////////////////////////////
	//                                   FOR TESTING
	////////////////////////////////////////////////////////////////////////////////////////

	param := ctx.Query("sem")

	if param != "0" && param != "1" {
		ctx.String(http.StatusBadRequest, "malformed query")
		return
	}

	semester, parse_err := strconv.Atoi(param)

	if parse_err != nil {
		ctx.String(http.StatusInternalServerError, "opps! we are not able to parse that semester")
		return
	}

	persistence := StorageSchedule.Persistence{ReaderService: &StorageSchedule.JsonReader{}}

	university_schedules, read_err := persistence.ReaderService.LoadSchedules(semester)

	if read_err != nil {
		ctx.String(http.StatusInternalServerError, "error reading the schedule")
		return
	}

	if university_schedules.IsEmpty() {
		ctx.String(http.StatusNotFound, "schedule was not found, please generate one first")
		return
	}

	for _, validation_err := range university_schedules.Validate() {
		if validation_err != nil {
			ctx.String(http.StatusConflict, "we detected an invalid schedule")
			return
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	serialized_data, body_err := io.ReadAll(ctx.Request.Body)

	if body_err != nil {
		ctx.String(http.StatusBadRequest, "we failed to read that post request body")
		return
	}

	received_university_schedules := Schedule.DeserializeUniversitySchedule(serialized_data)

	if len(university_schedules) != len(received_university_schedules) {
		ctx.String(http.StatusConflict, "schedule length mismatch")
		return
	}

	for section_idx := 0; section_idx < len(university_schedules); section_idx++ {
		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				if university_schedules[section_idx][day][time_slot].GetSubjectID() != received_university_schedules[section_idx][day][time_slot].GetSubjectID() {
					ctx.String(http.StatusConflict, "schedule subject id mismatch")
					return
				}

				if university_schedules[section_idx][day][time_slot].GetInstructorID() != received_university_schedules[section_idx][day][time_slot].GetInstructorID() {
					ctx.String(http.StatusConflict, "schedule subject id mismatch")
					return
				}

				if university_schedules[section_idx][day][time_slot].GetRoomID() != received_university_schedules[section_idx][day][time_slot].GetRoomID() {
					ctx.String(http.StatusConflict, "schedule subject id mismatch")
					return
				}
			}
		}
	}

	ctx.String(http.StatusOK, "schedule successfully matched")
}
