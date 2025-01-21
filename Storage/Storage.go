package Storage

import (
	// "fmt"
	// "testing"

	// storage "github.com/mrdcvlsc/scheduling-system-backend/Storage"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

type persistenceRepository interface {
	GetAllSubjects() []Curriculum.Subject
	GetAllCurriculum() []Curriculum.Curriculum
	GetAllInstructors() []Instructors.Instructor
	GetAllRooms() []Rooms.Room
}

type PersistenceService struct {
	Service persistenceRepository
}
