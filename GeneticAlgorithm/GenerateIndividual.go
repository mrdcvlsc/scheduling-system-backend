package GeneticAlgorithm

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
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

type ScheduleIndex struct {
	DepartmentID uint16
	CurriculumID uint16
	YearLevel    uint8
	Section      uint8
}

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
func NewIndividual(persistence_resources *StorageResources.Persistence, selected_semester, distribution_type int) (Schedule.UniTimeTables, error) {

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_room_type_to_rooms, err_dept_id_to_room_type_to_rooms := generate_map_dept_id_to_room_type_to_rooms(persistence_resources)

	if err_dept_id_to_room_type_to_rooms != nil {
		return nil, err_dept_id_to_room_type_to_rooms
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_instructors, err_dept_id_to_instructors := generate_map_dept_id_to_instructors(persistence_resources)

	if err_dept_id_to_instructors != nil {
		return nil, err_dept_id_to_instructors
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_department, err_dept_id_to_department := generate_map_dept_id_to_department(persistence_resources)

	if err_dept_id_to_department != nil {
		return nil, err_dept_id_to_department
	}

	////////////////////////////////////////////////////////////////////////////////////////

	counted_sections := 0
	individual_university_schedules := make(Schedule.UniTimeTables, 0, 64)

	curriculums, err_curriculums := persistence_resources.ReaderService.GetAllCurriculum()

	if err_curriculums != nil {
		return nil, err_curriculums
	}

	// the curriculums should always have the same order.

	for _, curriculum := range curriculums {

		room_type_to_rooms := dept_id_to_room_type_to_rooms[curriculum.DepartmentID]
		instructors := dept_id_to_instructors[curriculum.DepartmentID]

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
				var section_idx int

			section_loop:
				for section_idx = 0; section_idx < semester.Sections; section_idx++ {

					week_time_table := Schedule.WeekTimeTable{}

					// shuffle the rooms

					for _, rooms := range room_type_to_rooms {
						rng.Shuffle(len(rooms), func(i, j int) {
							rooms[i], rooms[j] = rooms[j], rooms[i]
						})
					}

					// shuffle the subjects

					rng.Shuffle(len(semester.Subjects), func(i, j int) {
						semester.Subjects[i], semester.Subjects[j] = semester.Subjects[j], semester.Subjects[i]
					})

					// TODO: implement distribution types
					//
					// front compressed distribution : start
					//
					// for now the this is a brute force implementation, yet it is still enough
					// considering the current low numbers of courses, instructors, rooms and sections.
					// but I believe this still can be optimize in the future if we wanted to.
					//
					// or just create another implementation for inidividual generation.

					// iterate over the subjects

					// TODO: the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

					subject_recorder := make(map[uint16]Curriculum.Subject)

					// fmt.Println("=============================================") // DEBUG SUBJECT UNASSIGNED PROBLEM
					for _, subject := range semester.Subjects {
						// fmt.Printf("assigning subject : %s", subject.Code) // DEBUG SUBJECT UNASSIGNED PROBLEM

						var selected_instructor *Instructors.Instructor

						// shuffle instructors

						rng.Shuffle(len(instructors), func(i, j int) {
							instructors[i], instructors[j] = instructors[j], instructors[i]
						})

						// sort the instructors based on the number of subjects they are assigned

						sort.Slice(instructors, func(i, j int) bool {
							return instructors[i].TotalTeachingHours < instructors[j].TotalTeachingHours
						})

						// iterate over the class type of the subject lec = 0 or lab = 1

						rand_class_type := int(rng.Int31n(2))

						for class_type_iter := 0; class_type_iter < 2; class_type_iter++ {

							class_type := (rand_class_type + class_type_iter) % 2

							// fmt.Printf("\tcI[%d], cT[%d]\t", class_type_iter, class_type) // DEBUG SUBJECT UNASSIGNED PROBLEM

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

								day_sched := week_time_table.GetDayTimeTable(day)

								for time_slot := 0; time_slot < (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots); time_slot++ {

									if !day_sched.IsTimeAvailable(time_slot, subject_total_time_slots) && (time_slot >= (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {
										if section_generation_retries < MAX_SECTION_SCHEDULE_GENERATION_RETRY {
											section_generation_retries++
											section_idx--
											continue section_loop
										}

										return individual_university_schedules, fmt.Errorf(
											"no time slot found for %s in %s for %s %s %s section[%d] after generating schedules for %d other sections",
											subject.Code, dept_id_to_department[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section_idx, counted_sections,
										)
									}

									if !day_sched.IsTimeAvailable(time_slot, subject_total_time_slots) {
										continue // if the current time slot is not available go to the next
									}

									/////////////////////////////////////////////////////////////////////////////////
									//                 FIND AVAILABLE INSTRUCTOR FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////

									// if the time slot is available proceed to assign an available instructor,
									// check if there is already an assigned instructor for the current subject

									instructor_search_iteration := 0

									// fmt.Printf("instructor & room searching for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

									selected_instructor_idx := -1
									var is_available_instructor bool

									for instructor_idx := range instructors {

										is_available_instructor = true

										if selected_instructor == nil {
											// fmt.Printf("searching the available time slot for the iterated instructor [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
											for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
												is_available_instructor = is_available_instructor && instructors[instructor_idx].Time.GetAvailability(day, instructor_time_slot)
											}

										} else {
											// fmt.Printf("searching the available time slot for the selected instructor [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
											for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
												is_available_instructor = is_available_instructor && selected_instructor.Time.GetAvailability(day, instructor_time_slot)
											}
										}

										instructor_search_iteration++

										if (!is_available_instructor && (instructor_idx == len(instructors)-1)) && (time_slot >= (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {

											if section_generation_retries < MAX_SECTION_SCHEDULE_GENERATION_RETRY {
												section_generation_retries++
												section_idx--
												continue section_loop
											}

											return individual_university_schedules, fmt.Errorf(
												"not enough instructors (%d) in %s for %s %s %s section[%d] after generating schedules for %d other sections",
												instructor_idx, dept_id_to_department[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section_idx, counted_sections,
											)
										}

										if !is_available_instructor && selected_instructor != nil {
											break // immediately find other time slots if there is already a selected instructor yet is not available
										}

										if !is_available_instructor {
											// fmt.Printf("No instructor found available for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
											continue // find another instructor if not available for the time slot
										}

										// fmt.Printf("instructor found available for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

										selected_instructor_idx = instructor_idx
										break
									}

									if !is_available_instructor {
										continue
									}

									/////////////////////////////////////////////////////////////////////////////////
									//                 FIND AVAILABLE ROOM FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////

									room_search_iteration := 0
									room_type := uint16(class_type)

									var has_available_room bool

									if subject.IsGymType() {
										// fmt.Printf("searching available gym rooms for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

										// search available gym for physical education subjects

										gym := dept_id_to_room_type_to_rooms[0][2]

										for room_idx := range gym {

											has_available_room = true

											for room_time_slot := time_slot; room_time_slot < (time_slot + subject_total_time_slots); room_time_slot++ {
												has_available_room = has_available_room && gym[room_idx].GetTimeSlotClassCount(day, room_time_slot) < uint8(gym[room_idx].Capacity)
											}

											if !has_available_room {
												continue
											}

											// fmt.Printf("selecting the available gym for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
											selected_room = &gym[room_idx]
											break
										}
									} else {

										// fmt.Printf("searching available room[type:%d] for the time slot [d:%d, ts:%d]...\n", room_type, day, time_slot) // DEBUG PRINTS

										// search for department specific rooms that are available

										for room_idx := range room_type_to_rooms[room_type] {

											has_available_room = true

											for room_time_slot := time_slot; room_time_slot < (time_slot + subject_total_time_slots); room_time_slot++ {
												has_available_room = has_available_room && room_type_to_rooms[room_type][room_idx].GetTimeSlotClassCount(day, room_time_slot) < uint8(room_type_to_rooms[room_type][room_idx].Capacity)
											}

											if !has_available_room {
												continue
											}

											// fmt.Printf("selecting the available room[type:%d] for the time slot [d:%d, ts:%d]...\n", room_type, day, time_slot) // DEBUG PRINTS
											selected_room = &room_type_to_rooms[room_type][room_idx]
											break
										}

										// TODO: [implement below] search for general rooms that are available (consult first)

										// TODO: [implement below] search available lab room for lecture subjects (consult first)
									}

									room_search_iteration++

									if !has_available_room && (time_slot >= (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {

										if section_generation_retries < MAX_SECTION_SCHEDULE_GENERATION_RETRY {
											section_generation_retries++
											section_idx--
											continue section_loop
										}

										return individual_university_schedules, fmt.Errorf(
											"not enough rooms (%d)-(type:%d) in %s for %s %s %s section[%d] after generating schedules for %d other sections",
											len(room_type_to_rooms[room_type]), room_type, dept_id_to_department[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section_idx, counted_sections,
										)
									}

									if !has_available_room {
										// fmt.Printf("No room found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
										continue // find another time slot
									}

									// fmt.Printf("room found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

									if selected_instructor == nil {
										// fmt.Printf("selecting the instructor found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
										selected_instructor = &instructors[selected_instructor_idx]
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
										selected_instructor.Time.SetAvailability(false, day, selected_time_slot)
										selected_room.IncTimeSlotClassCount(day, selected_time_slot)

										day_sched.GetTimeSlot(selected_time_slot).SetSubjectID(subject.ID)
										day_sched.GetTimeSlot(selected_time_slot).SetInstructorID(selected_instructor.InstructorID)
										day_sched.GetTimeSlot(selected_time_slot).SetRoomID(selected_room.RoomID)
									}

									selected_instructor.AssignedSubjects++
									selected_instructor.TotalTeachingHours += subject_hours

									subject_recorder[subject.ID] = subject

									// fmt.Printf("\t\tassigned : %s", subject.Code) // DEBUG SUBJECT UNASSIGNED PROBLEM

									/////////////////////////////////////////////////////////////////////////////////
									//                 BREAK day AND time_slot LOOP
									/////////////////////////////////////////////////////////////////////////////////

									section_generation_retries = 0
									day = 9999
									time_slot = 9999
								} // ------------- end of time_slot loop -------------
							} // ------------- end of day loop -------------
						} // ------------- end of class_type_iter loop -------------
						// fmt.Printf("\t\t<-- result for : %s\n", subject.Code) // DEBUG SUBJECT UNASSIGNED PROBLEM
					} // ------------- end of subject loop -------------

					// front compressed distribution : end

					if len(subject_recorder) != len(semester.Subjects) {

						// fmt.Printf("\n\nAssigned Subjects : %s %s %s section %d\n", curriculum.CurriculumName, year_level.Name, semester.Name, section_idx) // DEBUG SUBJECT UNASSIGNED PROBLEM

						// for _, s := range subject_recorder { // DEBUG SUBJECT UNASSIGNED PROBLEM
						// fmt.Printf("%v\n", s) // DEBUG SUBJECT UNASSIGNED PROBLEM
						// } // DEBUG SUBJECT UNASSIGNED PROBLEM

						// fmt.Printf("\n\nSubjects To Assign : %s %s %s section %d\n", curriculum.CurriculumName, year_level.Name, semester.Name, section_idx) // DEBUG SUBJECT UNASSIGNED PROBLEM

						// for _, s := range semester.Subjects { // DEBUG SUBJECT UNASSIGNED PROBLEM
						// fmt.Printf("%v\n", s) // DEBUG SUBJECT UNASSIGNED PROBLEM
						// } // DEBUG SUBJECT UNASSIGNED PROBLEM

						panic(fmt.Sprintf(
							"there are some subjects that was not assigned for some reason (%d/%d)", len(subject_recorder), len(semester.Subjects),
						))
					}

					individual_university_schedules = append(individual_university_schedules, week_time_table)
					counted_sections++
				} // ------------- end of section_idx loop -------------
			} // ------------- end of semester_idx loop -------------
		} // ------------- end of year_level loop -------------
	} // ------------- end of curriculum loop -------------

	return individual_university_schedules, nil
}

// TODO: when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying
// each sections in the generated university schedules, make the generated schedule and index global for access.
