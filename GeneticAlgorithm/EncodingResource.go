package GeneticAlgorithm

import (
	"fmt"
	"log"
	"sort"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
)

// struct type containing members that is use to track resource allocation/utilization in a schedule.
type EncodingResource struct {
	IsSchedIdxToSubIdToSkip map[uint16]map[uint16]bool
	DeptIdToInstructors     map[uint16][]Instructors.Instructor
	DeptIdToRoomtypeToRooms map[uint16]map[uint16][]Rooms.Room
}

// TODO: make a function to compare two `EncodingResource`.
func IsEqualEncodingResource(a, b *EncodingResource) bool {

	////////////////////////////////////////////////////////////////////////////////////////

	if len(a.IsSchedIdxToSubIdToSkip) != len(b.IsSchedIdxToSubIdToSkip) {
		return false
	}

	for a_out_k, a_out_v := range a.IsSchedIdxToSubIdToSkip {

		b_out_v, has_b_out_k := b.IsSchedIdxToSubIdToSkip[a_out_k]

		if !has_b_out_k {
			return false
		}

		if len(a_out_v) != len(b_out_v) {
			return false
		}

		for a_in_k, a_in_v := range a_out_v {

			b_in_v, has_b_in_k := b_out_v[a_in_k]

			if !has_b_in_k {
				return false
			}

			if b_in_v != a_in_v {
				return false
			}
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	for a_out_k, a_out_v := range a.DeptIdToRoomtypeToRooms {
		b_out_v, has_b_out_k := b.DeptIdToRoomtypeToRooms[a_out_k]

		if !has_b_out_k {
			return false
		}

		if len(a_out_v) != len(b_out_v) {
			return false
		}

		for a_in_k, a_in_v := range a_out_v {

			b_in_v, has_b_in_k := b_out_v[a_in_k]

			if !has_b_in_k {
				return false
			}

			if len(a_in_v) != len(b_in_v) {
				return false
			}

			sort.Slice(a_in_v, func(i, j int) bool {
				return a_in_v[i].RoomID < a_in_v[j].RoomID
			})

			sort.Slice(b_in_v, func(i, j int) bool {
				return b_in_v[i].RoomID < b_in_v[j].RoomID
			})

			for a_room_idx, a_room := range a_in_v {
				if a_room != b_in_v[a_room_idx] {
					for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
						for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
							if a_room.GetTimeSlotClassCount(day, time_slot) != b_in_v[a_room_idx].GetTimeSlotClassCount(day, time_slot) {
								fmt.Printf("RIDX: %d, room_id: %d | d(%d), t(%d) => a(%d), b(%d)\n", a_room_idx, a_room.RoomID, day, time_slot, a_room.GetTimeSlotClassCount(day, time_slot), b_in_v[a_room_idx].GetTimeSlotClassCount(day, time_slot))
								return false
							}
						}
					}
					return false
				}
			}
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	if len(a.DeptIdToInstructors) != len(b.DeptIdToInstructors) {
		return false
	}

	for a_k, a_v := range a.DeptIdToInstructors {

		b_v, has_b_k := b.DeptIdToInstructors[a_k]

		if !has_b_k {
			return false
		}

		if len(a_v) != len(b_v) {
			return false
		}

		sort.Slice(a_v, func(i, j int) bool {
			return a_v[i].InstructorID < a_v[j].InstructorID
		})

		sort.Slice(b_v, func(i, j int) bool {
			return b_v[i].InstructorID < b_v[j].InstructorID
		})

		for instructor_idx, instructor := range a_v {
			if instructor != b_v[instructor_idx] {
				return false
			}
		}
	}

	return true
}

/*
use to generate an `EncodingResource` for a `UniTimeTables`, can be used for user
configured university schedules or university schedules loaded from the persistence
since they don't have any associated `EncodingResource` instance.

NOTE TO SELF IN THE FUTURE:

if you want this function to support department specific encoding resource generation... DON'T DO IT!
because generating encoding resource can't be department specific, to give an example,
general instructors can be assigned to multiple class / section on different departments (like FITT teachers),
if you only generate encoding resource that is department specific, the other allocations to different
departments will not be reflected, this also applies to general rooms.
*/
func GenerateEncodingResourceFromUniTimeTable(
	university_schedules Schedule.UniTimeTables,
	curriculums []Curriculum.Curriculum,
	selected_semester int,
	resource_persistence *StorageResources.Persistence,
) (*EncodingResource, error) {

	encode_resource, err_read_default_encoding_resource := ReadDefaultEncodingResource(resource_persistence)

	if err_read_default_encoding_resource != nil {
		return nil, err_read_default_encoding_resource
	}

	//////////////////////////////////////////////////////////////////////////////////////
	//                              FLATTEN ENCODING RESOURCES
	//////////////////////////////////////////////////////////////////////////////////////

	////////////////////////////////////////////////////////////////////////////////////////

	room_id_to_room := make(map[uint16]*Rooms.Room)

	for out_key, out_v := range encode_resource.DeptIdToRoomtypeToRooms {
		for in_key, in_v := range out_v {
			for room_idx, room := range in_v {
				room_id_to_room[room.RoomID] = &encode_resource.DeptIdToRoomtypeToRooms[out_key][in_key][room_idx]
			}
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	instructor_id_to_instructor := make(map[uint16]*Instructors.Instructor)

	for k, v := range encode_resource.DeptIdToInstructors {
		for instructor_idx, instructor := range v {
			instructor_id_to_instructor[instructor.InstructorID] = &encode_resource.DeptIdToInstructors[k][instructor_idx]
		}
	}

	//////////////////////////////////////////////////////////////////////////////////////
	//                           RE-CREATE ENCODING RESOURCE DATA
	//////////////////////////////////////////////////////////////////////////////////////

	counted_sections := 0

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

					for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
						for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
							subject_id := university_schedules[counted_sections][day].GetTimeSlot(time_slot).GetSubjectID()

							if subject_id != 0 {
								instructor_id := university_schedules[counted_sections][day].GetTimeSlot(time_slot).GetInstructorID()
								room_id := university_schedules[counted_sections][day].GetTimeSlot(time_slot).GetRoomID()

								selected_instructor := instructor_id_to_instructor[instructor_id]
								selected_room := room_id_to_room[room_id]

								if instructor_id == 0 {
									log.Panic("there should be an instructor allocation here, why there is none?")
								}

								if room_id == 0 {
									log.Panic("there should be a room allocation here, why there is none?")
								}

								selected_instructor.Time.SetAvailability(false, day, time_slot)
								selected_room.IncTimeSlotClassCount(day, time_slot)

								_, has_sched_idx := encode_resource.IsSchedIdxToSubIdToSkip[uint16(counted_sections)]

								if !has_sched_idx {
									encode_resource.IsSchedIdxToSubIdToSkip[uint16(counted_sections)] = make(map[uint16]bool)
								}

								_, has_subject_id := encode_resource.IsSchedIdxToSubIdToSkip[uint16(counted_sections)][subject_id]

								if !has_subject_id {
									encode_resource.IsSchedIdxToSubIdToSkip[uint16(counted_sections)][subject_id] = true
									selected_instructor.AssignedSubjects++
								}

								selected_instructor.TotalTeachingHours += (1.0 / Const.N_HOUR_TIME_SLOTS)
							}
						} // ------------- end of time_slot loop -------------
					} // ------------- end of day loop -------------

					counted_sections++
				} // ------------- end of section_idx loop -------------
			} // ------------- end of semester_idx loop -------------
		} // ------------- end of year_level loop -------------
	} // ------------- end of curriculum loop -------------

	return encode_resource, nil
}

// reads the default values of `EncodingResource` saved in a persistence instance.
func ReadDefaultEncodingResource(resource_persistence *StorageResources.Persistence) (*EncodingResource, error) {
	dept_id_to_room_type_to_rooms, err_dept_id_to_room_type_to_rooms := GenerateMapDeptIdToRoomTypeToRooms(resource_persistence)

	if err_dept_id_to_room_type_to_rooms != nil {
		return nil, err_dept_id_to_room_type_to_rooms
	}

	dept_id_to_instructors, err_dept_id_to_instructors := GenerateMapDeptIdToInstructors(resource_persistence)

	if err_dept_id_to_instructors != nil {
		return nil, err_dept_id_to_instructors
	}

	return &EncodingResource{
		IsSchedIdxToSubIdToSkip: make(map[uint16]map[uint16]bool),
		DeptIdToInstructors:     dept_id_to_instructors,
		DeptIdToRoomtypeToRooms: dept_id_to_room_type_to_rooms,
	}, nil
}
