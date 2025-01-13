package Storage_test

import (
	"fmt"
	"testing"

	storage "github.com/mrdcvlsc/scheduling-system-backend/Storage"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

var TestPersistence storage.PersistenceService

func TestJsonFilePersistence_GetAllCurriculum(t *testing.T) {
	TestPersistence.Service = &storage.JsonFilePersistence{}
	curriculums := TestPersistence.Service.GetAllCurriculum()
	Utils.PrettyPrint(curriculums)
}

func TestJsonFilePersistence_GetAllInstructor(t *testing.T) {
	TestPersistence.Service = &storage.JsonFilePersistence{}
	instructors := TestPersistence.Service.GetAllInstructors()
	Utils.PrettyPrint(instructors)
	fmt.Println("\n\nTotal Number of Instructors : ", len(instructors))
}

func TestJsonFilePersistence_GetAllRoom(t *testing.T) {
	TestPersistence.Service = &storage.JsonFilePersistence{}
	rooms := TestPersistence.Service.GetAllRooms()
	Utils.PrettyPrint(rooms)
	fmt.Println("\n\nTotal Number of Rooms : ", len(rooms))
}
