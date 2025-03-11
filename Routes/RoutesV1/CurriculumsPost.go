package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
POST:

	"/curriculum_add"
*/
func PostCurriculum(ctx *gin.Context) {
	add_curriculum := Curriculum.Curriculum{}

	if err := ctx.BindJSON(&add_curriculum); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the curriculum to be added")
		return
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.CreateCurriculum(add_curriculum)

	if err != nil {
		log.Print(err)
		ctx.String(http.StatusBadRequest, "we are unable to properly add the curriculum")
		return
	}

	ctx.String(http.StatusOK, "curriculum added successfully")
}
