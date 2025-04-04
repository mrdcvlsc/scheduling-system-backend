package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

/*
DELETE:

	"/curriculum_remove?curriculum_id=[N>0]"
*/
func DeleteCurriculum(ctx *gin.Context) {
	curriculum_id, is_valid_curriculum_id_param := IsValidCurriculumID(ctx)

	if !is_valid_curriculum_id_param {
		return
	}

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

		// determine which schedule indices should be removed from the current university schedules

		schedule_idx := 0
		remove_starting_index := -1
		remove_chunk_length := 0

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

						if curriculum.CurriculumID == uint16(curriculum_id) {

							if remove_starting_index == -1 {
								remove_starting_index = schedule_idx
							}

							remove_chunk_length++
						}

						schedule_idx++
					}
				}
			}
		}

		// remove the to be deleted schedule indices from the university schedules

		log.Printf("DeleteCurriculum: removing university schedule index %d to %d (of size %d)", remove_starting_index, remove_starting_index+remove_chunk_length, remove_chunk_length)

		new_university_schedule, err_remove_chunk_in_slice := Utils.RemoveChunkInSlice(university_schedule, remove_starting_index, remove_chunk_length)

		if err_remove_chunk_in_slice != nil {
			ctx.String(http.StatusInternalServerError, "we're unable to remove the curriculum to the university schedules right now")
			return
		}

		// save the new university schedules

		err_save_schedules := RouteGlobals.SchedulePersistence.SaveService.SaveSchedules(new_university_schedule, selected_semester)

		if err_save_schedules != nil {
			ctx.String(http.StatusInternalServerError, "we're unable to save the deletion of the curriculum from the university schedules right now")
			return
		}

		err_set_cache := RouteGlobals.SetCachedUniversitySchedule(selected_semester, new_university_schedule)

		if err_set_cache != nil {
			log.Print("DeleteCurriculum: the new university schedule was saved, but we're unable to cache it")
		}
	}

	err := RouteGlobals.ResourcesPersistence.WriterService.DeleteCurriculum(uint16(curriculum_id))

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly remove the curriculum from the persistence")
		return
	}

	ctx.String(http.StatusOK, "curriculum deleted successfully")
}
