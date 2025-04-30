package RoutesV2

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Routes/RoutesV1"
)

/*
GET:

	"/class_json_schedule?department_id=[N>0]&semester=[0-N>=1]&curriculum_id=[N>0]&year_level_idx=[0-N>=1]&section_idx=[0-N>=1]"
*/
func GetJsonClassSchedule(ctx *gin.Context) {

	// parse curriculum id

	curriculum_id, is_valid_curriculum_id_param := RoutesV1.IsValidCurriculumID(ctx)

	if !is_valid_curriculum_id_param {
		return
	}

	// pasrse year level index

	param_year_level_idx, is_valid_year_level_idx_param := RoutesV1.IsValidIndex(ctx, "year_level_idx")

	if !is_valid_year_level_idx_param {
		return
	}

	// parse section index

	param_section_idx, is_valid_section_idx_param := RoutesV1.IsValidIndex(ctx, "section_idx")

	if !is_valid_section_idx_param {
		return
	}

	// parse semester parameter

	selected_semester, is_valid_semester_param := RoutesV1.IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	// parse department_id parameter

	department_id, is_valid_department_id_param := RoutesV1.IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	// load university schedules

	university_schedules, has_obtained := RoutesV1.ObtainUniversityScheduleNoHorizontalValidation(ctx, selected_semester)

	if !has_obtained {
		return
	}

	// cache the found university schedule for the semester

	err_set_cache := RouteGlobals.SetCachedUniversitySchedule(selected_semester, university_schedules)

	if err_set_cache != nil {
		log.Println(err_set_cache.Error())
	}

	// get all curriculums

	all_curriculums, err_read_all_curriculum := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_read_all_curriculum != nil {
		ctx.String(http.StatusInternalServerError, "unable to read curriculums for that department")
		return
	}

	// parse schedule_idx parameter

	schedule_idx := 0

curriculum_loop:
	for _, curriculum := range all_curriculums {
		for year_level_idx, year_level := range curriculum.YearLevels {
			if !year_level.IsActive {
				continue
			}

			for semester_idx, semester := range year_level.Semesters {
				if semester_idx != selected_semester {
					continue
				}

				for section_idx := 0; section_idx < semester.Sections; section_idx++ {

					if curriculum_id == int(curriculum.CurriculumID) && year_level_idx == param_year_level_idx && section_idx == param_section_idx {
						break curriculum_loop
					}

					schedule_idx++
				}
			}
		}
	}

	// extract selected schedule

	selected_class_schedule := university_schedules[schedule_idx:(schedule_idx + 1)]

	department_to_measure := make(map[uint16]bool)
	department_to_measure[uint16(department_id)] = true

	log.Printf("measured fitness : %f", GeneticAlgorithm.MeasureWeekTimeTableBasicFitness(
		selected_class_schedule[0],
	))

	// process schedule information

	sub_id_to_subject_code := make(map[uint16]string)

	{ // subjects
		subjects, err_read_all_subjects := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllSubjects()

		if err_read_all_subjects != nil {
			ctx.String(http.StatusInternalServerError, "we can not retrieve the subjects information right now")
			return
		}

		for _, subject := range subjects {
			sub_id_to_subject_code[subject.ID] = subject.Code
		}
	}

	instructor_id_to_instructor_name := make(map[uint16]string)

	{ // instructors
		instructors, err_read_all_instructors := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllInstructors()

		if err_read_all_instructors != nil {
			ctx.String(http.StatusInternalServerError, "we can not retrieve the instructors information right now")
			return
		}

		for _, instructor := range instructors {
			if instructor.DepartmentID == uint16(department_id) || instructor.DepartmentID == 0 {
				instructor_id_to_instructor_name[instructor.InstructorID] = instructor.LastName
			}
		}
	}

	room_id_to_room_name := make(map[uint16]string)

	{ // rooms
		rooms, err_read_all_rooms := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllRooms()

		if err_read_all_rooms != nil {
			ctx.String(http.StatusInternalServerError, "we can not retrieve the rooms information right now")
			return
		}

		for _, room := range rooms {
			if room.DepartmentID == uint16(department_id) || room.DepartmentID == 0 {
				room_id_to_room_name[room.RoomID] = room.Name
			}
		}
	}

	sub_assign_info := make([]RoutesV1.SubjectAssignmentInfo, 0)

	for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
		for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
			slot := selected_class_schedule[0][day].GetTimeSlot(time_slot)

			if slot.GetSubjectID() == 0 {
				continue
			}

			new_sub_assignment := RoutesV1.SubjectAssignmentInfo{
				SubjectID:          slot.GetSubjectID(),
				InstructorID:       slot.GetInstructorID(),
				RoomID:             slot.GetRoomID(),
				SubjectCode:        sub_id_to_subject_code[slot.GetSubjectID()],
				InstructorLastName: instructor_id_to_instructor_name[slot.GetInstructorID()],
				RoomName:           room_id_to_room_name[slot.GetRoomID()],
				DayIdx:             uint8(day),
				TimeSlotIdx:        uint8(time_slot),
				SubjectTimeSlots:   1,
			}

			for forward_time_slot := time_slot + 1; forward_time_slot < Const.N_DAILY_TIME_SLOTS; forward_time_slot++ {
				forward_slot := selected_class_schedule[0][day].GetTimeSlot(forward_time_slot)

				if forward_slot.GetSubjectID() == slot.GetSubjectID() && forward_slot.GetInstructorID() == slot.GetInstructorID() && forward_slot.GetRoomID() == slot.GetRoomID() {
					new_sub_assignment.SubjectTimeSlots++
				} else {
					time_slot = forward_time_slot - 1
					break
				}

				if forward_time_slot == (Const.N_DAILY_TIME_SLOTS - 1) {
					time_slot = 9999
					break
				}
			}

			sub_assign_info = append(sub_assign_info, new_sub_assignment)
		}
	}

	ctx.JSON(http.StatusOK, sub_assign_info)
}

/*
GET:

	"/validate_schedules?department_id=[N>0]&semester=[0-N>=1]"
*/
func GetValidateSchedules(ctx *gin.Context) {

	// parse parameters

	selected_semester, is_valid_semester_param := RoutesV1.IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	department_id, is_valid_department_id_param := RoutesV1.IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	department_to_horizontal_validate := make(map[uint16]bool)
	department_to_horizontal_validate[uint16(department_id)] = true

	// load university schedules

	university_schedules, err_obtain := RoutesV1.ObtainUniversityScheduleNoValidation(selected_semester)

	if err_obtain != nil {
		log.Print("GetClassScheduleValidate - obtain error:", err_obtain)
		ctx.String(http.StatusInternalServerError, err_obtain.Error())
		return
	}

	// cache the found university schedule for the semester

	err_set_cache := RouteGlobals.SetCachedUniversitySchedule(selected_semester, university_schedules)

	if err_set_cache != nil {
		log.Println("GetClassScheduleValidate - cache error:", err_set_cache.Error())
	}

	// get all curriculums

	// all_curriculums, err_read_all_curriculum := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	// if err_read_all_curriculum != nil {
	// 	ctx.String(http.StatusInternalServerError, "unable to read curriculums for that department")
	// 	return
	// }

	// extract selected schedule

	validation_results := make([]any, 0)

	if university_schedules.IsEmpty() {
		validation_results = append(validation_results, "all university schedules are empty")
		ctx.JSON(http.StatusNotFound, validation_results)
		return
	}

	rooms, err_read_all_rooms := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllRooms()

	if err_read_all_rooms != nil {
		ctx.String(http.StatusInternalServerError, "we can not retrieve the rooms information right now")
		return
	}

	curriculums, err_read_all_curriculum := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_read_all_curriculum != nil {
		ctx.String(http.StatusInternalServerError, "we're unable to read curriculums information right now")
		return
	}

	errs_vertical_validation := university_schedules.VerticalValidation(rooms)

	for _, err_vertical_validation := range errs_vertical_validation {
		if err_vertical_validation != nil {
			validation_results = append(validation_results, err_vertical_validation.Error())
		}
	}

	if len(errs_vertical_validation) > 0 {
		ctx.JSON(http.StatusConflict, validation_results)
		return
	}

	errs_horizontal_validation := university_schedules.HorizontalValidation(curriculums, department_to_horizontal_validate, selected_semester)

	for _, err_horizontal_validation := range errs_horizontal_validation {
		if err_horizontal_validation != nil {
			validation_results = append(validation_results, err_horizontal_validation.Error())
		}
	}

	if len(errs_horizontal_validation) > 0 {
		ctx.JSON(http.StatusConflict, validation_results)
		return
	}

	ctx.JSON(http.StatusOK, validation_results)
}
