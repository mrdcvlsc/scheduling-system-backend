package RoutesV1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// requires a 'semester' parameter then validate it.
func IsValidParameterSemesterIndex(ctx *gin.Context) (int, bool) {
	parameter := ctx.Query("semester")

	if parameter == "" {
		ctx.String(http.StatusBadRequest, "missing 'semester' parameter or parameter value")
		return -1, false
	}

	semester, semester_atoi_err := strconv.Atoi(parameter)

	if semester_atoi_err != nil {
		ctx.String(http.StatusBadRequest, "invalid 'semester' parameter value")
		return -1, false
	}

	if semester < 0 || semester >= 2 {
		ctx.String(http.StatusBadRequest, "invalid 'semester' index value")
		return -1, false
	}

	return semester, true
}

func IsValidParameterDepartmentID(ctx *gin.Context) (int, bool) {
	param_department_id := ctx.Query("department_id")

	if param_department_id == "" {
		ctx.String(http.StatusBadRequest, "missing 'department_id' parameter or parameter value")
		return -1, false
	}

	department_id, department_id_atoi_err := strconv.Atoi(param_department_id)

	if department_id_atoi_err != nil {
		ctx.String(http.StatusBadRequest, "invalid 'department_id' parameter value")
		return -1, false
	}

	if department_id <= 0 {
		ctx.String(http.StatusBadRequest, "invalid 'department_id' value")
		return -1, false
	}

	return department_id, true
}
