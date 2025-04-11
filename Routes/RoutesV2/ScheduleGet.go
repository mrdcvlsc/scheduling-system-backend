package RoutesV2

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Routes/RoutesV1"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
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

	// TODO: use department_id for authentication later on.

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

type WeekTimeTableSubjects []RoutesV1.SubjectAssignmentInfo

/*
POST:

	"/add_schedule_preference"
*/
func PostWeekTimeTableSurvery(ctx *gin.Context) {
	var configured_week_time_table []RoutesV1.SubjectAssignmentInfo

	if err := ctx.BindJSON(&configured_week_time_table); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read the prefered configured week time table")
		return
	}

	Utils.PrettyPrint(configured_week_time_table)

	///////

	var configured_week_time_table_records [][]RoutesV1.SubjectAssignmentInfo

	file_data, err := os.ReadFile("week-time-table-scedule-preferences.json")

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Printf("error reading file: %v\n", err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		configured_week_time_table_records = make([][]RoutesV1.SubjectAssignmentInfo, 0)
	} else {
		// File exists — try to parse the JSON
		if err := json.Unmarshal(file_data, &configured_week_time_table_records); err != nil {
			log.Printf("error parsing JSON: %v\n", err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	configured_week_time_table_records = append(configured_week_time_table_records, configured_week_time_table)

	updated, err := json.MarshalIndent(configured_week_time_table_records, "", " ")

	if err != nil {
		log.Printf("error marshaling: %v\n", err)
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := os.WriteFile("week-time-table-scedule-preferences.json", updated, 0644); err != nil {
		log.Printf("error writing file: %v\n", err)
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("save done | total number of records : %d", len(configured_week_time_table_records))

	//////////

	ctx.String(http.StatusOK, "week time table preference added")
}
