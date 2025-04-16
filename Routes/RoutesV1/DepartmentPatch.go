package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
PATCH:

	"/department_update"
*/
func PatchDepartment(ctx *gin.Context) {
	update_department := Departments.Department{}

	if err := ctx.BindJSON(&update_department); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the department updated data")
		return
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.UpdateDepartment(update_department)

	if err != nil {
		log.Print(err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.String(http.StatusOK, "department updated successfully")
}
