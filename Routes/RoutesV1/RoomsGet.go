package RoutesV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

type RoomTablePage struct {
	Rooms      []Rooms.Room `json:"Rooms"`
	TotalRooms int          `json:"TotalRooms"`
}

/*
GET:

	"/rooms?department_id=D&page_size=[N>0]&page[0-N>0]&name_match=[string]"
*/
func GetDepartmentRooms(ctx *gin.Context) {
	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)
	if !is_valid_department_id_param {
		return
	}

	page_size, is_valid_page_size_param := IsValidPageSize(ctx)
	if !is_valid_page_size_param {
		return
	}

	page, is_valid_page_param := IsValidPage(ctx)
	if !is_valid_page_param {
		return
	}

	name_parameter := ctx.Query("name_match")

	all_rooms, err_read_all_rooms := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllRooms()

	if err_read_all_rooms != nil {
		ctx.String(http.StatusInternalServerError, "we're currently unable to get the rooms")
		return
	}

	department_rooms_page := make([]Rooms.Room, 0)
	total_department_rooms := 0

	for _, room := range all_rooms {
		if room.DepartmentID != uint16(department_id) {
			continue
		}

		if len(name_parameter) > 0 {
			if !Utils.HasSubString(room.Name, name_parameter) {
				continue
			}
		}

		total_department_rooms++

		if (total_department_rooms - 1) < (page_size * page) {
			continue
		}

		if len(department_rooms_page) < page_size {
			department_rooms_page = append(department_rooms_page, room)
		}
	}

	room_table_page := &RoomTablePage{
		Rooms:      department_rooms_page,
		TotalRooms: total_department_rooms,
	}

	ctx.JSON(http.StatusOK, room_table_page)
}
