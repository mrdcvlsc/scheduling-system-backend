package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
PATCH:

	"/curriculum_update"
*/
func PatchCurriculum(ctx *gin.Context) {
	update_curriculum := Curriculum.Curriculum{}

	if err := ctx.BindJSON(&update_curriculum); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the curriculum updated data")
		return
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.UpdateCurriculum(update_curriculum)

	if err != nil {
		log.Print(err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.String(http.StatusOK, "curriculum updated successfully")
}
