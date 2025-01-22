package Schedule

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Storage"
)

// The type that represent all of the weekly schedules of each classes / sections
// in the whole university.
//
// this is just an array of `ScheduleWeek` types.
type UniTimeTables []WeekTimeTable

func NewUniTimeTables(num_of_time_tables uint) UniTimeTables {
	return make(UniTimeTables, num_of_time_tables)
}

func (uni_sched *UniTimeTables) Get(class_section_idx int) *WeekTimeTable {
	total_university_sections := len(*uni_sched)

	if class_section_idx < 0 || class_section_idx >= total_university_sections {
		panic(fmt.Sprintf(
			"GetSectionSchedule(section_idx = %d | min:max = 0:%d): error index out of bounds",
			class_section_idx, (total_university_sections - 1),
		))
	}

	return &(*uni_sched)[class_section_idx]
}

func (uni_sched *UniTimeTables) Validate() []error {
	list_of_errors := make([]error, 0, 16)

	// vertical checks
	for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
		for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {

			instructor_counter := make(map[uint16][]uint16)

			// TODO: enable the code below after room assigning is implemented.

			// room_counter := make(map[uint16]int)

			for section_idx := 0; section_idx < len(*uni_sched); section_idx++ {

				subject_id := (*uni_sched)[section_idx][day][time_slot].subjectID

				instructor_id := (*uni_sched)[section_idx][day][time_slot].instructorID

				// TODO: enable the code below after room assigning is implemented.

				// room_id := (*uni_sched)[section_idx][day][time_slot].roomID

				if subject_id == 0 && instructor_id != 0 {
					err_json := &UniInstructorValidationError{
						Msg:                 "an instructor was assigned, but no subject was scheduled for the time slot.",
						Day:                 day,
						TimeSlot:            time_slot,
						InstructorID:        instructor_id,
						OverlappingSections: nil,
					}

					json_err_str, err := json.Marshal(err_json)
					if err != nil {
						panic(err)
					}

					list_of_errors = append(list_of_errors, fmt.Errorf("%s",
						strings.Replace(string(json_err_str), `,"OverlappingSections":null`, "", 1),
					))
				}

				// TODO: enable the code below after room assigning is implemented.

				// if subject_id == 0 && room_id != 0 {
				// 	err_json := &UniRoomValidationError{
				// 		Msg:                 "a room was assigned, but no subject was scheduled for the time slot.",
				// 		Day:                 day,
				// 		TimeSlot:            time_slot,
				// 		RooomID:             room_id,
				// 		OverlappingSections: nil,
				// 	}

				// 	json_err_str, err := json.Marshal(err_json)
				// 	if err != nil {
				// 		panic(err)
				// 	}

				// 	list_of_errors = append(list_of_errors, fmt.Errorf("%s", json_err_str))
				// }

				// if there is an instructor assigned to a time slot add it to counter.
				if instructor_id > 0 {
					_, exist := instructor_counter[instructor_id]

					if !exist {
						instructor_counter[instructor_id] = make([]uint16, 0, 4)
					}

					instructor_counter[instructor_id] = append(instructor_counter[instructor_id], uint16(section_idx))
				}

				// TODO: enable the code below after room assigning is implemented.

				// if room_id > 0 {
				// 	_, exist := room_counter[room_id]

				// 	if !exist {
				// 		room_counter[room_id] = make([]uint16, 0, 4)
				// 	}

				// 	room_counter[room_id] = append(room_counter[room_id], uint16(section_idx))
				// }
			}

			for k, v := range instructor_counter {
				if len(v) > 1 {
					err_json := &UniInstructorValidationError{
						Msg:                 "overlapping instructor time slot",
						Day:                 day,
						TimeSlot:            time_slot,
						InstructorID:        k,
						OverlappingSections: v,
					}

					json_err_str, err := json.Marshal(err_json)
					if err != nil {
						panic(err)
					}

					list_of_errors = append(list_of_errors, fmt.Errorf("%s", json_err_str))
				}
			}

			// TODO: enable the code below after room assigning is implemented.

			// for k, v := range room_counter {
			// 	if len(v) > 1 {

			// 		err_json := &UniRoomValidationError{
			// 			Msg:                 "overlapping room time slot",
			// 			Day:                 day,
			// 			TimeSlot:            time_slot,
			// 			RooomID:             k,
			// 			OverlappingSections: v,
			// 		}

			// 		json_err_str, err := json.Marshal(err_json)
			// 		if err != nil {
			// 			panic(err)
			// 		}

			// 		list_of_errors = append(list_of_errors, fmt.Errorf("%s", json_err_str))
			// 	}
			// }
		}
	}

	persistence := Storage.PersistenceService{&Storage.JsonFilePersistence{}}
	persistence.Service.GetAllRooms()

	// for _, curriculum_items := range persistence.Service.GetAllCurriculum() {
	// 	for _, year_level := range curriculum_items.YearLevels {
	// 		for _, semester := range year_level.Semesters {
	// 			for _, subject := range semester.Subjects {

	// 			}
	// 		}
	// 	}
	// }

	// TODO: implement horizontal checks

	return list_of_errors
}

type UniInstructorValidationError struct {
	Msg                 string
	Day                 int
	TimeSlot            int
	InstructorID        uint16
	OverlappingSections []uint16
}

type UniRoomValidationError struct {
	Msg                 string
	Day                 int
	TimeSlot            int
	RooomID             uint16
	OverlappingSections []uint16
}
