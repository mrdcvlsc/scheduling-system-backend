package RoutesV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
DELETE:

	"/curriculum_remove?curriculum_id=[N>0]"
*/
func DeleteCurriculum(ctx *gin.Context) {
	curriculum_id, is_valid_curriculum_id_param := IsValidCurriculumID(ctx)

	if !is_valid_curriculum_id_param {
		return
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteCurriculum(uint16(curriculum_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the curriculum")
		return
	}

	ctx.String(http.StatusOK, "curriculum deleted successfully")
}
