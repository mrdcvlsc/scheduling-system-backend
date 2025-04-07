package RoutesV1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

/*
retrieves university schedule from cache or persistence, can return empty university schedule if there is not university schedule generated yet.

example usage inside a gin route:

	// setting departments_to_validate to nil will validate all departments.
	university_schedules, is_success := ObtainUniversitySchedule(ctx, nil, selected_semester)
	if !is_success {
		return
	}
*/
func ObtainUniversitySchedule(ctx *gin.Context, departments_to_validate map[uint16]bool, semester int) (Schedule.UniTimeTables, bool) {
	var university_schedules Schedule.UniTimeTables = nil

	cached_university_schedule, has_cache, err_get_cache := RouteGlobals.GetCachedUniversitySchedule(semester)

	if err_get_cache != nil {
		log.Println("ObtainUniversitySchedule:", err_get_cache.Error())
		ctx.String(http.StatusBadRequest, err_get_cache.Error())
		return nil, false
	}

	if has_cache {
		log.Println("ObtainUniversitySchedule: retrieving university schedule from cache.")
		university_schedules = cached_university_schedule
	} else {
		log.Println("ObtainUniversitySchedule: no cached detected loading from persistence")
		read_university_schedules, err_load_schedules := RouteGlobals.SchedulePersistence.LoadService.LoadSchedules(semester)

		if err_load_schedules != nil {
			log.Println("ObtainUniversitySchedule:", err_load_schedules)

			if errors.Is(err_load_schedules, os.ErrNotExist) {
				log.Printf("ObtainUniversitySchedule: schedule for semester index %d is not created yet, creating an empty schedule instead", semester)

				curriculums, err_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

				if err_curriculums != nil {
					log.Fatal("ObtainUniversitySchedule:", err_curriculums)
				}

				university_schedules = GeneticAlgorithm.NewEmptyIndividual(curriculums, semester)
			} else {
				ctx.String(http.StatusInternalServerError, "we failed to read that schedule")
				return nil, false
			}
		} else {
			university_schedules = read_university_schedules
		}
	}

	if university_schedules.IsEmpty() {
		return university_schedules, true
	}

	for _, err_vertical_validation := range university_schedules.VerticalValidation(RouteGlobals.ResourcesPersistence) {
		if err_vertical_validation != nil {
			log.Println("ObtainUniversitySchedule: invalid schedule detected, vertical overlap")
			ctx.String(http.StatusConflict, "server detected an invalid schedule with vertically overlapping data")
			return nil, false
		}
	}

	for _, err_horizontal_validation := range university_schedules.HorizontalValidation(RouteGlobals.ResourcesPersistence, departments_to_validate, semester) {
		if err_horizontal_validation != nil {
			log.Println("ObtainUniversitySchedule: invalid schedule detected, horizontal overlap")
			ctx.String(http.StatusConflict, "server detected an invalid schedule with wrong horizontal data allocations")
			return nil, false
		}
	}

	log.Println("ObtainUniversitySchedule: schedule found")
	return university_schedules, true
}

/*
retrieves university schedule from cache or persistence, can return empty university schedule if there is not university schedule generated yet.

example usage inside a gin route:

	// setting departments_to_validate to nil will validate all departments.
	university_schedules, is_success := ObtainUniversityScheduleNoHorizontalValidation(ctx, selected_semester)
	if !is_success {
		return
	}
*/
func ObtainUniversityScheduleNoHorizontalValidation(ctx *gin.Context, semester int) (Schedule.UniTimeTables, bool) {
	var university_schedules Schedule.UniTimeTables = nil

	cached_university_schedule, has_cache, err_get_cache := RouteGlobals.GetCachedUniversitySchedule(semester)

	if err_get_cache != nil {
		log.Println("ObtainUniversityScheduleNoHorizontalValidation:", err_get_cache.Error())
		ctx.String(http.StatusBadRequest, err_get_cache.Error())
		return nil, false
	}

	if has_cache {
		log.Println("ObtainUniversityScheduleNoHorizontalValidation: retrieving university schedule from cache.")
		university_schedules = cached_university_schedule
	} else {
		log.Println("ObtainUniversityScheduleNoHorizontalValidation: no cached detected loading from persistence")
		read_university_schedules, err_load_schedules := RouteGlobals.SchedulePersistence.LoadService.LoadSchedules(semester)

		if err_load_schedules != nil {
			log.Println("ObtainUniversityScheduleNoHorizontalValidation:", err_load_schedules)

			if errors.Is(err_load_schedules, os.ErrNotExist) {
				log.Printf("ObtainUniversityScheduleNoHorizontalValidation: schedule for semester index %d is not created yet, creating an empty schedule instead", semester)

				curriculums, err_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

				if err_curriculums != nil {
					log.Fatal("ObtainUniversityScheduleNoHorizontalValidation:", err_curriculums)
				}

				university_schedules = GeneticAlgorithm.NewEmptyIndividual(curriculums, semester)
			} else {
				ctx.String(http.StatusInternalServerError, "we failed to read that schedule")
				return nil, false
			}
		} else {
			university_schedules = read_university_schedules
		}
	}

	if university_schedules.IsEmpty() {
		return university_schedules, true
	}

	for _, err_vertical_validation := range university_schedules.VerticalValidation(RouteGlobals.ResourcesPersistence) {
		if err_vertical_validation != nil {
			log.Println("ObtainUniversityScheduleNoHorizontalValidation: invalid schedule detected - vertical overlap")
			ctx.String(http.StatusConflict, "server detected an invalid schedule with vertically overlapping data")
			return nil, false
		}
	}

	log.Println("ObtainUniversityScheduleNoHorizontalValidation: schedule found")
	return university_schedules, true
}

/*
retrieves university schedule from cache or persistence, can return empty university schedule if there is not university schedule generated yet.

this obtain function can return non-nil university schedules despite having a vertical or horizontal data overlaps

example usage inside a gin route:

	// setting departments_to_validate to nil will validate all departments.
	university_schedules, err_obtain_sched_no_ctx := ObtainUniversityScheduleNoContext(nil, selected_semester)
	if err_obtain_sched_no_ctx != nil {
		// handle error
		return
	}
*/
func ObtainUniversityScheduleNoContext(departments_to_validate map[uint16]bool, semester int) (Schedule.UniTimeTables, error) {
	var university_schedules Schedule.UniTimeTables = nil

	cached_university_schedule, has_cache, err_get_cache := RouteGlobals.GetCachedUniversitySchedule(semester)

	if err_get_cache != nil {
		log.Println("ObtainUniversityScheduleNoContext:", err_get_cache.Error())
		return nil, err_get_cache
	}

	if has_cache {
		log.Println("ObtainUniversityScheduleNoContext: retrieving university schedule from cache.")
		university_schedules = cached_university_schedule
	} else {
		log.Println("ObtainUniversityScheduleNoContext: no cached detected loading from persistence")
		read_university_schedules, err_load_schedules := RouteGlobals.SchedulePersistence.LoadService.LoadSchedules(semester)

		if err_load_schedules != nil {
			log.Println("ObtainUniversityScheduleNoContext:", err_load_schedules)

			if errors.Is(err_load_schedules, os.ErrNotExist) {
				log.Printf("ObtainUniversityScheduleNoContext: schedule for semester index %d is not created yet, creating an empty schedule instead", semester)

				curriculums, err_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

				if err_curriculums != nil {
					log.Fatal("ObtainUniversityScheduleNoContext:", err_curriculums)
				}

				university_schedules = GeneticAlgorithm.NewEmptyIndividual(curriculums, semester)
			} else {
				return nil, fmt.Errorf("failed to read the university schedule for the %s", Curriculum.SEMESTER_INDEX_NAME[semester])
			}
		} else {
			university_schedules = read_university_schedules
		}
	}

	if university_schedules.IsEmpty() {
		return university_schedules, nil
	}

	for _, err_vertical_validation := range university_schedules.VerticalValidation(RouteGlobals.ResourcesPersistence) {
		if err_vertical_validation != nil {
			log.Println("ObtainUniversityScheduleNoContext: invalid schedule detected, vertical overlaps")
			return university_schedules, errors.New("server detected an invalid schedule with vertically overlapping data")
		}
	}

	for _, err_horizontal_validation := range university_schedules.HorizontalValidation(RouteGlobals.ResourcesPersistence, departments_to_validate, semester) {
		if err_horizontal_validation != nil {
			log.Println("ObtainUniversityScheduleNoContext: invalid schedule detected, horizontal overlaps")
			return university_schedules, errors.New("server detected an invalid schedule with horizontally overlapping data")
		}
	}

	log.Println("ObtainUniversityScheduleNoContext: schedule found")
	return university_schedules, nil
}

/*
retrieves university schedule from cache or persistence, can return empty university schedule if there is not university schedule generated yet.

this obtain function can return non-nil university schedules despite having a vertical or horizontal data overlaps

example usage inside a gin route:

	// setting departments_to_validate to nil will validate all departments.
	university_schedules, err_obtain_sched_no_ctx := ObtainUniversityScheduleNoContext(nil, selected_semester)
	if err_obtain_sched_no_ctx != nil {
		// handle error
		return
	}
*/
func ObtainUniversityScheduleNoContextNoHorizontalValidation(semester int) (Schedule.UniTimeTables, error) {
	var university_schedules Schedule.UniTimeTables = nil

	cached_university_schedule, has_cache, err_get_cache := RouteGlobals.GetCachedUniversitySchedule(semester)

	if err_get_cache != nil {
		log.Println("ObtainUniversityScheduleNoContextNoHorizontalValidation:", err_get_cache.Error())
		return nil, err_get_cache
	}

	if has_cache {
		log.Println("ObtainUniversityScheduleNoContextNoHorizontalValidation: retrieving university schedule from cache.")
		university_schedules = cached_university_schedule
	} else {
		log.Println("ObtainUniversityScheduleNoContextNoHorizontalValidation: no cached detected loading from persistence")
		read_university_schedules, err_load_schedules := RouteGlobals.SchedulePersistence.LoadService.LoadSchedules(semester)

		if err_load_schedules != nil {
			log.Println("ObtainUniversityScheduleNoContextNoHorizontalValidation:", err_load_schedules)

			if errors.Is(err_load_schedules, os.ErrNotExist) {
				log.Printf("ObtainUniversityScheduleNoContextNoHorizontalValidation: schedule for semester index %d is not created yet, creating an empty schedule instead", semester)

				curriculums, err_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

				if err_curriculums != nil {
					log.Fatal("ObtainUniversityScheduleNoContextNoHorizontalValidation:", err_curriculums)
				}

				university_schedules = GeneticAlgorithm.NewEmptyIndividual(curriculums, semester)
			} else {
				return nil, fmt.Errorf("failed to read the university schedule for the %s", Curriculum.SEMESTER_INDEX_NAME[semester])
			}
		} else {
			university_schedules = read_university_schedules
		}
	}

	if university_schedules.IsEmpty() {
		return university_schedules, nil
	}

	for _, err_vertical_validation := range university_schedules.VerticalValidation(RouteGlobals.ResourcesPersistence) {
		if err_vertical_validation != nil {
			log.Println("ObtainUniversityScheduleNoContextNoHorizontalValidation: invalid schedule detected, vertical overlaps")
			return university_schedules, errors.New("server detected an invalid schedule with vertically overlapping data")
		}
	}

	log.Println("ObtainUniversityScheduleNoContextNoHorizontalValidation: schedule found")
	return university_schedules, nil
}
