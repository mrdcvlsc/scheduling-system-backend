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

const MAX_SECTION_SCHEDULE_GENERATION_RETRY int = 3

// Generate individual university schedules.
//
// Different return types:
//
// (nil, error) - if this function returns a (nil, error) that would mean there is an
// error that prevented the function to read the required data resources.
//
// (slice, error) - if this function returns a (slice, error) that would mean that it produced one
// invalid section schedule due to not having enough resources available during
// the schedule generation configuration.
//
// (slice, nil) - if this function returns a (slice, nil), that would mean it
// successfully generated a valid university schedules.
func NewIndividual(selected_semester, distribution_type int) (Schedule.UniTimeTables, error) {

	persistence := Storage.PersistenceService{ReaderService: &Storage.JsonFilePersistence{}}
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	////////////////////////////////////////////////////////////////////////////////////////

	list_of_all_rooms, err_all_rooms := generate_map_list_of_all_rooms(&persistence)

	if err_all_rooms != nil {
		return nil, err_all_rooms
	}

	////////////////////////////////////////////////////////////////////////////////////////

	list_of_all_instructors, err_all_instructors := generate_map_list_instructors(&persistence)

	if err_all_instructors != nil {
		return nil, err_all_instructors
	}

	////////////////////////////////////////////////////////////////////////////////////////

	map_department_id, err_map_department_id := generate_map_department_id(&persistence)

	if err_map_department_id != nil {
		return nil, err_map_department_id
	}

	////////////////////////////////////////////////////////////////////////////////////////

	counted_sections := 0
	individual := make(Schedule.UniTimeTables, 0, 64)

	all_curriculums, err_all_curriculums := persistence.ReaderService.GetAllCurriculum()

	if err_all_curriculums != nil {
		return nil, err_all_curriculums
	}

	// the curriculums should always have the same order.

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

				// generate week time table for each sections

				section_generation_retries := 0
				var section int

			section_loop:
				for section = 0; section < semester.Sections; section++ {

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

						rand_class_type := int(rng.Int31n(2))

						for class_type_iter := 0; class_type_iter < 2; class_type_iter++ {

							class_type := (rand_class_type + class_type_iter) % 2

							var selected_room *Rooms.Room

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

									// shuffle instructors

									rng.Shuffle(len(dept_teachers), func(i, j int) {
										dept_teachers[i], dept_teachers[j] = dept_teachers[j], dept_teachers[i]
									})

									// sort the instructors based on the number of subjects they are assigned

									sort.Slice(dept_teachers, func(i, j int) bool {
										return dept_teachers[i].TotalTeachingHours < dept_teachers[j].TotalTeachingHours
									})

									// fmt.Printf("instructor & room searching for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

									for instructor_idx := range dept_teachers {

										is_available_instructor := true

										if selected_instructor == nil {

											// fmt.Printf("searching the available time slot for the iterated instructor [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

											// IF NONE: iterate over all of the sorted instructors to find which
											// one is available, if there is no instructor available for the current
											// time slot, continue to the next iteration of the time slot loop.

											for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
												is_available_instructor = is_available_instructor && dept_teachers[instructor_idx].TimeSlotAvailability.GetAvailability(day, instructor_time_slot)
											}

										} else {
											// fmt.Printf("searching the available time slot for the selected instructor [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

											for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
												is_available_instructor = is_available_instructor && selected_instructor.TimeSlotAvailability.GetAvailability(day, instructor_time_slot)
											}
										}

										instructor_search_iteration++

										if (!is_available_instructor && (instructor_idx == len(dept_teachers)-1)) && (time_slot >= (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {

											if section_generation_retries < MAX_SECTION_SCHEDULE_GENERATION_RETRY {
												section_generation_retries++
												section--
												continue section_loop
											}

											return individual, fmt.Errorf(
												"not enough instructors (%d) in %s for %s %s %s section[%d] after generating schedules for %d other sections",
												instructor_idx, map_department_id[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section, counted_sections,
											)
										}

										if !is_available_instructor {
											// fmt.Printf("No instructor found available for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
											continue // find another instructor if not available for the time slot
										}

										// fmt.Printf("instructor found available for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

										/////////////////////////////////////////////////////////////////////////////////
										//                 FIND AVAILABLE ROOM FOR THE TIME SLOT
										/////////////////////////////////////////////////////////////////////////////////

										room_search_iteration := 0
										room_type := uint16(class_type)

										// shuffle rooms

										rng.Shuffle(len(dept_rooms[room_type]), func(i, j int) {
											dept_rooms[room_type][i], dept_rooms[room_type][j] = dept_rooms[room_type][j], dept_rooms[room_type][i]
										})

										// sort the rooms based on the number of class or sections assigned to it on a specific time slot.

										sort.Slice(dept_rooms[room_type], func(i, j int) bool {
											return dept_rooms[room_type][i].GetTimeSlotClassCount(day, time_slot) < dept_rooms[room_type][j].GetTimeSlotClassCount(day, time_slot)
										})

										is_available_room := true

										if subject.IsGymType() {
											// fmt.Printf("searching available gym rooms for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

											// search available gym for physical education subjects

											gym := list_of_all_rooms[0][2]

											for room_idx := range gym {

												for room_time_slot := time_slot; room_time_slot < (time_slot + subject_total_time_slots); room_time_slot++ {
													is_available_room = is_available_room && gym[room_idx].GetTimeSlotClassCount(day, room_time_slot) < uint8(gym[room_idx].Capacity)
												}

												if !is_available_room {
													continue
												}

												// fmt.Printf("selecting the available gym for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
												selected_room = &gym[room_idx]
												break
											}
										} else {

											// fmt.Printf("searching available room[type:%d] for the time slot [d:%d, ts:%d]...\n", room_type, day, time_slot) // DEBUG PRINTS

											// search for department specific rooms that are available

											for room_idx := range dept_rooms[room_type] {

												for room_time_slot := time_slot; room_time_slot < (time_slot + subject_total_time_slots); room_time_slot++ {
													is_available_room = is_available_room && dept_rooms[room_type][room_idx].GetTimeSlotClassCount(day, room_time_slot) < uint8(dept_rooms[room_type][room_idx].Capacity)
												}

												if !is_available_room {
													continue
												}

												// fmt.Printf("selecting the available room[type:%d] for the time slot [d:%d, ts:%d]...\n", room_type, day, time_slot) // DEBUG PRINTS
												selected_room = &dept_rooms[room_type][room_idx]
												break
											}

											// TODO: [implement below] search for general rooms that are available (consult first)

											// TODO: [implement below] search available lab room for lecture subjects (consult first)
										}

										room_search_iteration++

										if !is_available_room && (time_slot >= (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {

											if section_generation_retries < MAX_SECTION_SCHEDULE_GENERATION_RETRY {
												section_generation_retries++
												section--
												continue section_loop
											}

											return individual, fmt.Errorf(
												"not enough rooms (%d)-(type:%d) () in %s for %s %s %s section[%d] after generating schedules for %d other sections",
												len(dept_rooms[room_type]), room_type, map_department_id[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section, counted_sections,
											)
										}

										if !is_available_room {
											// fmt.Printf("No room found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
											continue // find another instructor or time slot
										}

										// fmt.Printf("room found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

										// the instructor loop will not reach here if there are no available instructors and rooms found

										if selected_instructor == nil {
											// fmt.Printf("selecting the instructor found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
											selected_instructor = &dept_teachers[instructor_idx]
										}

										break
									}

									if selected_instructor == nil || selected_room == nil {
										continue
									}

									/////////////////////////////////////////////////////////////////////////////////
									//                 ASSING AVAILABLE INSTRUCTOR FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////

									// fmt.Printf( // DEBUG PRINTS
									// "[%s]-[%s]-Section:[%d] | [%s][hours(%d):%d] : day(%d):timeslot(%d) | Instructor : (%s %s %s) found after %d iterations\n", // DEBUG PRINTS
									// curriculum.CurriculumCode,         // DEBUG PRINTS
									// year_level.Name,                   // DEBUG PRINTS
									// section,                           // DEBUG PRINTS
									// subject.Code,                      // DEBUG PRINTS
									// subject_hours,                     // DEBUG PRINTS
									// class_type,                        // DEBUG PRINTS
									// day,                               // DEBUG PRINTS
									// time_slot,                         // DEBUG PRINTS
									// selected_instructor.FirstName,     // DEBUG PRINTS
									// selected_instructor.MiddleInitial, // DEBUG PRINTS
									// selected_instructor.LastName,      // DEBUG PRINTS
									// instructor_search_iteration,       // DEBUG PRINTS
									// ) // DEBUG PRINTS

									// fmt.Printf("ONE SUBJECT FINISHED: assigning subject, instructor and room for the time slot [d:%d, ts:%d]...\n\n", day, time_slot) // DEBUG PRINTS

									for selected_time_slot := time_slot; selected_time_slot < (time_slot + subject_total_time_slots); selected_time_slot++ {
										selected_instructor.TimeSlotAvailability.SetAvailability(false, day, selected_time_slot)
										selected_room.IncTimeSlotClassCount(day, selected_time_slot)

										day_sched.Get(selected_time_slot).SetSubjectID(subject.ID)
										day_sched.Get(selected_time_slot).SetInstructorID(selected_instructor.InstructorID)
										day_sched.Get(selected_time_slot).SetRoomID(selected_room.RoomID)
									}

									selected_instructor.AssignedSubjects++
									selected_instructor.TotalTeachingHours += subject_hours

									/////////////////////////////////////////////////////////////////////////////////
									//                 BREAK day AND time_slot LOOP
									/////////////////////////////////////////////////////////////////////////////////

									section_generation_retries = 0
									day = 9999
									time_slot = 9999
								}
							}
						}
					}

					// front compressed distribution : end

					individual = append(individual, week_time_table)
					counted_sections++
				}
			}
		}

		// sort the instructors based on the number of subjects they are assigned
		sort.Slice(dept_teachers, func(i, j int) bool {
			return dept_teachers[i].TotalTeachingHours < dept_teachers[j].TotalTeachingHours
		})
	}

	return individual, nil
}
