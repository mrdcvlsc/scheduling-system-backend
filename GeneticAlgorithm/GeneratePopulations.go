package geneticalgorithm

import (
	// "fmt"
	// "testing"

	sd "github.com/mrdcvlsc/scheduling-system-backend/Resources/ScheduleDSA"
	persistence "github.com/mrdcvlsc/scheduling-system-backend/Storage"
)

func NewPopulation() sd.SchedulesInUniversity {
	// TODO: get all courses curriculum in the university
	p := persistence.PersistenceService{Service: &persistence.JsonFilePersistence{}}
	// TODO: get all sections / classes in the university
	// TODO: get all subjects
	return nil
}
