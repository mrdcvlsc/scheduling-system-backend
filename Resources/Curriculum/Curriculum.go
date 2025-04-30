package Curriculum

type Curriculum struct {
	CurriculumID   uint16      `json:"CurriculumID" bson:"CurriculumID"`     // non negative & non-zero unique number even if for example we have Computer Science (OLD) and Computer Science (New)
	CurriculumName string      `json:"CurriculumName" bson:"CurriculumName"` // e.g. Computer Science, Information Technology
	CurriculumCode string      `json:"CurriculumCode" bson:"CurriculumCode"` // e.g. BSCS, BSIT
	DepartmentID   uint16      `json:"DepartmentID" bson:"DepartmentID"`
	YearLevels     []YearLevel `json:"YearLevels" bson:"YearLevels"`
}

type YearLevel struct {
	Name      string     `json:"Name" bson:"Name"`
	IsActive  bool       `json:"IsActive" bson:"IsActive"`
	Semesters []Semester `json:"Semesters" bson:"Semesters"`
}

type Semester struct {
	Name     string    `json:"Name" bson:"Name"`
	Sections int       `json:"Sections" bson:"Sections"`
	Subjects []Subject `json:"Subjects" bson:"Subjects"`
}

var SEMESTER_INDEX_NAME [6]string = [6]string{
	"1st semester",
	"2nd semester",
	"3rd semester",
	"4th semester",
	"5th semester",
	"6th semester",
}

const SUPPORTED_SEMESTERS int = 2

var SECTION [26]string = [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func GetTotalNumberOfSections(curriculums []Curriculum, selected_semester int) int {
	section_count := 0

	for _, curriculum := range curriculums {
		for _, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue // skip inactive year levels
			}

			for semester_idx, semester := range year_level.Semesters {

				if selected_semester != semester_idx {
					continue // skip not selected semesters
				}

				section_count += semester.Sections
			} // ------------- end of semester_idx loop -------------
		} // ------------- end of year_level loop -------------
	} // ------------- end of curriculum loop -------------

	return section_count
}
