package Schedule

import (
	"encoding/binary"
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

const TIME_SLOT_BYTE_SIZE int = 6 // 3 uint16 = 6 bytes.

func SerializeUniversitySchedule(uni_sched *UniTimeTables) []byte {
	serialized_data := make([]byte, (len(*uni_sched) * Const.N_WEEKLY_TIME_SLOTS * TIME_SLOT_BYTE_SIZE))

	for section_idx := 0; section_idx < len(*uni_sched); section_idx++ {
		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				idx_2D_to_1D := (day * Const.N_DAILY_TIME_SLOTS) + time_slot
				serialized_time_slot_idx := (section_idx*Const.N_WEEKLY_TIME_SLOTS + idx_2D_to_1D) * TIME_SLOT_BYTE_SIZE

				binary.LittleEndian.PutUint16(serialized_data[serialized_time_slot_idx:serialized_time_slot_idx+2], (*uni_sched)[section_idx][day][time_slot].subjectID)
				binary.LittleEndian.PutUint16(serialized_data[serialized_time_slot_idx+2:serialized_time_slot_idx+4], (*uni_sched)[section_idx][day][time_slot].instructorID)
				binary.LittleEndian.PutUint16(serialized_data[serialized_time_slot_idx+4:serialized_time_slot_idx+6], (*uni_sched)[section_idx][day][time_slot].roomID)
			}
		}
	}

	return serialized_data
}

func DeserializeUniversitySchedule(serialized_data []byte) *UniTimeTables {
	uni_sched := make(UniTimeTables, (len(serialized_data) / (Const.N_WEEKLY_TIME_SLOTS * TIME_SLOT_BYTE_SIZE)))

	for section_idx := 0; section_idx < len(uni_sched); section_idx++ {
		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				idx_2D_to_1D := (day * Const.N_DAILY_TIME_SLOTS) + time_slot
				serialized_time_slot_idx := (section_idx*Const.N_WEEKLY_TIME_SLOTS + idx_2D_to_1D) * TIME_SLOT_BYTE_SIZE

				serialized_subject_id := serialized_data[serialized_time_slot_idx : serialized_time_slot_idx+2]
				serialized_instructor_id := serialized_data[serialized_time_slot_idx+2 : serialized_time_slot_idx+4]
				serialized_room_id := serialized_data[serialized_time_slot_idx+4 : serialized_time_slot_idx+6]

				uni_sched[section_idx][day][time_slot].subjectID = binary.LittleEndian.Uint16(serialized_subject_id)
				uni_sched[section_idx][day][time_slot].instructorID = binary.LittleEndian.Uint16(serialized_instructor_id)
				uni_sched[section_idx][day][time_slot].roomID = binary.LittleEndian.Uint16(serialized_room_id)
			}
		}
	}

	return &uni_sched
}

type room_count_and_capacity struct {
	OverlappingSections []uint16
	Capacity            uint16
}

func (uni_sched *UniTimeTables) IsEmpty() bool {

	if len(*uni_sched) == 0 {
		return true
	}

	empty_subject_count := 0

	for _, section_sched := range *uni_sched {
		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				if section_sched[day][time_slot].subjectID == 0 {
					empty_subject_count++
				}
			}
		}
	}

	return empty_subject_count == (Const.N_WEEKLY_TIME_SLOTS * len(*uni_sched))
}

func (uni_sched *UniTimeTables) Validate() []error {

	list_of_errors := make([]error, 0, 16)

	persistence := Storage.PersistenceService{ReaderService: &Storage.JsonFilePersistence{}}
	all_rooms, err_all_rooms := persistence.ReaderService.GetAllRooms()

	if err_all_rooms != nil {
		list_of_errors = append(list_of_errors, err_all_rooms)
		return list_of_errors
	}

	map_room_id_and_capacity := make(map[uint16]uint16)

	for _, room := range all_rooms {
		map_room_id_and_capacity[room.RoomID] = room.Capacity
	}

	/////////////////////////////////////////////////////////////////////////////////
	//                             VERTICAL CHECKS
	/////////////////////////////////////////////////////////////////////////////////

	for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
		for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {

			instructor_counter := make(map[uint16][]uint16)

			room_counter := make(map[uint16]*room_count_and_capacity)

			for section_idx := 0; section_idx < len(*uni_sched); section_idx++ {

				subject_id := (*uni_sched)[section_idx][day][time_slot].subjectID

				instructor_id := (*uni_sched)[section_idx][day][time_slot].instructorID

				room_id := (*uni_sched)[section_idx][day][time_slot].roomID

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

				if subject_id == 0 && room_id != 0 {
					err_json := &UniRoomValidationError{
						Msg:                 "a room was assigned, but no subject was scheduled for the time slot.",
						Day:                 day,
						TimeSlot:            time_slot,
						RooomID:             room_id,
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

				// if there is an instructor assigned to a time slot add it to counter.

				if instructor_id > 0 {
					_, exist := instructor_counter[instructor_id]

					if !exist {
						instructor_counter[instructor_id] = make([]uint16, 0, 4)
					}

					instructor_counter[instructor_id] = append(instructor_counter[instructor_id], uint16(section_idx))
				}

				if room_id > 0 {
					_, exist := room_counter[room_id]

					if !exist {
						room_counter[room_id] = &room_count_and_capacity{}
						room_counter[room_id].Capacity = map_room_id_and_capacity[room_id]
					}

					room_counter[room_id].OverlappingSections = append(room_counter[room_id].OverlappingSections, uint16(section_idx))
				}
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

			for k, v := range room_counter {
				if len(v.OverlappingSections) > int(v.Capacity) {

					err_json := &UniRoomValidationError{
						Msg:                 "overlapping room time slot",
						Day:                 day,
						TimeSlot:            time_slot,
						RooomID:             k,
						OverlappingSections: v.OverlappingSections,
					}

					json_err_str, err := json.Marshal(err_json)
					if err != nil {
						panic(err)
					}

					list_of_errors = append(list_of_errors, fmt.Errorf("%s", json_err_str))
				}
			}
		}
	}

	/////////////////////////////////////////////////////////////////////////////////
	//                            HORIZONTAL CHECKS
	/////////////////////////////////////////////////////////////////////////////////

	// TODO: implement horizontal checks - incomplete code below

	// persistence := Storage.PersistenceService{Service: &Storage.JsonFilePersistence{}}

	// subjects := persistence.Service.GetAllSubjects()
	// map_id_subjects := make(map[uint16]Curriculum.Subject)

	// for _, subject := range subjects {
	// 	map_id_subjects[subject.ID] = subject
	// }

	// rooms := persistence.Service.GetAllRooms()
	// map_id_rooms := make(map[uint16]Rooms.Room)

	// for _, room := range rooms {
	// 	map_id_rooms[room.RoomID] = room
	// }

	// for section_idx := 0; section_idx < len(*uni_sched); section_idx++ {
	// 	for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
	// 		for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
	// 			subject_id := (*uni_sched)[section_idx][day][time_slot].subjectID

	// 			instructor_id := (*uni_sched)[section_idx][day][time_slot].instructorID

	// 			room_id := (*uni_sched)[section_idx][day][time_slot].roomID
	// 		}
	// 	}
	// }

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
