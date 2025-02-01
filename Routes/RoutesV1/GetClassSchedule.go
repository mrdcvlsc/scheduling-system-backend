package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
)

type SubjectAssignmentInfo struct {
	SubjectCode        string `json:"SubjectCode"`
	InstructorLastName string `json:"InstructorLastName"`
	RoomName           string `json:"RoomName"`
	DayIdx             uint8  `json:"DayIdx"`
	TimeSlotIdx        uint8  `json:"TimeSlotIdx"`
	SubjectTimeSlots   uint8  `json:"SubjectTimeSlots"`
}

/*
GET :

	"/v1/class_json_schedule?department_id=D&semester=S&schedule_idx=I"
*/
func GetJsonClassSchedule(ctx *gin.Context) {

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

	university_schedules, has_obtained := ObtainUniversitySchedule(ctx, semester)

	if !has_obtained {
		return
	}

	// cache the found university schedule for the semester

	cache_err := RouteGlobals.SetCachedUniversitySchedule(semester, university_schedules)

	if cache_err != nil {
		log.Println(cache_err.Error())
	}

	// parse schedule_idx parameter

	schedule_idx, is_valid_idx := IsValidUniversityScheduleIndex(ctx, university_schedules)

	if !is_valid_idx {
		return
	}

	// extract selected schedule

	selected_class_schedule := university_schedules[schedule_idx:(schedule_idx + 1)]

	// process schedule information

	sub_id_to_subject_code := make(map[uint16]string)

	{ // subjects
		subjects, subjects_read_err := RouteGlobals.ResourcesPersistence.ReaderService.GetAllSubjects()

		if subjects_read_err != nil {
			ctx.String(http.StatusInternalServerError, "we can not retrieve the subjects information right now")
			return
		}

		for _, subject := range subjects {
			sub_id_to_subject_code[subject.ID] = subject.Code
		}
	}

	instructor_id_to_instructor_name := make(map[uint16]string)

	{ // instructors
		instructors, instructors_read_err := RouteGlobals.ResourcesPersistence.ReaderService.GetAllInstructors()

		if instructors_read_err != nil {
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
		rooms, rooms_read_err := RouteGlobals.ResourcesPersistence.ReaderService.GetAllRooms()

		if rooms_read_err != nil {
			ctx.String(http.StatusInternalServerError, "we can not retrieve the rooms information right now")
			return
		}

		for _, room := range rooms {
			if room.DepartmentID == uint16(department_id) || room.DepartmentID == 0 {
				room_id_to_room_name[room.RoomID] = room.Name
			}
		}
	}

	sub_assign_info := make([]SubjectAssignmentInfo, 0)

	for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
		for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
			slot := selected_class_schedule[0][day][time_slot]

			if slot.GetSubjectID() == 0 {
				continue
			}

			new_sub_assignment := SubjectAssignmentInfo{
				SubjectCode:        sub_id_to_subject_code[slot.GetSubjectID()],
				InstructorLastName: instructor_id_to_instructor_name[slot.GetInstructorID()],
				RoomName:           room_id_to_room_name[slot.GetRoomID()],
				DayIdx:             uint8(day),
				TimeSlotIdx:        uint8(time_slot),
				SubjectTimeSlots:   1,
			}

			for forward_time_slot := time_slot + 1; forward_time_slot < Const.N_DAILY_TIME_SLOTS; forward_time_slot++ {
				forward_slot := selected_class_schedule[0][day][forward_time_slot]

				if forward_slot.GetSubjectID() == slot.GetSubjectID() && forward_slot.GetRoomID() == slot.GetRoomID() {
					new_sub_assignment.SubjectTimeSlots++
				} else {
					time_slot = forward_time_slot - 1
					break
				}
			}

			sub_assign_info = append(sub_assign_info, new_sub_assignment)
		}
	}

	ctx.JSON(http.StatusOK, sub_assign_info)
}

// GET : /v1/class_schedule?department_id=D&semester=S&schedule_idx=I
func GetClassSchedule(ctx *gin.Context) {

	// parse semester parameter

	semester, is_valid_semester_param := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	// parse department_id parameter

	// TODO: use department_id for authentication later on.

	_, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	// load university schedules

	university_schedules, has_obtained := ObtainUniversitySchedule(ctx, semester)

	if !has_obtained {
		return
	}

	// cache the found university schedule for the semester

	cache_err := RouteGlobals.SetCachedUniversitySchedule(semester, university_schedules)

	if cache_err != nil {
		log.Println(cache_err.Error())
	}

	// parse schedule_idx parameter

	schedule_idx, is_valid_idx := IsValidUniversityScheduleIndex(ctx, university_schedules)

	if !is_valid_idx {
		return
	}

	// extract selected schedule

	serialized_schedule := Schedule.SerializeUniversitySchedule(university_schedules[schedule_idx:(schedule_idx + 1)])

	ctx.Data(http.StatusOK, "application/octet-stream", serialized_schedule)
}
