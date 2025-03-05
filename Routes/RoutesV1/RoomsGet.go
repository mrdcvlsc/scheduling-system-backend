package RoutesV1

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

type RoomTablePage struct {
	Rooms      []Rooms.Room `json:"Rooms"`
	TotalRooms int          `json:"TotalRooms"`
}

/*
GET:

	"/rooms?department_id=D&page_size=[N>0]&page[0-N>0]"
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

	dept_id_to_room_type_to_rooms, err_dept_id_to_room_type_to_rooms := GeneticAlgorithm.GenerateMapDeptIdToRoomTypeToRooms(
		RouteGlobals.ResourcesPersistence,
	)

	if err_dept_id_to_room_type_to_rooms != nil {
		log.Print(err_dept_id_to_room_type_to_rooms)
		ctx.String(http.StatusInternalServerError, "we're currently unable to get the rooms for the requested department")
		return
	}

	department_rooms_separated_by_type, has_department_id := dept_id_to_room_type_to_rooms[uint16(department_id)]

	if !has_department_id {
		ctx.String(http.StatusNotFound, "unable to get rooms for that department because that department does not exist")
		return
	}

	department_rooms := make([]Rooms.Room, 0)

	department_rooms = append(department_rooms, department_rooms_separated_by_type[Rooms.ROOM_TYPE_LEC]...)
	department_rooms = append(department_rooms, department_rooms_separated_by_type[Rooms.ROOM_TYPE_LAB]...)
	department_rooms = append(department_rooms, department_rooms_separated_by_type[Rooms.ROOM_TYPE_GYM]...)

	department_rooms_page := make([]Rooms.Room, 0)

	for i, room := range department_rooms {
		if i < (page_size * page) {
			continue
		}

		department_rooms_page = append(department_rooms_page, room)

		if len(department_rooms_page) >= page_size {
			break
		}
	}

	sort.Slice(department_rooms_page, func(i, j int) bool {
		return department_rooms_page[i].RoomID < department_rooms_page[j].RoomID
	})

	room_table_page := &RoomTablePage{
		Rooms:      department_rooms_page,
		TotalRooms: len(department_rooms),
	}

	ctx.JSON(http.StatusOK, room_table_page)
}
