package RoutesV2

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Routes/RoutesV1"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

type InstructorTableItem struct {
	InstructorID  uint16 `json:"InstructorID"`
	DepartmentID  uint16 `json:"DepartmentID"`
	FirstName     string `json:"FirstName"`
	MiddleInitial string `json:"MiddleInitial"`
	LastName      string `json:"LastName"`
}

type InstructorTablePage struct {
	Instructors      []InstructorTableItem `json:"Instructors"`
	TotalInstructors int                   `json:"TotalInstructors"`
}

/*
GET:

	"/instructors?
		department_id=D&
		page_size=[N>0]&
		page[0-N>0]&
		firstname_match=<string>&
		initial_match=<string>&
		lastname_match=<string>
	"
*/
func GetDepartmentInstructors(ctx *gin.Context) {
	department_id, is_valid_department_id_param := RoutesV1.IsValidParameterDepartmentID(ctx)
	if !is_valid_department_id_param {
		return
	}

	page_size, is_valid_page_size_param := RoutesV1.IsValidPageSize(ctx)
	if !is_valid_page_size_param {
		return
	}

	page, is_valid_page_param := RoutesV1.IsValidPage(ctx)
	if !is_valid_page_param {
		return
	}

	firstname_match := ctx.Query("firstname_match")
	initial_match := ctx.Query("initial_match")
	lastname_match := ctx.Query("lastname_match")

	all_instructors, err_read_all_instructors := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllInstructors()

	if err_read_all_instructors != nil {
		log.Println(err_read_all_instructors)
		ctx.String(http.StatusInternalServerError, "we are unable to retrieve the instructors right now")
		return
	}

	department_instructors_page := make([]InstructorTableItem, 0)

	total_instructors := 0

	for _, instructor := range all_instructors {
		if int(instructor.DepartmentID) != department_id {
			continue
		}

		if firstname_match != "" && !Utils.HasSubString(instructor.FirstName, firstname_match) {
			continue
		}

		if initial_match != "" && !Utils.HasSubString(instructor.MiddleInitial, initial_match) {
			continue
		}

		if lastname_match != "" && !Utils.HasSubString(instructor.LastName, lastname_match) {
			continue
		}

		total_instructors++

		if (total_instructors - 1) < (page_size * page) {
			continue
		}

		if len(department_instructors_page) < page_size {
			department_instructors_page = append(department_instructors_page, InstructorTableItem{
				InstructorID:  instructor.InstructorID,
				DepartmentID:  instructor.DepartmentID,
				FirstName:     instructor.FirstName,
				MiddleInitial: instructor.MiddleInitial,
				LastName:      instructor.LastName,
			})
		}
	}

	instructor_table_page := &InstructorTablePage{
		Instructors:      department_instructors_page,
		TotalInstructors: total_instructors,
	}

	ctx.JSON(http.StatusOK, instructor_table_page)
}

/*
GET:

	"/instructor?instructor_id=[N>0]&department_id=[N>0]"
*/
func GetInstructorResource(ctx *gin.Context) {
	instructor_id, is_valid_instructor_id_param := RoutesV1.IsValidInstructorID(ctx)
	if !is_valid_instructor_id_param {
		return
	}

	department_id, is_valid_department_id_param := RoutesV1.IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	all_instructors, err_read_all_instructors := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllInstructors()

	if err_read_all_instructors != nil {
		log.Println(err_read_all_instructors)
		ctx.String(http.StatusInternalServerError, "we are unable to retrieve the instructors right now")
		return
	}

	var selected_instructor_base *Instructors.Instructor

	for _, instructor := range all_instructors {
		if instructor.InstructorID == uint16(instructor_id) {
			selected_instructor_base = &instructor
			break
		}
	}

	if selected_instructor_base == nil {
		ctx.String(http.StatusNotFound, "that instructor does not exist")
		return
	}

	all_curriculums, err_read_all_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_read_all_curriculums != nil {
		ctx.String(http.StatusInternalServerError, "we are unable to get the curriculums needed to generate the instructor timeslot availability")
		return
	}

	time_encodings := make(map[string]any, 0)
	time_encodings["base"] = selected_instructor_base.Time.Stringify()

	department_to_validate := make(map[uint16]bool)
	department_to_validate[uint16(department_id)] = true

	sched_1st_sem, has_obtained_1st_sem := RoutesV1.ObtainUniversitySchedule(ctx, department_to_validate, GeneticAlgorithm.TERM_1ST_SEMESTER)

	if !has_obtained_1st_sem {
		return
	}

	if !sched_1st_sem.IsEmpty() {
		sem_1st_time_allocation, sub_assign, err_get_instructor_time_allocation := get_instructor_time_allocation(
			*selected_instructor_base,
			sched_1st_sem, all_curriculums,
			GeneticAlgorithm.TERM_1ST_SEMESTER,
		)

		if err_get_instructor_time_allocation != nil {
			ctx.String(http.StatusInternalServerError, "we are unable to recreated the instructor time allocation for the 1st semester")
			return
		}

		time_encodings["sem_1st"] = sem_1st_time_allocation.Stringify()
		time_encodings["sem_1st_sub_assign"] = sub_assign
	}

	sched_2nd_sem, has_obtained_2nd_sem := RoutesV1.ObtainUniversitySchedule(ctx, department_to_validate, GeneticAlgorithm.TERM_2ND_SEMESTER)

	if !has_obtained_2nd_sem {
		return
	}

	if !sched_2nd_sem.IsEmpty() {
		sem_2nd_time_allocation, sub_assign, err_get_instructor_time_allocation := get_instructor_time_allocation(
			*selected_instructor_base,
			sched_2nd_sem, all_curriculums,
			GeneticAlgorithm.TERM_2ND_SEMESTER,
		)

		if err_get_instructor_time_allocation != nil {
			ctx.String(http.StatusInternalServerError, "we are unable to recreated the instructor time allocation for the 2nd semester")
			return
		}

		time_encodings["sem_2nd"] = sem_2nd_time_allocation.Stringify()
		time_encodings["sem_2nd_sub_assign"] = sub_assign
	}

	ctx.JSON(http.StatusOK, time_encodings)
}

type InstructorSubjectAssignmentInfo struct {
	SubjectCode      string `json:"SubjectCode"`
	CourseSection    string `json:"CourseSection"`
	RoomName         string `json:"RoomName"`
	DayIdx           uint8  `json:"DayIdx"`
	TimeSlotIdx      uint8  `json:"TimeSlotIdx"`
	SubjectTimeSlots uint8  `json:"SubjectTimeSlots"`
}

func get_instructor_time_allocation(base_instructor Instructors.Instructor, university_schedules Schedule.UniTimeTables, all_curriculums []Curriculum.Curriculum, selected_semester int) (Instructors.InstructorTimeSlotBitMap, []InstructorSubjectAssignmentInfo, error) {

	sub_id_to_subject_code := make(map[uint16]string)

	{ // subjects
		subjects, err_read_all_subjects := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllSubjects()

		if err_read_all_subjects != nil {
			return base_instructor.Time, nil, errors.New("we can not retrieve the subjects information right now")
		}

		for _, subject := range subjects {
			sub_id_to_subject_code[subject.ID] = subject.Code
		}
	}

	room_id_to_room_name := make(map[uint16]string)

	{ // rooms
		rooms, err_read_all_rooms := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllRooms()

		if err_read_all_rooms != nil {
			return base_instructor.Time, nil, errors.New("we can not retrieve the rooms information right now")
		}

		for _, room := range rooms {
			if room.DepartmentID == base_instructor.DepartmentID || room.DepartmentID == 0 {
				room_id_to_room_name[room.RoomID] = room.Name
			}
		}
	}

	////////////////

	SECTION := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	counted_sections := 0

	sub_assign_info := make([]InstructorSubjectAssignmentInfo, 0)

	for _, curriculum := range all_curriculums {
		for year_level_idx, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue // skip inactive year levels
			}

			for semester_idx, semester := range year_level.Semesters {

				if selected_semester != semester_idx {
					continue // skip not selected semesters
				}

				for section_idx := range semester.Sections {

					for day := range Const.N_WEEKLY_SCHOOL_DAYS {

						for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {

							subject_id := university_schedules[counted_sections][day].GetTimeSlot(time_slot).GetSubjectID()
							instructor_id := university_schedules[counted_sections][day].GetTimeSlot(time_slot).GetInstructorID()

							if subject_id != 0 && instructor_id == base_instructor.InstructorID {
								if base_instructor.InstructorID == 0 {
									log.Panic("there should be an instructor allocation here, why there is none?")
								}

								base_instructor.Time.SetAvailability(false, day, time_slot)

								room_id := university_schedules[counted_sections][day].GetTimeSlot(time_slot).GetRoomID()

								new_sub_assignment := InstructorSubjectAssignmentInfo{
									SubjectCode:      sub_id_to_subject_code[subject_id],
									CourseSection:    fmt.Sprintf("%s-%d%s", curriculum.CurriculumCode, year_level_idx+1, SECTION[section_idx]),
									RoomName:         room_id_to_room_name[room_id],
									DayIdx:           uint8(day),
									TimeSlotIdx:      uint8(time_slot),
									SubjectTimeSlots: 1,
								}

								for forward_time_slot := time_slot + 1; forward_time_slot < Const.N_DAILY_TIME_SLOTS; forward_time_slot++ {
									forward_slot := university_schedules[counted_sections][day].GetTimeSlot(forward_time_slot)

									if forward_slot.GetSubjectID() == subject_id && forward_slot.GetInstructorID() == instructor_id && forward_slot.GetRoomID() == room_id {
										new_sub_assignment.SubjectTimeSlots++
										base_instructor.Time.SetAvailability(false, day, forward_time_slot)
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
						} // ------------- end of time_slot loop -------------
					} // ------------- end of day loop -------------

					counted_sections++
				} // ------------- end of section_idx loop -------------
			} // ------------- end of semester_idx loop -------------
		} // ------------- end of year_level loop -------------
	} // ------------- end of curriculum loop -------------

	return base_instructor.Time, sub_assign_info, nil
}

/*
GET:

	"/instructor_basic?instructor_id=[N>0]"
*/
func GetInstructorBasic(ctx *gin.Context) {
	instructor_id, is_valid_instructor_id_param := RoutesV1.IsValidInstructorID(ctx)
	if !is_valid_instructor_id_param {
		return
	}

	all_instructors, err_read_all_instructors := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllInstructors()

	if err_read_all_instructors != nil {
		log.Println(err_read_all_instructors)
		ctx.String(http.StatusInternalServerError, "we are unable to retrieve the instructors right now")
		return
	}

	var selected_instructor_base *Instructors.Instructor

	for _, instructor := range all_instructors {
		if instructor.InstructorID == uint16(instructor_id) {
			selected_instructor_base = &instructor
			break
		}
	}

	if selected_instructor_base == nil {
		ctx.String(http.StatusNotFound, "that instructor does not exist")
		return
	}

	ctx.JSON(http.StatusOK, selected_instructor_base)
}
