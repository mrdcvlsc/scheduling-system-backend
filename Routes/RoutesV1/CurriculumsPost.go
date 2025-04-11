package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

/*
POST:

	"/curriculum_add"
*/
func PostCurriculum(ctx *gin.Context) {
	add_curriculum := Curriculum.Curriculum{}

	if err := ctx.BindJSON(&add_curriculum); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the curriculum to be added")
		return
	}

	RouteGlobals.ReindexUniSchedMutex.Lock()
	defer RouteGlobals.ReindexUniSchedMutex.Unlock()

	// save new curriculum

	err := RouteGlobals.ResourcesPersistence.WriterService.CreateCurriculum(add_curriculum)

	if err != nil {
		log.Print(err)
		ctx.String(http.StatusBadRequest, "we are unable to properly add the curriculum")
		return
	}

	// rebuild university schedule index

	all_curriculums, err_read_all_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_read_all_curriculums != nil {
		ctx.String(http.StatusInternalServerError, "we are unable retrieve the curriculums right now")
		return
	}

	for selected_semester := range Curriculum.SUPPORTED_SEMESTERS {

		// obtain univesity schedules for each semester

		university_schedule, has_obtain := ObtainUniversityScheduleNoHorizontalValidation(ctx, selected_semester)

		if !has_obtain {
			return
		}

		// determine insert university schedule index for the new curriculum

		schedule_idx := 0
		insert_idx := -1
		insert_length := 0

		for _, curriculum := range all_curriculums {
			for _, year_level := range curriculum.YearLevels {

				if !year_level.IsActive {
					continue
				}

				for semester_idx, semester := range year_level.Semesters {
					if semester_idx != selected_semester {
						continue
					}

					for section_idx := 0; section_idx < semester.Sections; section_idx++ {

						is_equal_code := Utils.IsEqualStrCaseInsensitiveIgnoreWhiteSpace(curriculum.CurriculumCode, add_curriculum.CurriculumCode)
						is_equal_name := Utils.IsEqualStrCaseInsensitiveIgnoreWhiteSpace(curriculum.CurriculumName, add_curriculum.CurriculumName)

						if is_equal_code && is_equal_name {
							if insert_idx == -1 {
								insert_idx = schedule_idx
							}

							insert_length++
						}

						schedule_idx++
					}
				}
			}
		}

		if insert_idx < 0 {
			log.Print("PostCurriculum: unable to find the university schedule insert index for the new curriculum")
			ctx.String(http.StatusInternalServerError, "we are unable to rebuld the university schedule index")
			return
		}

		// rebuild index of the new university schedule with the new curriculum

		if insert_idx > len(university_schedule) {
			log.Printf("PostCurriculum: insert_idx %d exceeds university_schedule length %d", insert_idx, len(university_schedule))
			ctx.String(http.StatusInternalServerError, "unable to rebuild university schedule")
			return
		}

		new_university_schedule := make(Schedule.UniTimeTables, 0, len(university_schedule))
		new_university_schedule = append(new_university_schedule, university_schedule[:insert_idx]...)
		new_university_schedule = append(new_university_schedule, make(Schedule.UniTimeTables, insert_length)...)
		new_university_schedule = append(new_university_schedule, university_schedule[insert_idx:]...)

		// save the new university schedules

		err_save_schedules := RouteGlobals.SchedulePersistence.SaveService.SaveSchedules(new_university_schedule, selected_semester)

		if err_save_schedules != nil {
			ctx.String(http.StatusInternalServerError, "we're unable to save the deletion of the curriculum from the university schedules right now")
			return
		}

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(selected_semester, new_university_schedule)

		if err_set_cache != nil {
			log.Print("PostCurriculum: the new university schedule was saved, but we're unable to cache it")
		}
	}

	ctx.String(http.StatusOK, "curriculum added successfully")
}
