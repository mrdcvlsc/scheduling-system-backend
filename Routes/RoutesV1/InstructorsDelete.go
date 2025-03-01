package RoutesV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
POST:

	"/instructor_remove?instructor_id=[N>0]"
*/
func DeleteInstructor(ctx *gin.Context) {
	instructor_id, is_valid_instructor_id_param := IsValidInstructorID(ctx)

	if !is_valid_instructor_id_param {
		return
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteInstructor(uint16(instructor_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the instructor")
		return
	}

	ctx.String(http.StatusOK, "instructor deleted successfully")
}
