package Instructors

// type that will be use to read and save a single instructor's information from and to the database.
type Instructor struct {
	InstructorID  uint16 // this should never be zero, zero means empty, none or nothing.
	DepartmentID  uint16 `json:"DepartmentID"`
	FirstName     string `json:"FirstName"`
	MiddleInitial string `json:"MiddleInitial"`
	LastName      string `json:"LastName"`

	// the total number of assigned subjects to teach.
	AssignedSubjects int

	// the total teaching hours assigned for all assigned subjects.
	TotalTeachingHours int

	// AvailableDay, each bits represent the boolean value if the instructor
	// is available.
	//
	// The least significant bit represent MON then the preceding bits from
	// the least significant to the most significant bits are:
	//
	// TUE, WED, THR, FRI and SAT.
	//
	// The first 2 most significant bits are unused.
	//
	// BIT FORMAT: [2-bits unused][1-bit SAT][1-bit FRI][1-bit THR][1-bit WED][1-bit TUE][1-bit MON]
	AvailableDay uint8 `json:"AvailableDay"`

	// determines if the instructor is available for a certain time slot during schedule generation.
	//
	// if value corresponding bit value for that time slot is 1, then the instructor is available.
	// if value corresponding bit value for that time slot is 0, then the instructor is NOT available.
	Time InstructorTimeSlotBitMap
}

const DAY_MON uint8 = 1 << 5
const DAY_TUE uint8 = 1 << 4
const DAY_WED uint8 = 1 << 3
const DAY_THR uint8 = 1 << 2
const DAY_FRI uint8 = 1 << 1
const DAY_SAT uint8 = 1

const WHOLE_WEEK uint8 = DAY_MON | DAY_TUE | DAY_WED | DAY_THR | DAY_FRI | DAY_SAT
