package RoutesV1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Auth"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
DELETE:

	"/instructor_remove?instructor_id=[N>0]"
*/
func DeleteInstructor(ctx *gin.Context) {

	if is_success := Auth.IsAuthSuccess(ctx); !is_success {
		return
	}

	instructor_id, is_valid_instructor_id_param := IsValidInstructorID(ctx)

	if !is_valid_instructor_id_param {
		return
	}

	selected_instructor, err_read_instructor := RouteGlobals.ResourcesPersistence.ReaderService.ReadInstructor(uint16(instructor_id))

	if err_read_instructor != nil {
		ctx.String(http.StatusInternalServerError, "we're unable to find that instructor right now")
		return
	}

	if is_allowed := Auth.IsDepartmentAllowed(ctx, selected_instructor.DepartmentID); !is_allowed {
		return
	}

	for selected_semester := range Curriculum.SUPPORTED_SEMESTERS {
		university_schedules, _ := ObtainUniversityScheduleNoContext(nil, selected_semester)

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(selected_semester, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_instructor_assigned(university_schedules, uint16(instructor_id)) {
			ctx.String(
				http.StatusConflict,
				fmt.Sprintf("can not delete an instructor assigned to a schedule in %s", Curriculum.SEMESTER_INDEX_NAME[selected_semester]),
			)
			return
		}
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteInstructor(uint16(instructor_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the instructor")
		return
	}

	ctx.String(http.StatusOK, "instructor deleted successfully")
}

func is_instructor_assigned(university_schedules Schedule.UniTimeTables, instructor_id uint16) bool {
	for _, section_week_schedules := range university_schedules {
		for _, day_time_table := range section_week_schedules {
			for _, time_slot := range day_time_table {
				if instructor_id == time_slot.GetInstructorID() {
					return true
				}
			}
		}
	}

	return false
}
