package GeneticAlgorithm

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
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

/*
Generate individual university schedules.

Different return types:

	(nil, nil, error) // -> resources copy failed.
	(UniTimeTables, nil, error) // -> not enough resources.
	(UniTimeTables, EncodingResource, nil) // -> successfully generated valid university schedules.

to generate whole university schedules, set department to encode to nil:

	department_to_encode = nil
*/
func EncodeIndividualGenome(
	individual_university_schedules_arg Schedule.UniTimeTables,
	curriculums_arg []Curriculum.Curriculum,
	dept_id_to_department map[uint16]Departments.Department,
	input_encoding_resource *EncodingResource,
	department_to_encode map[uint16]bool,
	selected_semester,
	distribution_type int,
) (Schedule.UniTimeTables, *EncodingResource, error) {

	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums := make([]Curriculum.Curriculum, len(curriculums_arg))
	copied_curriculums := copy(curriculums, curriculums_arg)

	if copied_curriculums != len(curriculums_arg) {
		return nil, nil, fmt.Errorf("slice elements copied %d, internal curriculum copy operation failed in generate new individual function", copied_curriculums)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	individual_university_schedules := make(Schedule.UniTimeTables, len(individual_university_schedules_arg))

	copied_week_time_table := copy(individual_university_schedules, individual_university_schedules_arg)

	if copied_week_time_table != len(individual_university_schedules_arg) {
		return nil, nil, fmt.Errorf("slice elements copied %d, internal university schedule copy operation failed in generate new individual function", copied_week_time_table)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	encoding_resource, err_make_copy := input_encoding_resource.MakeCopy()

	if err_make_copy != nil {
		return nil, nil, err_make_copy
	}

	////////////////////////////////////////////////////////////////////////////////////////

	counted_sections := 0

	for _, curriculum := range curriculums {

		room_type_to_rooms := encoding_resource.DeptIdToRoomtypeToRooms[curriculum.DepartmentID]
		instructors := encoding_resource.DeptIdToInstructors[curriculum.DepartmentID]

		for _, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue // skip inactive year levels
			}

			for semester_idx, semester := range year_level.Semesters {

				if selected_semester != semester_idx {
					continue // skip not selected semesters
				}

				if len(semester.Subjects) == 0 {
					continue // skip semesters that don't have subjects
				}

				/////////////////////////////////////////////////////////////////////////////////////////////////////////
				//                           GENERATE WEEK TIME TABLE FOR EACH SECTIONS
				/////////////////////////////////////////////////////////////////////////////////////////////////////////

				for section_idx := range semester.Sections {

					if department_to_encode != nil {
						is_to_encode := department_to_encode[curriculum.DepartmentID]

						if !is_to_encode {
							counted_sections++
							continue
						}
					}

					week_time_table := individual_university_schedules[counted_sections]

					/////////////////////////////////////////////////////////////////////////////////////////////////////////
					//                                    SHUFFLE ROOMS AND SUBJECT
					/////////////////////////////////////////////////////////////////////////////////////////////////////////

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

					/////////////////////////////////////////////////////////////////////////////////////////////////////////

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

					for _, subject := range semester.Subjects {

						non_final_sched_idx := uint16(counted_sections)

						if _, has_sched_idx := encoding_resource.IsSchedIdxToSubIdToSkip[non_final_sched_idx]; has_sched_idx {
							_, has_sub_id := encoding_resource.IsSchedIdxToSubIdToSkip[non_final_sched_idx][subject.ID]
							if has_sub_id {
								if encoding_resource.IsSchedIdxToSubIdToSkip[non_final_sched_idx][subject.ID] {
									subject_recorder[subject.ID] = subject
									continue // skip since subject was already assigned
								}
							}
						}

						var selected_instructor *Instructors.Instructor

						/////////////////////////////////////////////////////////////////////////////////////////////////////////
						//               POPULATE SPECIALIZED INSTRUCTOR LIST FOR THE SUBJECT IF THEY EXIST
						/////////////////////////////////////////////////////////////////////////////////////////////////////////

						specialized_instructors := make([]*Instructors.Instructor, 0)

						if subject.DesignatedInstructors != nil {
							instructor_id_to_instructor := make(map[uint16]*Instructors.Instructor)

							for i := range instructors {
								instructor_id_to_instructor[instructors[i].InstructorID] = &instructors[i]
							}

							for i := range encoding_resource.DeptIdToInstructors[0] {
								instructor_id_to_instructor[encoding_resource.DeptIdToInstructors[0][i].InstructorID] = &encoding_resource.DeptIdToInstructors[0][i]
							}

							for _, specialized_id := range subject.DesignatedInstructors {
								if _, has_id := instructor_id_to_instructor[specialized_id]; has_id {
									specialized_instructors = append(specialized_instructors, instructor_id_to_instructor[specialized_id])
								}
							}

							if len(specialized_instructors) == 0 {
								return nil, nil, fmt.Errorf(
									"the specialized instructor(s) added in %s %s %s section[%d] subject %s are not found in the department instructors and general instructors",
									curriculum.CurriculumCode, year_level.Name, semester.Name, section_idx, subject.Code,
								)
							}

							// shuffle specialized instructors
							rng.Shuffle(len(specialized_instructors), func(i, j int) {
								specialized_instructors[i], specialized_instructors[j] = specialized_instructors[j], specialized_instructors[i]
							})

							// sort the specialized_instructors based on the number of subjects they are assigned
							sort.Slice(specialized_instructors, func(i, j int) bool {
								return specialized_instructors[i].TotalTeachingHours < specialized_instructors[j].TotalTeachingHours
							})
						} else {
							// shuffle instructors
							rng.Shuffle(len(instructors), func(i, j int) {
								instructors[i], instructors[j] = instructors[j], instructors[i]
							})

							// sort the instructors based on the number of subjects they are assigned
							sort.Slice(instructors, func(i, j int) bool {
								return instructors[i].TotalTeachingHours < instructors[j].TotalTeachingHours
							})
						}

						/////////////////////////////////////////////////////////////////////////////////////////////////////////
						//                         RANDOMIZE LECTURE AND LABORATORY ASSIGNMENT ORDER
						/////////////////////////////////////////////////////////////////////////////////////////////////////////

						// iterate over the class type of the subject lec = 0 or lab = 1
						is_subject_type_added_once := false

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

							/////////////////////////////////////////////////////////////////////////////////////////////////////////
							//                                ITERATE THROUGH THE WEEKLY TIME SLOTS
							/////////////////////////////////////////////////////////////////////////////////////////////////////////

							for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {

								day_sched := week_time_table.GetDayTimeTable(day)

								for time_slot := 0; time_slot <= (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots); time_slot++ {

									/////////////////////////////////////////////////////////////////////////////////////////////////////////
									//                      CHECK IF CURRENT TIME SLOT IS AVAILABLE FOR THE SUBJECT
									/////////////////////////////////////////////////////////////////////////////////////////////////////////

									is_time_slot_available := day_sched.IsTimeAvailable(time_slot, subject_total_time_slots)

									if !is_time_slot_available && (time_slot > (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {
										return individual_university_schedules, nil, fmt.Errorf(
											"no time slot found for %s in %s for %s %s %s section[%d] after generating schedules for %d other sections",
											subject.Code, dept_id_to_department[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section_idx, counted_sections,
										)
									}

									if !is_time_slot_available {
										continue // if the current time slot is not available, go to the next time slot
									}

									/////////////////////////////////////////////////////////////////////////////////////////////////////////
									//                          FIND AVAILABLE INSTRUCTOR FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////////////////////////////

									// if the time slot is available proceed to find then assign an available instructor

									instructor_search_iteration := 0

									selected_instructor_idx := -1
									var is_available_instructor bool

									if subject.DesignatedInstructors != nil {

										/////////////////////////////////////////////////////////////////////////////////////////////////////////
										//                               find available specialized instructors
										/////////////////////////////////////////////////////////////////////////////////////////////////////////

										for instructor_idx := range specialized_instructors {

											is_available_instructor = true

											if selected_instructor == nil {
												for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
													is_available_instructor = is_available_instructor && specialized_instructors[instructor_idx].Time.GetAvailability(day, instructor_time_slot)
												}
											} else {
												for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
													is_available_instructor = is_available_instructor && selected_instructor.Time.GetAvailability(day, instructor_time_slot)
												}
											}

											instructor_search_iteration++

											if (!is_available_instructor && ((instructor_idx == len(specialized_instructors)-1) || selected_instructor != nil)) && (time_slot > (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {
												return individual_university_schedules, nil, fmt.Errorf(
													"not enough specialized_instructors (%d) in %s for %s %s %s section[%d] after generating schedules for %d other sections",
													instructor_idx, dept_id_to_department[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section_idx, counted_sections,
												)
											}

											if !is_available_instructor && selected_instructor != nil {
												break // immediately find other time slots if there is already a selected instructor yet is not available
											}

											if !is_available_instructor {
												continue // find another instructor if not available for the time slot
											}

											selected_instructor_idx = instructor_idx
											break
										}
									} else {
										/////////////////////////////////////////////////////////////////////////////////////////////////////////
										//                                find available department instructors
										/////////////////////////////////////////////////////////////////////////////////////////////////////////

										for instructor_idx := range instructors {

											is_available_instructor = true

											if selected_instructor == nil {
												for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
													is_available_instructor = is_available_instructor && instructors[instructor_idx].Time.GetAvailability(day, instructor_time_slot)
												}

											} else {
												for instructor_time_slot := time_slot; instructor_time_slot < (time_slot + subject_total_time_slots); instructor_time_slot++ {
													is_available_instructor = is_available_instructor && selected_instructor.Time.GetAvailability(day, instructor_time_slot)
												}
											}

											instructor_search_iteration++

											if (!is_available_instructor && ((instructor_idx == len(instructors)-1) || selected_instructor != nil)) && (time_slot > (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {
												return individual_university_schedules, nil, fmt.Errorf(
													"not enough instructors (%d) in %s for %s %s %s section[%d] after generating schedules for %d other sections",
													instructor_idx, dept_id_to_department[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section_idx, counted_sections,
												)
											}

											if !is_available_instructor && selected_instructor != nil {
												break // immediately find other time slots if there is already a selected instructor yet is not available
											}

											if !is_available_instructor {
												continue // find another instructor if not available for the time slot
											}

											selected_instructor_idx = instructor_idx
											break
										}
									}

									// TODO: analyze if this is really needed?

									if !is_available_instructor {
										continue // find other time slot if there is no available instructor
									}

									/////////////////////////////////////////////////////////////////////////////////////////////////////////
									//                             FIND AVAILABLE ROOM FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////////////////////////////

									room_search_iteration := 0
									room_type := uint16(class_type)

									var has_available_room bool

									if subject.IsGymType() {

										/////////////////////////////////////////////////////////////////////////////////////////////////////////
										//                                       find available gym room
										/////////////////////////////////////////////////////////////////////////////////////////////////////////

										// search available gym for physical education subjects

										gym := encoding_resource.DeptIdToRoomtypeToRooms[0][2]

										for room_idx := range gym {

											has_available_room = true

											for room_time_slot := time_slot; room_time_slot < (time_slot + subject_total_time_slots); room_time_slot++ {
												has_available_room = has_available_room && gym[room_idx].GetTimeSlotClassCount(day, room_time_slot) < uint8(gym[room_idx].Capacity)
											}

											if !has_available_room {
												continue
											}

											selected_room = &gym[room_idx]
											break
										}
									} else {

										/////////////////////////////////////////////////////////////////////////////////////////////////////////
										//                                  find available department rooms
										/////////////////////////////////////////////////////////////////////////////////////////////////////////

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

									if !has_available_room && (time_slot > (Const.N_DAILY_TIME_SLOTS - subject_total_time_slots - 1)) && (day >= (Const.N_WEEKLY_SCHOOL_DAYS - 1)) {
										return individual_university_schedules, nil, fmt.Errorf(
											"not enough rooms (%d)-(type:%d) in %s for %s %s %s section[%d] after generating schedules for %d other sections",
											len(room_type_to_rooms[room_type]), room_type, dept_id_to_department[curriculum.DepartmentID].Name, curriculum.CurriculumCode, semester.Name, year_level.Name, section_idx, counted_sections,
										)
									}

									if !has_available_room {
										// fmt.Printf("No room found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS
										continue // find another time slot
									}

									/////////////////////////////////////////////////////////////////////////////////////////////////////////
									//   if there is an available room then finalize instructor selection if there is no one selected yet
									/////////////////////////////////////////////////////////////////////////////////////////////////////////

									// fmt.Printf("room found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

									if selected_instructor == nil {
										// fmt.Printf("selecting the instructor found for the time slot [d:%d, ts:%d]...\n", day, time_slot) // DEBUG PRINTS

										if subject.DesignatedInstructors != nil {
											selected_instructor = specialized_instructors[selected_instructor_idx]
										} else {
											selected_instructor = &instructors[selected_instructor_idx]
										}
									}

									/////////////////////////////////////////////////////////////////////////////////////////////////////////
									//                       ALLOCATE THE FINAL AVAILABLE INSTRUCTOR FOR THE TIME SLOT
									/////////////////////////////////////////////////////////////////////////////////////////////////////////

									time_slot_assignment_sanity_counter := 0

									for selected_time_slot := time_slot; selected_time_slot < (time_slot + subject_total_time_slots); selected_time_slot++ {
										if day_sched.GetTimeSlot(selected_time_slot).GetSubjectID() != 0 {
											panic("woah woah woah! you are overwriting a subject allocated in that time slot")
										}

										if day_sched.GetTimeSlot(selected_time_slot).GetInstructorID() != 0 {
											panic("woah woah woah! you are overwriting an instructor allocated in that time slot")
										}

										if day_sched.GetTimeSlot(selected_time_slot).GetRoomID() != 0 {
											panic("woah woah woah! you are overwriting a room allocated in that time slot")
										}

										selected_instructor.Time.SetAvailability(false, day, selected_time_slot)
										selected_room.IncTimeSlotClassCount(day, selected_time_slot)

										if subject.ID == 0 {
											panic(fmt.Sprintf(
												"%s %s %s section[%d] %s %s's subject id should never be zero",
												curriculum.CurriculumCode, semester.Name, year_level.Name, section_idx, subject.Code, subject.Name,
											))
										}

										day_sched.GetTimeSlot(selected_time_slot).SetSubjectID(subject.ID)
										day_sched.GetTimeSlot(selected_time_slot).SetInstructorID(selected_instructor.InstructorID)
										day_sched.GetTimeSlot(selected_time_slot).SetRoomID(selected_room.RoomID)

										time_slot_assignment_sanity_counter++
									}

									if time_slot_assignment_sanity_counter != subject_total_time_slots {
										panic(
											"total time slot assigned did not match the subject total time slot",
										)
									}

									if !is_subject_type_added_once {
										selected_instructor.AssignedSubjects++
										is_subject_type_added_once = true
									}

									selected_instructor.TotalTeachingHours += float32(subject_hours)

									subject_recorder[subject.ID] = subject

									/////////////////////////////////////////////////////////////////////////////////////////////////////////
									//                                  BREAK day AND time_slot LOOP
									/////////////////////////////////////////////////////////////////////////////////////////////////////////

									day = 9999
									time_slot = 9999
								} // ------------- end of time_slot loop -------------
							} // ------------- end of day loop -------------
						} // ------------- end of class_type_iter loop -------------

						// map encoding resource that this subject is already assigned.

						if _, has_sched_idx := encoding_resource.IsSchedIdxToSubIdToSkip[non_final_sched_idx]; !has_sched_idx {
							encoding_resource.IsSchedIdxToSubIdToSkip[non_final_sched_idx] = make(map[uint16]bool)
							encoding_resource.IsSchedIdxToSubIdToSkip[non_final_sched_idx][subject.ID] = true
						} else {
							if _, has_sub_id := encoding_resource.IsSchedIdxToSubIdToSkip[non_final_sched_idx][subject.ID]; !has_sub_id {
								encoding_resource.IsSchedIdxToSubIdToSkip[non_final_sched_idx][subject.ID] = true
							} else {
								panic("woah woah woah!, you're not supposed to be here")
							}
						}
					} // ------------- end of subject loop -------------

					// front compressed distribution : end

					if len(subject_recorder) != len(semester.Subjects) {
						panic(fmt.Sprintf(
							"there are some subjects in %s %s %s %s that was not assigned for some reason s(%d/%d), i(%d), r(%d)",
							dept_id_to_department[curriculum.DepartmentID].Code,
							curriculum.CurriculumCode,
							year_level.Name,
							semester.Name,
							len(subject_recorder), len(semester.Subjects),
							len(instructors),
							len(room_type_to_rooms[Rooms.ROOM_TYPE_LAB])+len(room_type_to_rooms[Rooms.ROOM_TYPE_LEC]),
						))
					}

					individual_university_schedules[counted_sections] = week_time_table
					counted_sections++
				} // ------------- end of section_idx loop -------------
			} // ------------- end of semester_idx loop -------------
		} // ------------- end of year_level loop -------------
	} // ------------- end of curriculum loop -------------

	return individual_university_schedules, encoding_resource, nil
}

// TODO: when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying
// each sections in the generated university schedules, make the generated schedule and index global for access.

func NewEmptyIndividual(
	curriculums []Curriculum.Curriculum,
	selected_semester int,
) Schedule.UniTimeTables {

	individual_university_schedules := make(Schedule.UniTimeTables, 0, 128)

	for _, curriculum := range curriculums {
		for _, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue // skip inactive year levels
			}

			for semester_idx, semester := range year_level.Semesters {

				if selected_semester != semester_idx {
					continue // skip not selected semesters
				}

				for section_idx := 0; section_idx < semester.Sections; section_idx++ {
					week_time_table := Schedule.WeekTimeTable{}
					individual_university_schedules = append(individual_university_schedules, week_time_table)
				} // ------------- end of section_idx loop -------------
			} // ------------- end of semester_idx loop -------------
		} // ------------- end of year_level loop -------------
	} // ------------- end of curriculum loop -------------

	return individual_university_schedules
}
