package RoutesV1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Auth"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
DELETE:

	"/room_remove?room_id=[N>0]"
*/
func DeleteRoom(ctx *gin.Context) {

	if is_success := Auth.IsAuthSuccess(ctx); !is_success {
		return
	}

	room_id, is_valid_room_id_param := IsValidRoomID(ctx)

	if !is_valid_room_id_param {
		return
	}

	selected_room, err_read_room := RouteGlobals.ResourcesPersistence.ReaderService.ReadRoom(uint16(room_id))

	if err_read_room != nil {
		ctx.String(http.StatusInternalServerError, "we're unable to find that room right now")
		return
	}

	if is_allowed := Auth.IsDepartmentAllowed(ctx, selected_room.DepartmentID); !is_allowed {
		return
	}

	for selected_semester := range Curriculum.SUPPORTED_SEMESTERS {
		university_schedules, _ := ObtainUniversityScheduleNoContext(nil, selected_semester)

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(selected_semester, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_room_assigned(university_schedules, uint16(room_id)) {
			ctx.String(
				http.StatusConflict,
				fmt.Sprintf("can not delete a room assigned to a schedule in %s", Curriculum.SEMESTER_INDEX_NAME[selected_semester]),
			)
			return
		}
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteRoom(uint16(room_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the room")
		return
	}

	ctx.String(http.StatusOK, "room deleted successfully")
}

func is_room_assigned(university_schedules Schedule.UniTimeTables, room_id uint16) bool {
	for _, section_week_schedules := range university_schedules {
		for _, day_time_table := range section_week_schedules {
			for _, time_slot := range day_time_table {
				if room_id == time_slot.GetRoomID() {
					return true
				}
			}
		}
	}

	return false
}
