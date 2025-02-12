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
	TotalTeachingHours float32

	// determines if the instructor is available for a certain time slot during schedule generation.
	//
	// if value corresponding bit value for that time slot is 1, then the instructor is available.
	// if value corresponding bit value for that time slot is 0, then the instructor is NOT available.
	Time InstructorTimeSlotBitMap
}

type InstructorWithTimeString struct {
	InstructorID  uint16
	DepartmentID  uint16   `json:"DepartmentID"`
	FirstName     string   `json:"FirstName"`
	MiddleInitial string   `json:"MiddleInitial"`
	LastName      string   `json:"LastName"`
	Time          []string `json:"Time"`
}
