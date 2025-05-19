package RoutesV1

import (
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

	"/subject_remove?subject_id=[N>0]"
*/
func DeleteSubject(ctx *gin.Context) {

	if is_success := Auth.IsAuthSuccess(ctx); !is_success {
		return
	}

	subject_id, is_valid_subject_id_param := IsValidSubjectID(ctx)

	if !is_valid_subject_id_param {
		return
	}

	if RouteGlobals.IsGeneratingSchedule.Load() {
		log.Print("DeleteSubject: [busy] you or other department(s) are still generating a schedule, please wait until the process is finished")
		ctx.String(http.StatusForbidden, "we're unable to delete the subject right now, you or other department(s) are still generating a schedule, please wait a little while until those process are done")
		return
	}

	RouteGlobals.ReindexUniSchedMutex.Lock()
	defer RouteGlobals.ReindexUniSchedMutex.Unlock()

	for semester := range Curriculum.SUPPORTED_SEMESTERS {
		university_schedules, _ := ObtainUniversityScheduleNoContext(nil, semester)

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(semester, university_schedules)

		if err_set_cache != nil {
			log.Println(err_set_cache.Error())
		}

		if is_subject_assigned(university_schedules, uint16(subject_id)) {
			ctx.String(
				http.StatusConflict,
				"can not delete a subject assigned to a schedule in %s, your department or other departments are still using this subject, clear the schedules first",
				Curriculum.SEMESTER_INDEX_NAME[semester],
			)

			return
		}
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteSubject(uint16(subject_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the subject")
		return
	}

	ctx.String(http.StatusOK, "subject deleted successfully")
}

func is_subject_assigned(university_schedules Schedule.UniTimeTables, subject_id uint16) bool {
	for _, section_week_schedules := range university_schedules {
		for _, day_time_table := range section_week_schedules {
			for _, time_slot := range day_time_table {
				if subject_id == time_slot.GetSubjectID() {
					return true
				}
			}
		}
	}

	return false
}
