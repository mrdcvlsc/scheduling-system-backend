package Schedule

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

// The type that represent all of the weekly schedules of each classes / sections
// in the whole university.
//
// this is just an array of `ScheduleWeek` types.
type UniTimeTables []WeekTimeTable

func NewUniTimeTables(num_of_time_tables uint) UniTimeTables {
	return make(UniTimeTables, num_of_time_tables)
}

func (university_sched UniTimeTables) GetWeekTimeTable(class_section_idx int) *WeekTimeTable {
	total_university_sections := len(university_sched)

	if class_section_idx < 0 || class_section_idx >= total_university_sections {
		panic(fmt.Sprintf(
			"GetSectionSchedule(section_idx = %d | min:max = 0:%d): error index out of bounds",
			class_section_idx, (total_university_sections - 1),
		))
	}

	return &university_sched[class_section_idx]
}

const TIME_SLOT_BYTE_SIZE int = 6 // 3 uint16 = 6 bytes.

func SerializeUniversitySchedule(university_sched UniTimeTables) []byte {
	serialized_data := make([]byte, (len(university_sched) * Const.N_WEEKLY_TIME_SLOTS * TIME_SLOT_BYTE_SIZE))

	for section_idx := 0; section_idx < len(university_sched); section_idx++ {
		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				idx_2D_to_1D := (day * Const.N_DAILY_TIME_SLOTS) + time_slot
				serialized_time_slot_idx := (section_idx*Const.N_WEEKLY_TIME_SLOTS + idx_2D_to_1D) * TIME_SLOT_BYTE_SIZE

				binary.LittleEndian.PutUint16(
					serialized_data[serialized_time_slot_idx:serialized_time_slot_idx+2],
					university_sched[section_idx][day][time_slot].subjectID,
				)

				binary.LittleEndian.PutUint16(
					serialized_data[serialized_time_slot_idx+2:serialized_time_slot_idx+4],
					university_sched[section_idx][day][time_slot].instructorID,
				)

				binary.LittleEndian.PutUint16(
					serialized_data[serialized_time_slot_idx+4:serialized_time_slot_idx+6],
					university_sched[section_idx][day][time_slot].roomID,
				)

			}
		}
	}

	return serialized_data
}

func DeserializeUniversitySchedule(serialized_data []byte) UniTimeTables {
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

	return uni_sched
}

type roomCountAndCapacity struct {
	OverlappingSections []uint16
	Capacity            uint16
}

// returns true if the length is 0 or if there are no subjects allocated in the time tables. false otherwise.
//
// can be used on incomplete university schedules (uni sched without other departments schedule).
func (university_sched UniTimeTables) IsEmpty() bool {

	if len(university_sched) == 0 {
		return true
	}

	subject_count := 0

	for _, section_sched := range university_sched {
		for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
			for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
				if section_sched[day][time_slot].GetSubjectID() != 0 {
					subject_count++
				}
			}
		}
	}

	return subject_count == 0
}

// validate rooms and instructors time slot availability, this function detects overlapping instructor or room time slots.
func (university_sched UniTimeTables) VerticalRangedValidation(
	rooms []Rooms.Room,
	day_start, day_size,
	time_slot_start, time_slot_size int,
) []error {
	errs_slice := make([]error, 0, 16)

	room_id_to_capacity := make(map[uint16]uint16)

	for _, room := range rooms {
		room_id_to_capacity[room.RoomID] = room.Capacity
	}

	/////////////////////////////////////////////////////////////////////////////////
	//                             VERTICAL CHECKS
	/////////////////////////////////////////////////////////////////////////////////

	for day := day_start; day < (day_start + day_size); day++ {
		for time_slot := time_slot_start; time_slot < (time_slot_start + time_slot_size); time_slot++ {

			instructor_counter := make(map[uint16][]uint16)

			room_counter := make(map[uint16]*roomCountAndCapacity)

			for section_idx := 0; section_idx < len(university_sched); section_idx++ {

				subject_id := university_sched[section_idx][day][time_slot].GetSubjectID()

				instructor_id := university_sched[section_idx][day][time_slot].GetInstructorID()

				room_id := university_sched[section_idx][day][time_slot].GetRoomID()

				if subject_id == 0 && instructor_id != 0 {
					err_json := &UniInstructorValidationError{
						Msg:                 "an instructor was assigned, but no subject was scheduled for the time slot.",
						Day:                 day,
						TimeSlot:            time_slot,
						InstructorID:        instructor_id,
						OverlappingSections: nil,
					}

					msg, err := json.Marshal(err_json)
					if err != nil {
						panic(err)
					}

					errs_slice = append(errs_slice, fmt.Errorf("%s",
						strings.Replace(string(msg), `,"OverlappingSections":null`, "", 1),
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

					msg, err := json.Marshal(err_json)
					if err != nil {
						panic(err)
					}

					errs_slice = append(errs_slice, fmt.Errorf("%s",
						strings.Replace(string(msg), `,"OverlappingSections":null`, "", 1),
					))
				}

				// if there is an instructor assigned to a time slot add it to counter.

				if instructor_id > 0 {
					_, has_instructor_id := instructor_counter[instructor_id]

					if !has_instructor_id {
						instructor_counter[instructor_id] = make([]uint16, 0, 4)
					}

					instructor_counter[instructor_id] = append(instructor_counter[instructor_id], uint16(section_idx))
				}

				if room_id > 0 {
					_, has_room_id := room_counter[room_id]

					if !has_room_id {
						room_counter[room_id] = &roomCountAndCapacity{}
						room_counter[room_id].Capacity = room_id_to_capacity[room_id]
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

					msg, err := json.Marshal(err_json)
					if err != nil {
						panic(err)
					}

					errs_slice = append(errs_slice, fmt.Errorf("%s", msg))
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

					msg, err := json.Marshal(err_json)
					if err != nil {
						panic(err)
					}

					errs_slice = append(errs_slice, fmt.Errorf("%s", msg))
				}
			}
		}
	}

	return errs_slice
}

// validate rooms and instructors time slot availability, this function detects overlapping instructor or room time slots.
func (university_sched UniTimeTables) VerticalValidation(rooms []Rooms.Room) []error {
	return university_sched.VerticalRangedValidation(
		rooms,
		0, Const.N_WEEKLY_SCHOOL_DAYS,
		0, Const.N_DAILY_TIME_SLOTS,
	)
}

/*
validate assigned subjects to every section schedules in the whole university.

to validate whole university schedules, set department to encode to nil:

	department_to_validate = nil
*/
func (university_sched UniTimeTables) HorizontalValidation(
	curriculums []Curriculum.Curriculum,
	department_to_validate map[uint16]bool, selected_semester int,
) []error {

	errs_slice := make([]error, 0, 16)

	/////////////////////////////////////////////////////////////////////////////////
	//                            HORIZONTAL CHECKS
	/////////////////////////////////////////////////////////////////////////////////

	total_university_sections := Curriculum.GetTotalNumberOfSections(curriculums, selected_semester)

	if total_university_sections != len(university_sched) {
		errs_slice = append(errs_slice, fmt.Errorf(
			"read total university sections (%d) in persistence did not match the university schedule instance (%d)",
			total_university_sections, len(university_sched),
		))
	}

	schedule_idx := 0
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

					if department_to_validate != nil {
						is_to_validate := department_to_validate[curriculum.DepartmentID]

						if !is_to_validate {
							schedule_idx++
							continue // skip section that does not need horizontal validation
						}
					}

					subject_id_to_time_slot_count := make(map[uint16]int)

					for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
						for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
							subject_id := university_sched[schedule_idx][day][time_slot].GetSubjectID()

							if subject_id == 0 {
								continue
							}

							_, has_subjsubject_id := subject_id_to_time_slot_count[subject_id]

							if !has_subjsubject_id {
								subject_id_to_time_slot_count[subject_id] = 1
							} else {
								subject_id_to_time_slot_count[subject_id]++
							}
						}
					}

					if len(semester.Subjects) != len(subject_id_to_time_slot_count) {
						errs_slice = append(errs_slice, fmt.Errorf(
							"detected %d missing subject(s) in %s, %s, %s, section %s (usi:%d)",
							len(semester.Subjects)-len(subject_id_to_time_slot_count),
							curriculum.CurriculumCode,
							semester.Name,
							year_level.Name,
							Curriculum.SECTION[section_idx],
							schedule_idx,
						))
					}

					for _, subject := range semester.Subjects {
						_, has_subject_id := subject_id_to_time_slot_count[subject.ID]

						if !has_subject_id {
							errs_slice = append(errs_slice, fmt.Errorf(
								"the subject %s was not assigned to %s, %s, %s, section %s (usi:%d)",
								subject.Code, curriculum.CurriculumCode, year_level.Name, semester.Name, Curriculum.SECTION[section_idx], schedule_idx,
							))
						} else if ((subject.LecHours + subject.LabHours) * Const.N_HOUR_TIME_SLOTS) > uint8(subject_id_to_time_slot_count[subject.ID]) {
							errs_slice = append(errs_slice, fmt.Errorf(
								"the subject %s has missing time slot allocations, expecting %d, but only found %d in %s, %s, %s, section %s (usi:%d)",
								subject.Code,
								((subject.LecHours+subject.LabHours)*Const.N_HOUR_TIME_SLOTS), uint8(subject_id_to_time_slot_count[subject.ID]),
								curriculum.CurriculumCode, year_level.Name, semester.Name, Curriculum.SECTION[section_idx], schedule_idx,
							))
						} else if ((subject.LecHours + subject.LabHours) * Const.N_HOUR_TIME_SLOTS) < uint8(subject_id_to_time_slot_count[subject.ID]) {
							errs_slice = append(errs_slice, fmt.Errorf(
								"the subject %s has extra time slot allocations, expecting only %d, but found %d in %s, %s, %s, section %s (usi:%d)",
								subject.Code,
								((subject.LecHours+subject.LabHours)*Const.N_HOUR_TIME_SLOTS), uint8(subject_id_to_time_slot_count[subject.ID]),
								curriculum.CurriculumCode, year_level.Name, semester.Name, Curriculum.SECTION[section_idx], schedule_idx,
							))
						}
					}

					schedule_idx++
				} // ------------- end of section_idx loop -------------
			} // ------------- end of semester_idx loop -------------
		} // ------------- end of year_level loop -------------
	} // ------------- end of curriculum loop -------------

	return errs_slice
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
