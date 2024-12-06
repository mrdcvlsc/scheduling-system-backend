package geneticalgorithm

import "fmt"

// A bitmask to identify the constraint flag of an attribute.
// When this bit is set, the associated attribute is considered constrained (e.g., fixed or locked).
const ATTRIBUTE_CONSTRAINT_MASK uint16 = 1 << 15 // 0b1000000000000000

// A bitmask to extract the actual value of an attribute by clearing the constraint flag bit.
const ATTRIBUTE_VALUE_MASK uint16 = ^ATTRIBUTE_CONSTRAINT_MASK // 0b0111111111111111

// Defines the maximum valid value for an attribute, derived from ATTRIBUTE_VALUE_MASK.
const ATTRIBUTE_MAX_VALUE = ATTRIBUTE_VALUE_MASK

/*
Represents a scheduling unit containing a subject, instructor, and room.

Each attribute is encoded as a 16-bit unsigned integer with the most significant
bit used as a constraint flag. When this bit is set to 1, the associated attribute
is considered constrained (e.g., fixed or locked), and will not be modified by the
genetic algorithm.

The position of the time slot is also fixed and will not be moved to other time slots
by the genetic algorithm if one attribute is fixed/constrained.
*/
type TimeSlot struct {
	subjectID    uint16
	instructorID uint16
	roomID       uint16
}

// ============================= CONSTRUCTOR =============================

func NewTimeSlot(subject_id, instructor_id, room_id uint16) *TimeSlot {
	return &TimeSlot{
		subjectID:    subject_id,
		instructorID: instructor_id,
		roomID:       room_id,
	}
}

// ============================= GET COPY WITHOUT CONSTRAIN FLAGS =============================

// Returns a copy of the TimeSlot with all constraint flags cleared.
// This is useful for obtaining the pure attribute values without the constraint metadata.
func (time_slot TimeSlot) CopyWithoutFlags() TimeSlot {
	time_slot.subjectID = time_slot.subjectID & ATTRIBUTE_VALUE_MASK
	time_slot.instructorID = time_slot.instructorID & ATTRIBUTE_VALUE_MASK
	time_slot.roomID = time_slot.roomID & ATTRIBUTE_VALUE_MASK
	return time_slot
}

// ============================= GET VALUE METHODS =============================

// Retrieves the subject ID value of the TimeSlot, excluding the constraint flag.
func (time_slot *TimeSlot) GetSubjectID() uint16 {
	return ATTRIBUTE_VALUE_MASK & time_slot.subjectID
}

// Retrieves the instructor ID value of the TimeSlot, excluding the constraint flag.
func (time_slot *TimeSlot) GetInstructorID() uint16 {
	return ATTRIBUTE_VALUE_MASK & time_slot.instructorID
}

// Retrieves the room ID value of the TimeSlot, excluding the constraint flag.
func (time_slot *TimeSlot) GetRoomID() uint16 {
	return ATTRIBUTE_VALUE_MASK & time_slot.roomID
}

// ============================= SET VALUE METHODS =============================

// Updates the subject ID value of the TimeSlot.
// If the value exceeds ATTRIBUTE_MAX_VALUE, an error is returned.
func (time_slot *TimeSlot) SetSubjectID(subject_id uint16) error {
	if subject_id > ATTRIBUTE_MAX_VALUE {
		return fmt.Errorf("SetSubjectID(%d) : Invalid Parameter Value", subject_id)
	}

	time_slot.subjectID = (time_slot.subjectID & ATTRIBUTE_CONSTRAINT_MASK) | subject_id
	return nil
}

// Updates the instructor ID value of the TimeSlot.
// If the value exceeds ATTRIBUTE_MAX_VALUE, an error is returned.
func (time_slot *TimeSlot) SetInstructorID(instructor_id uint16) error {
	if instructor_id > ATTRIBUTE_MAX_VALUE {
		return fmt.Errorf("SetInstructorID(%d) : Invalid Parameter Value", instructor_id)
	}

	time_slot.instructorID = (time_slot.instructorID & ATTRIBUTE_CONSTRAINT_MASK) | instructor_id
	return nil
}

// Updates the room ID value of the TimeSlot.
// If the value exceeds ATTRIBUTE_MAX_VALUE, an error is returned.
func (time_slot *TimeSlot) SetRoomID(room_id uint16) error {
	if room_id > ATTRIBUTE_MAX_VALUE {
		return fmt.Errorf("SetRoomID(%d) : Invalid Parameter Value", room_id)
	}

	time_slot.roomID = (time_slot.roomID & ATTRIBUTE_CONSTRAINT_MASK) | room_id
	return nil
}

// ============================= GET CONSTRAINT FLAG METHODS =============================

// Checks whether the constraint flag is set for the subject ID.
// When this bit is set, the associated attribute is considered
// constrained (e.g., fixed or locked) by the genetic algorithm.
func (time_slot *TimeSlot) IsConstSubjectID() bool {
	return ATTRIBUTE_CONSTRAINT_MASK&time_slot.subjectID != 0
}

// Checks whether the constraint flag is set for the instructor ID.
// When this bit is set, the associated attribute is considered
// constrained (e.g., fixed or locked) by the genetic algorithm.
func (time_slot *TimeSlot) IsConstInstructorID() bool {
	return ATTRIBUTE_CONSTRAINT_MASK&time_slot.instructorID != 0
}

// Checks whether the constraint flag is set for the room ID.
// When this bit is set, the associated attribute is considered
// constrained (e.g., fixed or locked) by the genetic algorithm.
func (time_slot *TimeSlot) IsConstRoomID() bool {
	return ATTRIBUTE_CONSTRAINT_MASK&time_slot.roomID != 0
}

// ============================= TOGGLE CONSTRAINT FLAG METHODS =============================

// Flips the constraint flag bit for the subject ID.
func (time_slot *TimeSlot) ToggleCFlagSubjectID() {
	constraint_flag := time_slot.subjectID & ATTRIBUTE_CONSTRAINT_MASK
	constraint_flag = ^constraint_flag
	constraint_flag &= ATTRIBUTE_CONSTRAINT_MASK

	time_slot.subjectID = (time_slot.subjectID & ATTRIBUTE_VALUE_MASK) | constraint_flag
}

// Flips the constraint flag bit for the instructor ID.
func (time_slot *TimeSlot) ToggleCFlagInstructorID() {
	constraint_flag := time_slot.instructorID & ATTRIBUTE_CONSTRAINT_MASK
	constraint_flag = ^constraint_flag
	constraint_flag &= ATTRIBUTE_CONSTRAINT_MASK

	time_slot.instructorID = (time_slot.instructorID & ATTRIBUTE_VALUE_MASK) | constraint_flag
}

// Flips the constraint flag bit for the room ID.
func (time_slot *TimeSlot) ToggleCFlagRoomID() {
	constraint_flag := time_slot.roomID & ATTRIBUTE_CONSTRAINT_MASK
	constraint_flag = ^constraint_flag
	constraint_flag &= ATTRIBUTE_CONSTRAINT_MASK

	time_slot.roomID = (time_slot.roomID & ATTRIBUTE_VALUE_MASK) | constraint_flag
}
