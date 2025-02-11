package StorageResources

import (
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

type readerRepository interface {
	// the order of curriculums returned by this method is always sorted by curriculum ID.
	GetAllCurriculum() ([]Curriculum.Curriculum, error)
	GetAllSubjects() ([]Curriculum.Subject, error)
	GetAllDepartments() ([]Departments.Department, error)
	GetAllInstructors() ([]Instructors.Instructor, error)
	GetDepartmentInstructors(department_id int) ([]Instructors.Instructor, error)
	GetAllRooms() ([]Rooms.Room, error)
}

type Persistence struct {
	ReaderService readerRepository
}
