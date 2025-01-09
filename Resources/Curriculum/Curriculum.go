package curriculum

type Curriculum struct {
	CurriculumID   uint16      // unique number even if for example we have Computer Science (OLD) and Computer Science (New)
	CurriculumName string      `json:"CurriculumName"` // e.g. Computer Science, Information Technology
	CurriculumCode string      `json:"CurriculumCode"` // e.g. BSCS, BSIT
	DepartmentID   uint16      `json:"DepartmentID"`
	YearLevels     []YearLevel `json:"YearLevels"`
}

type YearLevel struct {
	Name      string     `json:"Name"`
	IsActive  bool       `json:"IsActive"`
	Semesters []Semester `json:"Semesters"`
}

type Semester struct {
	Name     string    `json:"Name"`
	Sections int       `json:"Sections"`
	Subjects []Subject `json:"Subjects"`
}
