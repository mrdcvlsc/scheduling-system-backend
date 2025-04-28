package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
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

	if RouteGlobals.IsGeneratingSchedule.Load() {
		log.Print("PostCurriculum: [busy] you or other department(s) are still generating a schedule, please wait until the process is finished")
		ctx.String(http.StatusForbidden, "we're unable to add a curriculum right now, you or other department(s) are still generating a schedule, please wait a little while until those process are done")
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

		insert_idx := -1
		insert_length := 0

		GeneticAlgorithm.IterateSectionsWeekSchedule(university_schedule, all_curriculums, selected_semester, nil, nil,
			func(indicies GeneticAlgorithm.IterIndices, values GeneticAlgorithm.IterValues) GeneticAlgorithm.IterReturnType {

				is_equal_code := Utils.IsEqualStrCaseInsensitiveIgnoreWhiteSpace(values.Curriculum.CurriculumCode, add_curriculum.CurriculumCode)
				is_equal_name := Utils.IsEqualStrCaseInsensitiveIgnoreWhiteSpace(values.Curriculum.CurriculumName, add_curriculum.CurriculumName)

				if is_equal_code && is_equal_name {
					if insert_idx == -1 {
						insert_idx = indicies.Usi
					}

					insert_length++
				}

				return GeneticAlgorithm.IterProceed
			},
		)

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
