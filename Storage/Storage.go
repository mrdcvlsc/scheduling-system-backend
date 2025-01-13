package Storage

import (
	// "fmt"
	// "testing"

	// storage "github.com/mrdcvlsc/scheduling-system-backend/Storage"
	curriculum "github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	instructor "github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	room "github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

type persistenceRepository interface {
	GetAllCurriculum() []curriculum.Curriculum
	GetAllInstructors() []instructor.Instructor
	GetAllRooms() []room.Room
}

type PersistenceService struct {
	Service persistenceRepository
}
