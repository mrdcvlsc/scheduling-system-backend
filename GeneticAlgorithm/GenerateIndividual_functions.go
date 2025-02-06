package GeneticAlgorithm

import (
	"fmt"
	"math"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func GenerateMapDeptIdToRoomTypeToRooms(persistence *StorageResources.Persistence) (map[uint16]map[uint16][]Rooms.Room, error) {
	department_id_to_room_type_to_rooms := make(map[uint16]map[uint16][]Rooms.Room)

	rooms, err := persistence.ReaderService.GetAllRooms()

	if err != nil {
		return nil, err
	}

	for _, room := range rooms {
		_, has_department_id := department_id_to_room_type_to_rooms[room.DepartmentID]

		if !has_department_id {
			department_id_to_room_type_to_rooms[room.DepartmentID] = make(map[uint16][]Rooms.Room)
		}

		_, has_room_type := department_id_to_room_type_to_rooms[room.DepartmentID][room.RoomType]

		if !has_room_type {
			department_id_to_room_type_to_rooms[room.DepartmentID][room.RoomType] = make([]Rooms.Room, 1, 8)
			department_id_to_room_type_to_rooms[room.DepartmentID][room.RoomType][0] = room
		} else {
			department_id_to_room_type_to_rooms[room.DepartmentID][room.RoomType] = append(
				department_id_to_room_type_to_rooms[room.DepartmentID][room.RoomType], room,
			)
		}
	}

	return department_id_to_room_type_to_rooms, nil
}

func GenerateMapDeptIdToInstructors(persistence *StorageResources.Persistence) (map[uint16][]Instructors.Instructor, error) {
	department_id_to_instructors := make(map[uint16][]Instructors.Instructor)

	instructors, err := persistence.ReaderService.GetAllInstructors()

	if err != nil {
		return nil, err
	}

	for _, instructor := range instructors {
		_, has_department_id := department_id_to_instructors[instructor.DepartmentID]

		if !has_department_id {
			department_id_to_instructors[instructor.DepartmentID] = make([]Instructors.Instructor, 1, 8)
			department_id_to_instructors[instructor.DepartmentID][0] = instructor
		} else {
			department_id_to_instructors[instructor.DepartmentID] = append(
				department_id_to_instructors[instructor.DepartmentID], instructor,
			)
		}
	}

	return department_id_to_instructors, nil
}

func GenerateMapDeptIdToDepartment(persistence *StorageResources.Persistence) (map[uint16]Departments.Department, error) {
	department_id_to_department := make(map[uint16]Departments.Department)

	departments, err := persistence.ReaderService.GetAllDepartments()

	if err != nil {
		return nil, err
	}

	for _, department := range departments {
		department_id_to_department[department.DepartmentID] = department
	}

	return department_id_to_department, nil
}

////////////////////////////////////////////////////////////////////////////////////
//              CHECK IF THERE IS ENOUGH INSTRUCTORS FOR THE SCHEDULES
////////////////////////////////////////////////////////////////////////////////////

// The minimum recommended difference between the total available room hours in a department
// and the total lecture hours (or the total laboratory hours) for all subjects in the department.
// This ensures that schedules can be generated with minimal risk of resource shortages.
//
// Tested Values :
//
// * MIN_SUBJECT_ROOM_HOUR_BUFFER = 200, MIN_SUBJECT_INSTRUCTOR_HOUR_BUFFER = 300
//
// : 1024 generations => 65% - 68% valid schedules.
const MIN_SUBJECT_ROOM_HOUR_BUFFER int = 300

// The minimum recommended difference between the total available instructor hours in a department
// and the total lecture and laboratory hours combined for all subjects in the department.
// This ensures that schedules can be generated with minimal risk of resource shortages.
const MIN_SUBJECT_INSTRUCTOR_HOUR_BUFFER int = 300

type Totals struct {
	DepartmentID uint16
	SectionCount int

	LecRoomHours int
	LabRoomHours int
	GymRoomHours int

	InstructorHours float64
	InstructorCount int

	SubjectLecHours int
	SubjectLabHours int
	SubjectGymHours int

	RoomCapacity int
	RoomLabCount int
	RoomLecCount int

	Semester        int
	DepartmentName  string
	Courses         int
	LecSubjectCount int
	LabSubjectCount int
}

func EstimateResourceAvailability(persistence *StorageResources.Persistence, selected_semester, distribution_type int) []error {

	curriculums, err_curriculum := persistence.ReaderService.GetAllCurriculum()

	if err_curriculum != nil {
		list_of_returned_errors := make([]error, 0, 2)
		list_of_returned_errors = append(list_of_returned_errors, err_curriculum)
		return list_of_returned_errors
	}

	department_id_to_department, err_department_id_to_department := GenerateMapDeptIdToDepartment(persistence)

	if err_department_id_to_department != nil {
		error_slice := make([]error, 0, 2)
		error_slice = append(error_slice, err_department_id_to_department)
		return error_slice
	}

	instructors, err_instructors := persistence.ReaderService.GetAllInstructors()

	if err_instructors != nil {
		list_of_returned_errors := make([]error, 0, 2)
		list_of_returned_errors = append(list_of_returned_errors, err_instructors)
		return list_of_returned_errors
	}

	rooms, err_rooms := persistence.ReaderService.GetAllRooms()

	if err_rooms != nil {
		list_of_returned_errors := make([]error, 0, 2)
		list_of_returned_errors = append(list_of_returned_errors, err_rooms)
		return list_of_returned_errors
	}

	totals := make(map[uint16]*Totals)

	for _, room := range rooms {
		_, has_department_id := totals[room.DepartmentID]

		if !has_department_id {
			if room.RoomType == Rooms.ROOM_TYPE_LEC {
				totals[room.DepartmentID] = &Totals{
					LecRoomHours: Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity),
					RoomCapacity: int(room.Capacity),
					RoomLecCount: 1,
				}
			} else if room.RoomType == Rooms.ROOM_TYPE_LAB {
				totals[room.DepartmentID] = &Totals{
					LabRoomHours: Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity),
					RoomCapacity: int(room.Capacity),
					RoomLabCount: 1,
				}
			} else if room.RoomType == Rooms.ROOM_TYPE_GYM {
				totals[room.DepartmentID] = &Totals{
					GymRoomHours: Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity),
					RoomCapacity: int(room.Capacity),
				}
			}
		} else {
			if room.RoomType == Rooms.ROOM_TYPE_LEC {
				totals[room.DepartmentID].LecRoomHours += (Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity))
				totals[room.DepartmentID].RoomLecCount += int(room.Capacity)
			} else if room.RoomType == Rooms.ROOM_TYPE_LAB {
				totals[room.DepartmentID].LabRoomHours += (Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity))
				totals[room.DepartmentID].RoomLabCount += int(room.Capacity)
			} else if room.RoomType == Rooms.ROOM_TYPE_GYM {
				totals[room.DepartmentID].GymRoomHours += (Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity))
			}

			totals[room.DepartmentID].RoomCapacity += int(room.Capacity)
		}
	}

	for _, instructor := range instructors {
		_, has_department_id := totals[instructor.DepartmentID]

		if !has_department_id {
			totals[instructor.DepartmentID] = &Totals{
				InstructorCount: 1,
			}
		} else {
			totals[instructor.DepartmentID].InstructorCount++
		}

		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				if instructor.Time.GetAvailability(day, time_slot) {
					_, ok := totals[instructor.DepartmentID]

					if !ok {
						totals[instructor.DepartmentID] = &Totals{
							InstructorHours: 1.0 / Const.N_HOUR_TIME_SLOTS,
						}
					} else {
						totals[instructor.DepartmentID].InstructorHours += (1.0 / Const.N_HOUR_TIME_SLOTS)
					}
				}
			}
		}
	}

	for _, curriculum := range curriculums {
		_, has_department_id := totals[curriculum.DepartmentID]

		if !has_department_id {
			totals[curriculum.DepartmentID] = &Totals{
				DepartmentID: curriculum.DepartmentID,
				Courses:      1,
			}
		} else {
			totals[curriculum.DepartmentID].DepartmentID = curriculum.DepartmentID
			totals[curriculum.DepartmentID].Courses++

			if len(totals[curriculum.DepartmentID].DepartmentName) == 0 {
				totals[curriculum.DepartmentID].DepartmentName = department_id_to_department[curriculum.DepartmentID].Name
			}
		}

		for _, year_level := range curriculum.YearLevels {
			for semester_idx, semester := range year_level.Semesters {
				if semester_idx == selected_semester {
					total_lec_hours_for_course_level := 0
					total_lab_hours_for_course_level := 0
					total_gym_hours_for_course_level := 0

					for _, subject := range semester.Subjects {
						if !subject.IsGymType() {
							total_lec_hours_for_course_level += int(subject.LecHours)
						} else {
							total_gym_hours_for_course_level += int(subject.LecHours)
						}

						total_lab_hours_for_course_level += int(subject.LabHours)

						if subject.LecHours > 0 {
							totals[curriculum.DepartmentID].LecSubjectCount++
						}

						if subject.LabHours > 0 {
							totals[curriculum.DepartmentID].LabSubjectCount++
						}

					}

					total_lec_hours_for_course_level *= semester.Sections
					total_lab_hours_for_course_level *= semester.Sections
					total_gym_hours_for_course_level *= semester.Sections

					totals[curriculum.DepartmentID].Semester = semester_idx
					totals[curriculum.DepartmentID].SubjectLecHours += total_lec_hours_for_course_level
					totals[curriculum.DepartmentID].SubjectLabHours += total_lab_hours_for_course_level
					totals[curriculum.DepartmentID].SubjectGymHours += total_gym_hours_for_course_level
					totals[curriculum.DepartmentID].SectionCount += semester.Sections
				}
			}
		}
	}

	Utils.PrettyPrint(totals)

	list_of_returned_errors := make([]error, 0, 8)

	for k, v := range totals {
		if k != 0 && (v.SubjectLecHours > 0 || v.SubjectLabHours > 0 || v.SubjectGymHours > 0) {
			if (v.SubjectLecHours + MIN_SUBJECT_ROOM_HOUR_BUFFER) > v.LecRoomHours {
				recommended_room_hours := v.SubjectLecHours + MIN_SUBJECT_ROOM_HOUR_BUFFER
				needed_hours := (recommended_room_hours - v.LecRoomHours)

				rooms_to_add := needed_hours / (Const.N_DAILY_SCHOOL_HOURS * Const.N_WEEKLY_SCHOOL_DAYS)

				if (needed_hours % (Const.N_DAILY_SCHOOL_HOURS * Const.N_WEEKLY_SCHOOL_DAYS)) > 0 {
					rooms_to_add++
				}

				list_of_returned_errors = append(list_of_returned_errors,
					fmt.Errorf(`{"Msg":`+
						`"not enough lecture rooms for the '%s', need %d more capacity for lecture subjects"}`,
						v.DepartmentName, rooms_to_add,
					),
				)
			}

			if (v.SubjectLabHours + MIN_SUBJECT_ROOM_HOUR_BUFFER) > v.LabRoomHours {
				recommended_room_hours := v.SubjectLabHours + MIN_SUBJECT_ROOM_HOUR_BUFFER
				needed_hours := (recommended_room_hours - v.LabRoomHours)

				rooms_to_add := needed_hours / (Const.N_DAILY_SCHOOL_HOURS * Const.N_WEEKLY_SCHOOL_DAYS)

				if (needed_hours % (Const.N_DAILY_SCHOOL_HOURS * Const.N_WEEKLY_SCHOOL_DAYS)) > 0 {
					rooms_to_add++
				}

				list_of_returned_errors = append(list_of_returned_errors,
					fmt.Errorf(`{"Msg":`+
						`"not enough laboratory rooms for the '%s', need %d more capacity for laboratory subjects"}`,
						v.DepartmentName, rooms_to_add,
					),
				)
			}

			if (v.SubjectGymHours + MIN_SUBJECT_ROOM_HOUR_BUFFER) > totals[0].GymRoomHours {
				recommended_room_hours := v.SubjectGymHours + MIN_SUBJECT_ROOM_HOUR_BUFFER
				needed_hours := (recommended_room_hours - totals[0].GymRoomHours)

				rooms_to_add := needed_hours / (Const.N_DAILY_SCHOOL_HOURS * Const.N_WEEKLY_SCHOOL_DAYS)

				if (needed_hours % (Const.N_DAILY_SCHOOL_HOURS * Const.N_WEEKLY_SCHOOL_DAYS)) > 0 {
					rooms_to_add++
				}

				list_of_returned_errors = append(list_of_returned_errors,
					fmt.Errorf(`{"Msg":`+
						`"not enough gym capacity, need %d more capacity for gym subjects"}`,
						rooms_to_add,
					),
				)
			}

			if (v.SubjectLecHours + v.SubjectLabHours + MIN_SUBJECT_INSTRUCTOR_HOUR_BUFFER) > int(math.Round(v.InstructorHours)) {
				recommended_instructor_hours := v.SubjectLecHours + v.SubjectLabHours + MIN_SUBJECT_INSTRUCTOR_HOUR_BUFFER
				needed_hours := (recommended_instructor_hours - int(math.Floor(v.InstructorHours)))

				needed_full_time_instructor := needed_hours / (Const.N_DAILY_SCHOOL_HOURS * Const.N_WEEKLY_SCHOOL_DAYS)

				if (needed_hours % (Const.N_DAILY_SCHOOL_HOURS * Const.N_WEEKLY_SCHOOL_DAYS)) > 0 {
					needed_full_time_instructor++
				}

				list_of_returned_errors = append(list_of_returned_errors,
					fmt.Errorf(`{"Msg":`+
						`"not enough instructors for the '%s', need %d more full time instructors, `+
						`to fill up the missing %d hours of duty, or encourage existing multiple `+
						`instructors to extend their work time"`,
						v.DepartmentName, needed_full_time_instructor, needed_hours,
					),
				)
			}
		}
	}

	return list_of_returned_errors
}
