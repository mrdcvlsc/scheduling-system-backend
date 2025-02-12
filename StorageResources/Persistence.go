package StorageResources

import (
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

type readerRepository interface {
	// the order of curriculums returned by this method is always sorted by curriculum ID.
	ReadAllCurriculum() ([]Curriculum.Curriculum, error)
	ReadAllSubjects() ([]Curriculum.Subject, error)
	ReadAllDepartments() ([]Departments.Department, error)
	ReadAllInstructors() ([]Instructors.Instructor, error)
	ReadDepartmentInstructors(department_id int) ([]Instructors.Instructor, error)
	ReadAllRooms() ([]Rooms.Room, error)
}

type writerRepository interface {
	CreateInstructor(instructor Instructors.Instructor) error
	UpdateInstructor(instructor Instructors.Instructor) error
}

type Persistence struct {
	ReaderService readerRepository
	WriterService writerRepository
}
