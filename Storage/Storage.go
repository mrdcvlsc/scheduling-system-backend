package storage

import (
	// "fmt"
	// "testing"

	// storage "github.com/mrdcvlsc/scheduling-system-backend/Storage"
	c "github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
)

type persistenceRepository interface {
	GetAllCurriculum() []c.Curriculum
}

type PersistenceService struct {
	Service persistenceRepository
}
