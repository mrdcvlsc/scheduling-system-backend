package RoutesV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Auth"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
PATCH:

	"/instructor_update"
*/
func PatchInstructor(ctx *gin.Context) {

	if is_success := Auth.IsAuthSuccess(ctx); !is_success {
		return
	}

	update_instructor_with_time_str := Instructors.InstructorWithTimeString{}

	if err := ctx.BindJSON(&update_instructor_with_time_str); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the instructor updated data")
		return
	}

	selected_instructor, err_read_instructor := RouteGlobals.ResourcesPersistence.ReaderService.ReadInstructor(update_instructor_with_time_str.InstructorID)

	if err_read_instructor != nil {
		ctx.String(http.StatusInternalServerError, "we're unable to find that instructor right now")
		return
	}

	if is_allowed := Auth.IsDepartmentAllowed(ctx, selected_instructor.DepartmentID); !is_allowed {
		return
	}

	update_instructor := Instructors.Instructor{
		InstructorID:  update_instructor_with_time_str.InstructorID,
		DepartmentID:  update_instructor_with_time_str.DepartmentID,
		FirstName:     update_instructor_with_time_str.FirstName,
		MiddleInitial: update_instructor_with_time_str.MiddleInitial,
		LastName:      update_instructor_with_time_str.LastName,
	}

	update_instructor.Time.StringParse(update_instructor_with_time_str.Time)

	err := RouteGlobals.ResourcesPersistence.WriterService.UpdateInstructor(update_instructor)

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly updated the instructor")
		return
	}

	ctx.String(http.StatusOK, "instructor updated successfully")
}
