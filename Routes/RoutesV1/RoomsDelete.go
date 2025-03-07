package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
DELETE:

	"/room_remove?room_id=[N>0]"
*/
func DeleteRoom(ctx *gin.Context) {
	room_id, is_valid_room_id_param := IsValidRoomID(ctx)

	if !is_valid_room_id_param {
		return
	}

	{ // check if room is assign in the first semester subjects
		university_schedules, has_obtained := ObtainUniversitySchedule(ctx, nil, GeneticAlgorithm.TERM_1ST_SEMESTER)

		if !has_obtained {
			return
		}

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(GeneticAlgorithm.TERM_1ST_SEMESTER, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_room_assigned(university_schedules, uint16(room_id)) {
			ctx.String(http.StatusConflict, "can not delete a room assigned to a schedule  in 1st semester")
			return
		}
	}

	{ // check if room is assign in the second semester subjects
		university_schedules, has_obtained := ObtainUniversitySchedule(ctx, nil, GeneticAlgorithm.TERM_2ND_SEMESTER)

		if !has_obtained {
			return
		}

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(GeneticAlgorithm.TERM_2ND_SEMESTER, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_room_assigned(university_schedules, uint16(room_id)) {
			ctx.String(http.StatusConflict, "can not delete a room assigned to a schedule in 2nd semester")
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
