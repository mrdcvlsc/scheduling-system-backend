package storage

import (
	// "fmt"
	// "testing"

	// storage "github.com/mrdcvlsc/scheduling-system-backend/Storage"
	curriculum "github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	instructor "github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
)

type persistenceRepository interface {
	GetAllCurriculum() []curriculum.Curriculum
	GetAllInstructors() []instructor.Instructor
}

type PersistenceService struct {
	Service persistenceRepository
}
