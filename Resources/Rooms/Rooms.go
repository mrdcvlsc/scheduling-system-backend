package Rooms

import "github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"

// the maximum room capacity.
const MAX_ROOM_CAPACITY int = 15 // = 0b1111 (4 bits only)

// TODO: (we only need 4 bits to store up to 15 hours) - Test with array dim : [6][24]
const TIME_SLOT_CLASS_COUNTER_SIZE int = (Const.N_DAILY_TIME_SLOTS * Const.N_WEEKLY_SCHOOL_DAYS) / 2

type Room struct {
	RoomID             uint16
	DepartmentID       uint16                              `json:"DepartmentID"`
	Capacity           uint16                              `json:"Capacity"` // maximum numbers of classes or sections a room can hold in a single time slot.
	RoomType           uint16                              `json:"RoomType"` // determines the room type.
	Name               string                              `json:"Name"`
	timeSlotClassCount [TIME_SLOT_CLASS_COUNTER_SIZE]uint8 // records the numbers of classes or sections allocated in the room for a specific timeslot
}

// set the current number of classes or sections allocated in the room for a specific time slot.
func (room *Room) SetTimeSlotClassCount(day, time_slot int, class_count uint8) {
	idx_2D_to_1D := (day * Const.N_DAILY_TIME_SLOTS) + time_slot

	// in one byte or uint8 we can use the two set of 4-bits (higher and lower) to store two class or
	// section count for a specific time slot, hence we divide by the number of bits of uint8 to 2.
	limb_idx := idx_2D_to_1D / 2
	shift_multiplier := idx_2D_to_1D % 2

	room.timeSlotClassCount[limb_idx] &= (0b11110000 >> (4 * shift_multiplier))
	room.timeSlotClassCount[limb_idx] |= class_count << (4 * shift_multiplier)
}

// increase by 1 the current number of classes or sections allocated in the room for a specific time slot.
//
// warning incrementing until the max room capacity 15 would overflow the whole
// uint8 which cause serialized data corruption due to overflow of the uint8 type.
func (room *Room) IncTimeSlotClassCount(day, time_slot int) {
	idx_2D_to_1D := (day * Const.N_DAILY_TIME_SLOTS) + time_slot
	limb_idx := idx_2D_to_1D / 2
	shift_multiplier := idx_2D_to_1D % 2

	room.timeSlotClassCount[limb_idx] += (0b1 << (4 * shift_multiplier))
}

// get the current number of classes or sections allocated in the room for a specific time slot.
func (room *Room) GetTimeSlotClassCount(day, time_slot int) uint8 {
	idx_2D_to_1D := (day * Const.N_DAILY_TIME_SLOTS) + time_slot
	limb_idx := idx_2D_to_1D / 2
	shift_multiplier := idx_2D_to_1D % 2

	return (room.timeSlotClassCount[limb_idx] & (0b1111 << (4 * shift_multiplier))) >> (4 * shift_multiplier)
}
