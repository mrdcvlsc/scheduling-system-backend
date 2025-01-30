package RoutesV1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

func TestRead(ctx *gin.Context) {
	if RouteGlobals.Test == 0 {
		ctx.String(http.StatusNotFound, "uninitialized global variable")
		return
	} else {
		ctx.String(http.StatusOK, fmt.Sprintf("value is now %d", RouteGlobals.Test))
		return
	}
}

func TestWrite(ctx *gin.Context) {
	if ctx.Query("num") == "" {
		ctx.String(http.StatusBadRequest, "malformed request")
	} else {
		num, err := strconv.Atoi(ctx.Query("num"))

		if err != nil {
			ctx.String(http.StatusInternalServerError, "atoi failed")
			return
		}

		RouteGlobals.Test = num

		ctx.String(http.StatusOK, fmt.Sprintf("success in assigning global variable value of %d", num))
		return
	}
}
