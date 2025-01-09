package storage_test

import (
	"fmt"
	"testing"

	p "github.com/mrdcvlsc/scheduling-system-backend/Storage"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func TestJsonFilePersistence_GetAllCurriculum(t *testing.T) {
	json_persistence := p.JsonFilePersistence{}
	curriculums := json_persistence.GetAllCurriculum()
	Utils.PrettyPrint(curriculums)
}

func TestJsonFilePersistence_GetAllInstructor(t *testing.T) {
	json_persistence := p.JsonFilePersistence{}
	instructors := json_persistence.GetAllInstructors()
	Utils.PrettyPrint(instructors)
	fmt.Println("\n\nTotal Number of Instructors : ", len(instructors))
}
