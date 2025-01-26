package Storage

import (
	// "fmt"
	// "testing"

	// storage "github.com/mrdcvlsc/scheduling-system-backend/Storage"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

type readerRepository interface {
	GetAllSubjects() ([]Curriculum.Subject, error)
	GetAllCurriculum() ([]Curriculum.Curriculum, error)
	GetAllDepartments() ([]Departments.Department, error)
	GetAllInstructors() ([]Instructors.Instructor, error)
	GetAllRooms() ([]Rooms.Room, error)
}

type writerRepository interface {
	SaveSchedules() error
}

type PersistenceService struct {
	ReaderService readerRepository
	WriterService writerRepository
}
