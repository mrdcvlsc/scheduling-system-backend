package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
GET:

	"/v1/clear_class_schedule?department_id=D&semester=S&schedule_idx=I"

the `schedule_idx` for a section can be fetch from
`GetDepartmentData` function using rest api GET request:

	"/v1/department_data?department_id=[N>0]&semester=[0-1]"
*/
func DeleteClearClassSchedule(ctx *gin.Context) {

	// parse semester parameter

	semester, is_valid_semester_param := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	// parse department_id parameter

	// TODO: use department_id for authentication later on.

	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	// load university schedules

	department_to_validate := make(map[uint16]bool)
	department_to_validate[uint16(department_id)] = true

	university_schedules, has_obtained := ObtainUniversityScheduleNoHorizontalValidation(ctx, semester)

	if !has_obtained {
		return
	}

	// cache the found university schedule for the semester

	err_set_cache := RouteGlobals.SetCachedUniversitySchedule(semester, university_schedules)

	if err_set_cache != nil {
		log.Println(err_set_cache.Error())
	}

	// parse schedule_idx parameter

	schedule_idx, is_valid_idx := IsValidUniversityScheduleIndex(ctx, university_schedules)

	if !is_valid_idx {
		return
	}

	// TODO: authenticate `schedule_idx` parameter, when someone deletes a curriculum

	// clear class schedule

	university_schedules[schedule_idx] = Schedule.WeekTimeTable{}

	// save schedule

	err_save_schedules := RouteGlobals.SchedulePersistence.SaveService.SaveSchedules(university_schedules, semester)

	if err_save_schedules != nil {
		log.Print("DeleteClearClassSchedule: (save error) ", err_save_schedules.Error())
		ctx.String(http.StatusOK, "we're unable to clear that schedule's week time table")
		return
	}

	ctx.String(http.StatusOK, "weekly time table schedule was successfully cleared")
}

/*
GET:

	"/v1/clear_department_schedules?department_id=[N>0]&semester=[0-1]"
*/
func DeleteClearDepartmentSchedule(ctx *gin.Context) {

	// parse semester parameter

	semester, is_valid_semester_param := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	// parse department_id parameter

	// TODO: use department_id for authentication later on.

	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	// load university schedules

	department_to_validate := make(map[uint16]bool)
	department_to_validate[uint16(department_id)] = true

	university_schedules, has_obtained := ObtainUniversityScheduleNoHorizontalValidation(ctx, semester)

	if !has_obtained {
		return
	}

	// cache found for university schedule for the current semester

	err_set_cache := RouteGlobals.SetCachedUniversitySchedule(semester, university_schedules)

	if err_set_cache != nil {
		log.Println(err_set_cache.Error())
	}

	// TODO: authenticate `schedule_idx` parameter, when someone deletes a curriculum

	// get all curriculums

	all_curriculums, err_read_all_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_read_all_curriculums != nil {
		ctx.String(http.StatusInternalServerError, "we are unable retrieve the curriculums right now")
		return
	}

	// clear class schedule

	schedule_idx := 0

	for _, curriculum := range all_curriculums {
		for _, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue
			}

			for semester_idx, semester_element := range year_level.Semesters {
				if semester_idx != semester {
					continue
				}

				for section_idx := 0; section_idx < semester_element.Sections; section_idx++ {

					if curriculum.DepartmentID == uint16(department_id) {
						university_schedules[schedule_idx] = Schedule.WeekTimeTable{}
					}

					schedule_idx++
				}
			}
		}
	}

	// save schedule

	err_save_schedules := RouteGlobals.SchedulePersistence.SaveService.SaveSchedules(university_schedules, semester)

	if err_save_schedules != nil {
		log.Print("DeleteClearClassSchedule: (save error) ", err_save_schedules.Error())
		ctx.String(http.StatusOK, "we're unable to clear the department schedules")
		return
	}

	ctx.String(http.StatusOK, "all department schedule was successfully cleared")
}
