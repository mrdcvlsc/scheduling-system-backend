package geneticalgorithm

import (
	"fmt"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Storage"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func generate_map_list_of_all_rooms(persistence *Storage.PersistenceService) map[uint16]map[uint16][]Rooms.Room {

	list_of_all_rooms := make(map[uint16]map[uint16][]Rooms.Room)

	all_rooms := persistence.Service.GetAllRooms()

	for _, room := range all_rooms {
		_, dept_key_exist := list_of_all_rooms[room.DepartmentID]

		if !dept_key_exist {
			list_of_all_rooms[room.DepartmentID] = make(map[uint16][]Rooms.Room)
		}

		_, roomtype_key_exist := list_of_all_rooms[room.DepartmentID][room.RoomType]

		if !roomtype_key_exist {
			list_of_all_rooms[room.DepartmentID][room.RoomType] = make([]Rooms.Room, 1, 8)
			list_of_all_rooms[room.DepartmentID][room.RoomType][0] = room
		} else {
			list_of_all_rooms[room.DepartmentID][room.RoomType] = append(
				list_of_all_rooms[room.DepartmentID][room.RoomType], room,
			)
		}
	}

	return list_of_all_rooms
}

func generate_map_list_instructors(persistence *Storage.PersistenceService) map[uint16][]Instructors.Instructor {
	list_of_all_instructors := make(map[uint16][]Instructors.Instructor)

	all_instructors := persistence.Service.GetAllInstructors()

	for _, instructor := range all_instructors {
		_, exist := list_of_all_instructors[instructor.DepartmentID]

		if !exist {
			list_of_all_instructors[instructor.DepartmentID] = make([]Instructors.Instructor, 1, 8)
			list_of_all_instructors[instructor.DepartmentID][0] = instructor
		} else {
			list_of_all_instructors[instructor.DepartmentID] = append(
				list_of_all_instructors[instructor.DepartmentID], instructor,
			)
		}
	}

	return list_of_all_instructors
}

////////////////////////////////////////////////////////////////////////////////////
//              CHECK IF THERE IS ENOUGH INSTRUCTORS FOR THE SCHEDULES
////////////////////////////////////////////////////////////////////////////////////

type Totals struct {
	TotalLecRoomHours    int
	TotalLabRoomHours    int
	TotalGymRoomHours    int
	TotalInstructorHours float64
	TotalLecHours        int
	TotalLabHours        int
	TotalGymHours        int
	TotalSections        int
	TotalInstructors     int
	TotalSectionCapRooms int
}

func EstimateResourceAvailability(selected_semester, distribution_type int) {
	persistence := Storage.PersistenceService{Service: &Storage.JsonFilePersistence{}}

	curriculums := persistence.Service.GetAllCurriculum()
	instructors := persistence.Service.GetAllInstructors()
	rooms := persistence.Service.GetAllRooms()

	totals := make(map[uint16]*Totals)

	for _, room := range rooms {
		_, ok := totals[room.DepartmentID]

		if !ok {
			if room.RoomType == Rooms.ROOM_TYPE_LEC {
				totals[room.DepartmentID] = &Totals{
					TotalLecRoomHours:    Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity),
					TotalSectionCapRooms: int(room.Capacity),
				}
			} else if room.RoomType == Rooms.ROOM_TYPE_LAB {
				totals[room.DepartmentID] = &Totals{
					TotalLabRoomHours:    Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity),
					TotalSectionCapRooms: int(room.Capacity),
				}
			} else if room.RoomType == Rooms.ROOM_TYPE_GYM {
				totals[room.DepartmentID] = &Totals{
					TotalGymRoomHours:    Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity),
					TotalSectionCapRooms: int(room.Capacity),
				}
			}
		} else {
			if room.RoomType == Rooms.ROOM_TYPE_LEC {
				totals[room.DepartmentID].TotalLecRoomHours += (Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity))
			} else if room.RoomType == Rooms.ROOM_TYPE_LAB {
				totals[room.DepartmentID].TotalLabRoomHours += (Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity))
			} else if room.RoomType == Rooms.ROOM_TYPE_GYM {
				totals[room.DepartmentID].TotalGymRoomHours += (Const.N_WEEKLY_SCHOOL_DAYS * Const.N_DAILY_SCHOOL_HOURS * int(room.Capacity))
			}
		}
	}

	for _, instructor := range instructors {

		_, ok := totals[instructor.DepartmentID]

		if !ok {
			totals[instructor.DepartmentID] = &Totals{
				TotalInstructors: 1,
			}
		} else {
			totals[instructor.DepartmentID].TotalInstructors++
		}

		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				if instructor.TimeSlotAvailability.GetAvailability(day, time_slot) {
					_, ok := totals[instructor.DepartmentID]

					if !ok {
						totals[instructor.DepartmentID] = &Totals{
							TotalInstructorHours: 1.0 / Const.N_HOUR_TIME_SLOTS,
						}
					} else {
						totals[instructor.DepartmentID].TotalInstructorHours += (1.0 / Const.N_HOUR_TIME_SLOTS)
					}
				}
			}
		}
	}

	for _, curriculum := range curriculums {
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
					}

					total_lec_hours_for_course_level *= semester.Sections
					total_lab_hours_for_course_level *= semester.Sections
					total_gym_hours_for_course_level *= semester.Sections

					_, ok := totals[curriculum.DepartmentID]

					if !ok {
						totals[curriculum.DepartmentID] = &Totals{
							TotalLecHours: total_lec_hours_for_course_level,
							TotalLabHours: total_lab_hours_for_course_level,
							TotalGymHours: total_gym_hours_for_course_level,
							TotalSections: semester.Sections,
						}
					} else {
						totals[curriculum.DepartmentID].TotalLecHours += total_lec_hours_for_course_level
						totals[curriculum.DepartmentID].TotalLabHours += total_lab_hours_for_course_level
						totals[curriculum.DepartmentID].TotalGymHours += total_gym_hours_for_course_level
						totals[curriculum.DepartmentID].TotalSections += semester.Sections
					}
				}
			}
		}
	}

	fmt.Print("\n\n================ EstimateResourceAvailability ================ \n\n")
	Utils.PrettyPrint(totals)
}
