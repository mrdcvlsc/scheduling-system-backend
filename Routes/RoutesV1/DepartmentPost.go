package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
POST:

	"/department_add"
*/
func PostDepartment(ctx *gin.Context) {
	add_department := Departments.Department{}

	if err := ctx.BindJSON(&add_department); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the department to be added")
		return
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.CreateDepartment(add_department)

	if err != nil {
		log.Print(err)
		ctx.String(http.StatusBadRequest, "we are unable to properly add the department")
		return
	}

	ctx.String(http.StatusOK, "department added successfully")
}
