package GeneticAlgorithm

import (
	"fmt"
	"log"
	"reflect"
	"slices"
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

	IdToInstructor map[uint16]*Instructors.Instructor // flatten `DeptIdToInstructors`
	IdToRoom       map[uint16]*Rooms.Room             // flatten `DeptIdToRoomtypeToRooms`
}

func (s *EncodingResource) MakeCopy() (*EncodingResource, error) {
	////////////////////////////////////////////////////////////////////////////////////////

	is_sched_idx_to_sub_id_to_skip := make(map[uint16]map[uint16]bool)

	for out_k, out_v := range s.IsSchedIdxToSubIdToSkip {
		is_sched_idx_to_sub_id_to_skip[out_k] = make(map[uint16]bool)

		for in_k, in_v := range out_v {
			is_sched_idx_to_sub_id_to_skip[out_k][in_k] = in_v
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_room_type_to_rooms := make(map[uint16]map[uint16][]Rooms.Room)

	for out_k, out_v := range s.DeptIdToRoomtypeToRooms {
		dept_id_to_room_type_to_rooms[out_k] = make(map[uint16][]Rooms.Room)

		for in_k, in_v := range out_v {
			dept_id_to_room_type_to_rooms[out_k][in_k] = make([]Rooms.Room, len(in_v))
			copies := copy(dept_id_to_room_type_to_rooms[out_k][in_k], in_v)

			if copies != len(in_v) {
				log.Printf("copies : %d\tlen(dept_id_to_room_type_to_rooms[out_k][in_k] = %d/%d = in_v)\n", copies, len(dept_id_to_room_type_to_rooms[out_k][in_k]), len(in_v))
				return nil, fmt.Errorf("slice elements copied %d, internal department id to rooms map copy operation failed in generate new individual function", copies)
			}
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	dept_id_to_instructors := make(map[uint16][]Instructors.Instructor)

	for k, v := range s.DeptIdToInstructors {
		dept_id_to_instructors[k] = make([]Instructors.Instructor, len(v))
		copies := copy(dept_id_to_instructors[k], v)
		if copies != len(v) {
			return nil, fmt.Errorf("slice elements copied %d, internal department id to instructors map copy operation failed in generate new individual function", copies)
		}
	}

	//////////////////////////////////////////////////////////////////////////////////////
	//                              FLATTEN ENCODING RESOURCES
	//////////////////////////////////////////////////////////////////////////////////////

	room_id_to_room := make(map[uint16]*Rooms.Room)

	for out_key, out_v := range dept_id_to_room_type_to_rooms {
		for in_key, in_v := range out_v {
			for room_idx, room := range in_v {
				room_id_to_room[room.RoomID] = &dept_id_to_room_type_to_rooms[out_key][in_key][room_idx]
			}
		}
	}

	instructor_id_to_instructor := make(map[uint16]*Instructors.Instructor)

	for k, v := range dept_id_to_instructors {
		for instructor_idx, instructor := range v {
			instructor_id_to_instructor[instructor.InstructorID] = &dept_id_to_instructors[k][instructor_idx]
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	return &EncodingResource{
		IsSchedIdxToSubIdToSkip: is_sched_idx_to_sub_id_to_skip,
		DeptIdToInstructors:     dept_id_to_instructors,
		DeptIdToRoomtypeToRooms: dept_id_to_room_type_to_rooms,

		IdToInstructor: instructor_id_to_instructor,
		IdToRoom:       room_id_to_room,
	}, nil
}

func IsEqualEncodingResource(a, b *EncodingResource) bool {

	////////////////////////////////////////////////////////////////////////////////////////

	remove_false_subject_to_skip := func(m map[uint16]map[uint16]bool) map[uint16][]uint16 {
		out := make(map[uint16][]uint16, len(m))

		for usi, subject_id_to_skip := range m {

			for subject_id, to_skip := range subject_id_to_skip {
				if !to_skip {
					continue
				}

				out[usi] = append(out[usi], subject_id)
			}

			slices.Sort(out[usi])
		}

		return out
	}

	skip_subjects_a := remove_false_subject_to_skip(a.IsSchedIdxToSubIdToSkip)
	skip_subjects_b := remove_false_subject_to_skip(b.IsSchedIdxToSubIdToSkip)

	if !reflect.DeepEqual(skip_subjects_a, skip_subjects_b) {
		log.Print("IsEqualEncodingResource: not equal IsSchedIdxToSubIdToSkip")
		return false
	}

	////////////////////////////////////////////////////////////////////////////////////////

	for a_out_k, a_out_v := range a.DeptIdToRoomtypeToRooms {
		b_out_v, has_b_out_k := b.DeptIdToRoomtypeToRooms[a_out_k]

		if !has_b_out_k {
			log.Print("IsEqualEncodingResource: not equal DeptIdToRoomtypeToRooms - case 1")
			return false
		}

		if len(a_out_v) != len(b_out_v) {
			log.Print("IsEqualEncodingResource: not equal DeptIdToRoomtypeToRooms - case 2")
			return false
		}

		for a_in_k, a_in_v := range a_out_v {

			b_in_v, has_b_in_k := b_out_v[a_in_k]

			if !has_b_in_k {
				log.Print("IsEqualEncodingResource: not equal DeptIdToRoomtypeToRooms - case 3")
				return false
			}

			if len(a_in_v) != len(b_in_v) {
				log.Print("IsEqualEncodingResource: not equal DeptIdToRoomtypeToRooms - case 4")
				return false
			}

			sort.Slice(a_in_v, func(i, j int) bool {
				return a_in_v[i].RoomID < a_in_v[j].RoomID
			})

			sort.Slice(b_in_v, func(i, j int) bool {
				return b_in_v[i].RoomID < b_in_v[j].RoomID
			})

			for a_room_idx, a_room := range a_in_v {
				if !reflect.DeepEqual(a_room, b_in_v[a_room_idx]) {
					for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
						for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
							if a_room.GetTimeSlotClassCount(day, time_slot) != b_in_v[a_room_idx].GetTimeSlotClassCount(day, time_slot) {
								fmt.Printf("RIDX: %d, room_id: %d | d(%d), t(%d) => a(%d), b(%d)\n", a_room_idx, a_room.RoomID, day, time_slot, a_room.GetTimeSlotClassCount(day, time_slot), b_in_v[a_room_idx].GetTimeSlotClassCount(day, time_slot))
								return false
							}
						}
					}

					log.Print("IsEqualEncodingResource: not equal DeptIdToRoomtypeToRooms - case 5")
					return false
				}
			}
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	if len(a.DeptIdToInstructors) != len(b.DeptIdToInstructors) {
		log.Print("IsEqualEncodingResource: not equal DeptIdToInstructors - case 1")
		return false
	}

	for a_k, a_v := range a.DeptIdToInstructors {

		b_v, has_b_k := b.DeptIdToInstructors[a_k]

		if !has_b_k {
			log.Print("IsEqualEncodingResource: not equal DeptIdToInstructors - case 2")
			return false
		}

		if len(a_v) != len(b_v) {
			log.Print("IsEqualEncodingResource: not equal DeptIdToInstructors - case 3")
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
				log.Print("IsEqualEncodingResource: not equal DeptIdToInstructors - case 4")
				return false
			}
		}
	}

	return true
}

/*
@param - `default_empty_encoding_resource` - should strictly

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
	default_empty_encoding_resource *EncodingResource,
) (*EncodingResource, error) {

	encode_resource, err_make_copy := default_empty_encoding_resource.MakeCopy()

	if err_make_copy != nil {
		return nil, err_make_copy
	}

	//////////////////////////////////////////////////////////////////////////////////////
	//                              FLATTEN ENCODING RESOURCES
	//////////////////////////////////////////////////////////////////////////////////////

	room_id_to_room := encode_resource.IdToRoom
	instructor_id_to_instructor := encode_resource.IdToInstructor

	//////////////////////////////////////////////////////////////////////////////////////
	//                           RE-CREATE ENCODING RESOURCE DATA
	//////////////////////////////////////////////////////////////////////////////////////

	IterateSectionsWeekSchedule(university_schedules, curriculums, selected_semester, nil, nil,
		func(indicies IterIndices, values IterValues) IterReturnType {
			for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
				for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
					subject_id := university_schedules[indicies.Usi][day].GetTimeSlot(time_slot).GetSubjectID()

					if subject_id != 0 {
						instructor_id := university_schedules[indicies.Usi][day].GetTimeSlot(time_slot).GetInstructorID()
						room_id := university_schedules[indicies.Usi][day].GetTimeSlot(time_slot).GetRoomID()

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

						_, has_sched_idx := encode_resource.IsSchedIdxToSubIdToSkip[uint16(indicies.Usi)]

						if !has_sched_idx {
							encode_resource.IsSchedIdxToSubIdToSkip[uint16(indicies.Usi)] = make(map[uint16]bool)
						}

						_, has_subject_id := encode_resource.IsSchedIdxToSubIdToSkip[uint16(indicies.Usi)][subject_id]

						if !has_subject_id {
							encode_resource.IsSchedIdxToSubIdToSkip[uint16(indicies.Usi)][subject_id] = true
							selected_instructor.AssignedSubjects++
						}

						selected_instructor.TotalTeachingHours += (1.0 / Const.N_HOUR_TIME_SLOTS)
					}
				} // ------------- end of time_slot loop -------------
			} // ------------- end of day loop -------------

			return IterProceed
		},
	)

	return encode_resource, nil
}

// reads the default values of `EncodingResource` saved in a persistence instance.
func ReadDefaultEncodingResource(resource_persistence *StorageResources.Persistence) (*EncodingResource, error) {

	rooms, err_read_rooms := resource_persistence.ReaderService.ReadAllRooms()

	if err_read_rooms != nil {
		return nil, err_read_rooms
	}

	dept_id_to_room_type_to_rooms := GenerateMapDeptIdToRoomTypeToRooms(rooms)

	instructors, err_read_instructors := resource_persistence.ReaderService.ReadAllInstructors()

	if err_read_instructors != nil {
		return nil, err_read_instructors
	}

	dept_id_to_instructors := GenerateMapDeptIdToInstructors(instructors)

	//////////////////////////////////////////////////////////////////////////////////////
	//                              FLATTEN ENCODING RESOURCES
	//////////////////////////////////////////////////////////////////////////////////////

	room_id_to_room := make(map[uint16]*Rooms.Room)

	for out_key, out_v := range dept_id_to_room_type_to_rooms {
		for in_key, in_v := range out_v {
			for room_idx, room := range in_v {
				room_id_to_room[room.RoomID] = &dept_id_to_room_type_to_rooms[out_key][in_key][room_idx]
			}
		}
	}

	instructor_id_to_instructor := make(map[uint16]*Instructors.Instructor)

	for k, v := range dept_id_to_instructors {
		for instructor_idx, instructor := range v {
			instructor_id_to_instructor[instructor.InstructorID] = &dept_id_to_instructors[k][instructor_idx]
		}
	}

	////////////////////////////////////////////////////////////////////////////////////////

	return &EncodingResource{
		IsSchedIdxToSubIdToSkip: make(map[uint16]map[uint16]bool),
		DeptIdToInstructors:     dept_id_to_instructors,
		DeptIdToRoomtypeToRooms: dept_id_to_room_type_to_rooms,

		IdToInstructor: instructor_id_to_instructor,
		IdToRoom:       room_id_to_room,
	}, nil
}
