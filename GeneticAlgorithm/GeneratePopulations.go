package geneticalgorithm

import (
	"math/rand"
	"time"

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

func NewIndividual(selected_semester, distribution_type int) (Schedule.UniTimeTables, error) {

	// test.
	persistence := Storage.PersistenceService{Service: &Storage.JsonFilePersistence{}}
	rng := rand.New(rand.NewSource(time.Now().UnixMilli()))

	// process rooms

	// department_rooms = Map of rooms => [Key:DepartmentId][Key:RoomType]
	department_rooms := make(map[uint16]map[uint16][]Rooms.Room)

	{
		all_rooms := persistence.Service.GetAllRooms()

		for _, room := range all_rooms {
			_, dept_key_exist := department_rooms[room.DepartmentID]

			if !dept_key_exist {
				department_rooms[room.DepartmentID] = make(map[uint16][]Rooms.Room)
			}

			_, roomtype_key_exist := department_rooms[room.DepartmentID][room.RoomType]

			if !roomtype_key_exist {
				department_rooms[room.DepartmentID][room.RoomType] = make([]Rooms.Room, 1, 8)
				department_rooms[room.DepartmentID][room.RoomType][0] = room
			} else {
				department_rooms[room.DepartmentID][room.RoomType] = append(
					department_rooms[room.DepartmentID][room.RoomType], room,
				)
			}
		}
	}

	// process instructors

	// department_instructors = Map of instructors => [Key:DepartmentId]
	department_instructors := make(map[uint16][]Instructors.Instructor)

	{
		all_instructors := persistence.Service.GetAllInstructors()

		for _, instructor := range all_instructors {
			_, exist := department_instructors[instructor.DepartmentID]

			if !exist {
				department_instructors[instructor.DepartmentID] = make([]Instructors.Instructor, 1, 8)
				department_instructors[instructor.DepartmentID][0] = instructor
			} else {
				department_instructors[instructor.DepartmentID] = append(
					department_instructors[instructor.DepartmentID], instructor,
				)
			}
		}
	}

	// generate time tables

	// fmt.Println("department_rooms Map:")
	// Utils.PrettyPrint(department_rooms)
	// fmt.Print("\n\n")

	total_number_of_sections := 0
	individual := make(Schedule.UniTimeTables, 0, 64)

	all_curriculums := persistence.Service.GetAllCurriculum()

	for _, curriculum := range all_curriculums {
		for _, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue
			}

			for semester_idx, semester := range year_level.Semesters {

				if selected_semester != semester_idx {
					continue
				}

				total_number_of_sections += semester.Sections

				for section := 0; section < semester.Sections; section++ {

					// shuffle the subjects
					rng.Shuffle(len(semester.Subjects), func(i, j int) {
						semester.Subjects[i], semester.Subjects[j] = semester.Subjects[j], semester.Subjects[i]
					})

					// TODO: add a property in the instructor struct that will record the
					// total number of subjects that instructor will teach.

					// TODO: replace instructor slice shuffle with sort instead

					// shuffle instructors / maybe change to sort based on current numbers of subject to teach, prioritize the lowest
					rng.Shuffle(len(department_instructors[curriculum.DepartmentID]), func(i, j int) {
						department_instructors[curriculum.DepartmentID][i], department_instructors[curriculum.DepartmentID][j] = department_instructors[curriculum.DepartmentID][j], department_instructors[curriculum.DepartmentID][i]
					})

					// shuffle rooms
					for _, rooms_group_by_type := range department_rooms[curriculum.DepartmentID] {
						rng.Shuffle(len(rooms_group_by_type), func(i, j int) {
							rooms_group_by_type[i], rooms_group_by_type[j] = rooms_group_by_type[j], rooms_group_by_type[i]
						})
					}

					week_time_table := Schedule.WeekTimeTable{}

					// front compressed distribution : start

					// for now the this is a brute force implementation, yet it is still enough
					// considering the current low numbers of courses, instructors, rooms and sections.
					// but I believe this still can be optimize in the future if we wanted to.

					// or just create another implementation for inidividual generation.

					for _, subject := range semester.Subjects {

						var selected_instructor Instructors.Instructor

						for class_type := 0; class_type < 2; class_type++ {

							var class_hour int

							if class_type == 0 {
								class_hour = int(subject.LecHours)
							} else {
								class_hour = int(subject.LabHours)
							}

							if class_hour == 0 {
								continue
							}

							// class_type == 0 == lecture
							// class_type == 1 == laboratory

							// search an available time slot for the current subject

							for day := 0; day < Schedule.N_WEEKLY_SCHOOL_DAYS; day++ {

								day_sched := week_time_table.Get(day)

								for time_slot := 0; time_slot < Schedule.N_DAILY_TIME_SLOTS; time_slot++ {

									// assign a subject to a time slot

									if !day_sched.Availability(time_slot, class_hour) {
										continue
									}

									// TODO: check if there is already an assigned instructor for the subject.
									//
									// IF NONE: iterate over all of the sorted instructors to find which
									// one is available, if there is no instructor available for the current
									// time slot, continue to the next iteration of the time slot loop.
									//
									// ELSE: if there is already a selected instructor but that instructor
									// is not available for the current time slot, continue to the next
									// iteration of the time slot loop.

									// TODO: if this is the last time slot and there is still no
									// instructor available, throw an error saying there is not
									// enough instructors.

									// TODO: iterate over all of sorted rooms to find which one is
									// available, if there is no room available for the current time
									// slot, continue to the next iteration of the loop.

									// TODO: if this is the last time slot and there is still no
									// room available, throw an error saying there is not enough room.

									// TODO: after all of the checks and searches, if there are instructors
									// and rooms that are available for a specific subject and time slot,
									// then assign the selected subject, instructor and room to the time slot.
								}
							}

							if selected_instructor.InstructorID == 0 {
								// select an instructor if there is no one selected yet

							}
						}
					}

					// front compressed distribution : end

					individual = append(individual, week_time_table)
				}
			}
		}
	}

	// fmt.Println("Total Number of Sections : ", total_number_of_sections)

	// TODO: during the first schedule generation throw error if there is not enough instructors

	// TODO: during the first schedule generation throw error if there is not enough rooms

	return individual, nil
}
