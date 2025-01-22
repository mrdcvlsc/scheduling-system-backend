package geneticalgorithm

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Storage"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

const (
	TERM_1ST_SEMESTER int = 0 // `selected_semester` option.
	TERM_2ND_SEMESTER int = 1 // `selected_semester` option.
)

const (
	DIST_FRONT_COMPRESSED int = 0
	DIST_BACK_COMPRESSED  int = 1
	DIST_FRONT_LOOSE      int = 2
	DIST_BACK_LOOSE       int = 3
)

func NewIndividual(selected_semester, distribution_type int) (Schedule.UniTimeTables, error) {

	persistence := Storage.PersistenceService{Service: &Storage.JsonFilePersistence{}}
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	list_of_all_rooms := make(map[uint16]map[uint16][]Rooms.Room)

	{
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
	}

	list_of_all_instructors := make(map[uint16][]Instructors.Instructor)

	{
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
	}

	counted_sections := 0
	individual := make(Schedule.UniTimeTables, 0, 64)

	all_curriculums := persistence.Service.GetAllCurriculum()

	for _, curriculum := range all_curriculums {

		dept_rooms := list_of_all_rooms[curriculum.DepartmentID]
		dept_teachers := list_of_all_instructors[curriculum.DepartmentID]

		for _, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue // skip inactive year levels
			}

			for semester_idx, semester := range year_level.Semesters {

				if selected_semester != semester_idx {
					continue // skip not selected semesters
				}

				counted_sections += semester.Sections

				// generate week time table for each sections

				for section := 0; section < semester.Sections; section++ {

					week_time_table := Schedule.WeekTimeTable{}

					// shuffle the rooms

					for _, rooms_group_by_type := range dept_rooms {
						rng.Shuffle(len(rooms_group_by_type), func(i, j int) {
							rooms_group_by_type[i], rooms_group_by_type[j] = rooms_group_by_type[j], rooms_group_by_type[i]
						})
					}

					// shuffle the subjects

					rng.Shuffle(len(semester.Subjects), func(i, j int) {
						semester.Subjects[i], semester.Subjects[j] = semester.Subjects[j], semester.Subjects[i]
					})

					// front compressed distribution : start
					//
					// for now the this is a brute force implementation, yet it is still enough
					// considering the current low numbers of courses, instructors, rooms and sections.
					// but I believe this still can be optimize in the future if we wanted to.
					//
					// or just create another implementation for inidividual generation.

					// iterate over the subjects

					for _, subject := range semester.Subjects {

						var selected_instructor *Instructors.Instructor

						// iterate over the class type of the subject lec = 0 or lab = 1

						for class_type := 0; class_type < 2; class_type++ {

							var subject_hours int

							if class_type == 0 {
								subject_hours = int(subject.LecHours)
							} else {
								subject_hours = int(subject.LabHours)
							}

							if subject_hours == 0 {
								continue // skip subject class type if there is no contact hours
							}

							subject_total_time_slots := subject_hours * Const.N_HOUR_TIME_SLOTS

							// search an available time slot for the current subject

							for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {

								day_sched := week_time_table.Get(day)

								for time_slot := 0; time_slot < (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots); time_slot++ {

									if !day_sched.Availability(time_slot, subject_total_time_slots) {
										continue // if the current time slot is not available go to the next
									}

									/////////////////////////////////////////////////////////////////////////////////
									//                 FIND AVAILABLE INSTRUCTOR FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////

									// if the time slot is available proceed to assign an available instructor,
									// check if there is already an assigned instructor for the current subject

									instructor_search_iteration := 0

									if selected_instructor == nil {

										// shuffle instructors
										rng.Shuffle(len(dept_teachers), func(i, j int) {
											dept_teachers[i], dept_teachers[j] = dept_teachers[j], dept_teachers[i]
										})

										// sort the instructors based on the number of subjects they are assigned
										sort.Slice(dept_teachers, func(i, j int) bool {
											return dept_teachers[i].TotalTeachingHours < dept_teachers[j].TotalTeachingHours
										})

										// IF NONE: iterate over all of the sorted instructors to find which
										// one is available, if there is no instructor available for the current
										// time slot, continue to the next iteration of the time slot loop.

										for instructor_idx := range dept_teachers {

											is_available_instructor := true

											for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
												is_available_instructor = is_available_instructor && dept_teachers[instructor_idx].TimeSlotAvailability.GetAvailability(day, instructor_time_slot)
											}

											instructor_search_iteration++

											if is_available_instructor {
												selected_instructor = &dept_teachers[instructor_idx]
												break
											}
										}

										if (selected_instructor == nil) && (time_slot == (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day == (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {

											// TODO: if this is the last time slot and there is still no
											// instructor available, throw an error saying there is not
											// enough instructors, true error handling not panic.

											fmt.Printf("\n\n!Panic At The Disco:\n\n")
											Utils.PrettyPrint(dept_teachers)

											panic(fmt.Sprintf(
												"Not Enough Instructors after %d sections for %s %s section-%d, instructor iteration %d",
												counted_sections, curriculum.CurriculumCode, year_level.Name, section, instructor_search_iteration,
											))
										}

										if selected_instructor == nil {
											continue // find another time slot if no instructor is available
										}

									} else {

										is_available_instructor := true

										for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
											is_available_instructor = is_available_instructor && selected_instructor.TimeSlotAvailability.GetAvailability(day, instructor_time_slot)
										}

										if !is_available_instructor {
											// if there is already a selected instructor but that instructor
											// is not available for the current time slot, continue to the
											// next iteration of the time slot loop.
											continue
										}
									}

									/////////////////////////////////////////////////////////////////////////////////
									//                 FIND AVAILABLE ROOM FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////

									/////////////////////////////////////////////////////////////////////////////////
									//                 ASSING AVAILABLE INSTRUCTOR FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////

									fmt.Printf(
										"[%s]-[%s]-Section:[%d] | [%s][hours(%d):%d] : day(%d):timeslot(%d) | Instructor : (%s %s %s) found after %d iterations\n",
										curriculum.CurriculumCode,
										year_level.Name,
										section,
										subject.Code,
										subject_hours,
										class_type,
										day,
										time_slot,
										selected_instructor.FirstName,
										selected_instructor.MiddleInitial,
										selected_instructor.LastName,
										instructor_search_iteration,
									)

									for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
										selected_instructor.TimeSlotAvailability.SetAvailability(false, day, instructor_time_slot)

										// HINT : this might be the final loop where the subjects, instructors and rooms will be assigned.
										day_sched.Get(instructor_time_slot).SetSubjectID(subject.ID)
										day_sched.Get(instructor_time_slot).SetInstructorID(selected_instructor.InstructorID)

									}

									selected_instructor.AssignedSubjects++
									selected_instructor.TotalTeachingHours += subject_hours

									// temporary write assignment data to instructor - this should be done only after all rooms are slected (developement) : start

									/////////////////////////////////////////////////////////////////////////////////
									//                 ASSIGN AVAILABLE ROOM FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////

									// TODO: iterate over all of sorted rooms to find which one is
									// available, if there is no room available for the current time
									// slot, continue to the next iteration of the loop.

									// TODO: if this is the last time slot and there is still no
									// room available, throw an error saying there is not enough room.

									// TODO: after all of the checks and searches, if there are instructors
									// and rooms that are available for a specific subject and time slot,
									// then assign the selected subject, instructor and room to the time slot.

									/////////////////////////////////////////////////////////////////////////////////
									//                 BREAK day AND time_slot LOOP
									/////////////////////////////////////////////////////////////////////////////////

									day = 9999
									time_slot = 9999
								}
							}
						}
					}

					// front compressed distribution : end

					individual = append(individual, week_time_table)
				}
			}
		}

		// sort the instructors based on the number of subjects they are assigned
		sort.Slice(dept_teachers, func(i, j int) bool {
			return dept_teachers[i].TotalTeachingHours < dept_teachers[j].TotalTeachingHours
		})

		// fmt.Print("\n\nDepartment Teachers : \n\n")
		// for _, dept_teacher_iter := range dept_teachers {
		// 	fmt.Printf(
		// 		"Assigned Subjects : %d = %d hours | %s %s. %s | %d\n",
		// 		dept_teacher_iter.AssignedSubjects,
		// 		dept_teacher_iter.TotalTeachingHours,
		// 		dept_teacher_iter.FirstName,
		// 		dept_teacher_iter.MiddleInitial,
		// 		dept_teacher_iter.LastName,
		// 		curriculum.DepartmentID,
		// 	)
		// }
	}

	// fmt.Println("Total Number of Sections : ", total_number_of_sections)

	// TODO: during the first schedule generation throw error if there is not enough instructors

	// TODO: during the first schedule generation throw error if there is not enough rooms

	return individual, nil
}
