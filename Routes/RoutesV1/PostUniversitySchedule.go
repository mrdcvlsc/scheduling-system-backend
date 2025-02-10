package RoutesV1

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
POST:

	"/university_schedule"
*/
func PostUniversitySchedule(ctx *gin.Context) {

	////////////////////////////////////////////////////////////////////////////////////////
	//                                   FOR TESTING
	////////////////////////////////////////////////////////////////////////////////////////

	selected_semester, is_valid_semester_para := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_para {
		return
	}

	university_schedules, read_err := RouteGlobals.SchedulePersistence.ReaderService.LoadSchedules(selected_semester)

	if read_err != nil {
		ctx.String(http.StatusInternalServerError, "error reading the schedule")
		return
	}

	if university_schedules.IsEmpty() {
		ctx.String(http.StatusNotFound, "schedule was empty, please generate one first")
		return
	}

	for _, validation_err := range university_schedules.VerticalValidation(RouteGlobals.ResourcesPersistence) {
		if validation_err != nil {
			ctx.String(http.StatusConflict, "we detected an invalid schedule")
			return
		}
	}

	for _, validation_err := range university_schedules.HorizontalValidation(RouteGlobals.ResourcesPersistence, selected_semester) {
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
					ctx.String(http.StatusConflict, "schedule instructor id mismatch")
					return
				}

				if university_schedules[section_idx][day][time_slot].GetRoomID() != received_university_schedules[section_idx][day][time_slot].GetRoomID() {
					ctx.String(http.StatusConflict, "schedule room id mismatch")
					return
				}
			}
		}
	}

	ctx.String(http.StatusOK, "schedule successfully matched")
}
