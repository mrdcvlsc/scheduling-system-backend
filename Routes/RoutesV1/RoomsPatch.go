package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
PATCH:

	"/room_update"
*/
func PatchRoom(ctx *gin.Context) {
	update_room := Rooms.Room{}

	if err := ctx.BindJSON(&update_room); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the room updated data")
		return
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.UpdateRoom(update_room)

	if err != nil {
		log.Print(err)
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.String(http.StatusOK, "room updated successfully")
}
