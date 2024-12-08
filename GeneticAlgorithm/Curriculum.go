package geneticalgorithm

type CurriculumSubject struct {
	SubjectID uint16

	// if there are no designated instructor IDs here, the algorithm
	// will assign random instructors from the department.
	//
	// if there are some designated instructor IDs here, the algorithm
	// will immediately assign the instructor to the allocated subject
	// time slot, in the InstructorMonitor
	DesignatedInstructorsID []uint16

	DesignatedRoomTypes []uint16
}

type Semester []CurriculumSubject

type CurriculumDbRecord struct {
	// unique number even if for example we have Computer Science (OLD) and Computer Science (New)
	CourseID uint8

	DepartmentID uint8

	// e.g. Computer Science, Information Technology
	CourseName string

	// e.g. BSCS, BSIT
	CourseCode string

	Semesters []Semester
}
