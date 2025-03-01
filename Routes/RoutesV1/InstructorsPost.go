package RoutesV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
POST:

	"/instructor_add"
*/
func PostInstructor(ctx *gin.Context) {
	add_instructor_with_time_str := Instructors.InstructorWithTimeString{}

	if err := ctx.BindJSON(&add_instructor_with_time_str); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the instructor to be add")
		return
	}

	add_instructor := Instructors.Instructor{
		InstructorID:  add_instructor_with_time_str.InstructorID,
		DepartmentID:  add_instructor_with_time_str.DepartmentID,
		FirstName:     add_instructor_with_time_str.FirstName,
		MiddleInitial: add_instructor_with_time_str.MiddleInitial,
		LastName:      add_instructor_with_time_str.LastName,
	}

	add_instructor.Time.StringParse(add_instructor_with_time_str.Time)

	err := RouteGlobals.ResourcesPersistence.WriterService.CreateInstructor(add_instructor)

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly add the instructor")
		return
	}

	ctx.String(http.StatusOK, "instructor added successfully")
}
