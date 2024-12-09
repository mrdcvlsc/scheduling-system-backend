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

	// determines if the instructor is available for a certain time slot during schedule generation.
	//
	// if value corresponding bit value for that time slot is 1, then the instructor is available.
	// if value corresponding bit value for that time slot is 0, then the instructor is NOT available.
	TimeSlotAvailability InstructorWeekTimeSlotAvailabilityMap

	// DayAvailability, each bits represent the boolean value if the instructor
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
