package RoutesV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Auth"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
DELETE:

	"/department_remove?department_id=[N>0]"
*/
func DeleteDepartment(ctx *gin.Context) {

	if is_success := Auth.IsAuthSuccess(ctx); !is_success {
		return
	}

	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	if department_id == 0 {
		ctx.String(http.StatusForbidden, "deleting the general department is now allowed")
		return
	}

	if is_allowed := Auth.IsDepartmentAllowed(ctx, uint16(department_id)); !is_allowed {
		return
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteDepartment(uint16(department_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the department")
		return
	}

	ctx.String(http.StatusOK, "department deleted successfully")
}
