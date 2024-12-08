package geneticalgorithm

const DAY_MON uint8 = 1 << 5
const DAY_TUE uint8 = 1 << 4
const DAY_WED uint8 = 1 << 3
const DAY_THR uint8 = 1 << 2
const DAY_FRI uint8 = 1 << 1
const DAY_SAT uint8 = 1
const WHOLE_WEEK uint8 = DAY_MON | DAY_TUE | DAY_WED | DAY_THR | DAY_FRI | DAY_SAT

type InstructorMonitor struct {
	// this should never be zero, zero means empty, none or nothing.
	InstructorID uint16

	// determines if the instructor is available in the current timeslot during schedule generation.
	//
	// if value is 1, then the instructor is available.
	// if less than 1, the instructor is not available.
	//
	// This attribute should be incremented everytime we move one time slot
	// forward during the instructor assignement phase in schedule generation.
	// If the value is already 1 this attribute should not be incremented.
	TimeSlotAvailability int8

	// TODO: use bitset/bitmaps for recording time slot availability of instructors.
	// by using bitmap you can remove the instructor assignment phase that is isolated
	// from subject assignment phase, which could lead to faster performance at the cost
	// of slightly higer memory consumption.

	// Available days, each bits represent the boolean value if the instructor
	// is available.
	//
	// The least significant bit represent SAT then the preceding bits from
	// the least significant to the most significant bits are:
	//
	// FRI, THR, WED, TUE and MON.
	//
	// The first 2 most significant bits are unused.
	//
	// BIT FORMAT: [2 bits unused][1-bit MON][1-bit TUE][1-bit WED][1-bit THR][1-bit FRI][1-bit SAT]
	DayAvailability uint8
}

type InstructorDbRecord struct {
	InstructorID    uint16
	DayAvailability uint8
	DepartmentID    uint8
	InstructorName  string
}
