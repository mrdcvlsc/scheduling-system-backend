package Storage_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/Storage"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

var TestPersistence Storage.PersistenceService

func TestJsonFilePersistence_GetAllSubjects(t *testing.T) {
	TestPersistence.Service = &Storage.JsonFilePersistence{}
	subjects := TestPersistence.Service.GetAllSubjects()

	for _, subject := range subjects {
		if strings.Contains(subject.Code, "FITT") {
			if ((subject.BitFlags & 1) != 1) || !subject.IsGymType() {
				t.Errorf("Subject %s - should have one 1 bit in the least significant bit", subject.Code)
			}
		} else {
			if ((subject.BitFlags & 1) != 0) || subject.IsGymType() {
				t.Errorf("Subject %s - should have a 0 bit in the least significant bit", subject.Code)
			}
		}
	}

	Utils.PrettyPrint(subjects)
}

func TestJsonFilePersistence_GetAllCurriculum(t *testing.T) {
	TestPersistence.Service = &Storage.JsonFilePersistence{}
	curriculums := TestPersistence.Service.GetAllCurriculum()

	for _, curriculum := range curriculums {
		for _, year_level := range curriculum.YearLevels {
			for _, semester := range year_level.Semesters {
				for _, subject := range semester.Subjects {
					if strings.Contains(subject.Code, "FITT") {
						if ((subject.BitFlags & 1) != 1) || !subject.IsGymType() {
							t.Errorf("Subject %s - should have one 1 bit in the least significant bit", subject.Code)
						}
					} else {
						if ((subject.BitFlags & 1) != 0) || subject.IsGymType() {
							t.Errorf("Subject %s - should have a 0 bit in the least significant bit", subject.Code)
						}
					}
				}
			}
		}
	}

	Utils.PrettyPrint(curriculums)
}

func TestJsonFilePersistence_GetAllInstructor(t *testing.T) {
	TestPersistence.Service = &Storage.JsonFilePersistence{}
	instructors := TestPersistence.Service.GetAllInstructors()
	Utils.PrettyPrint(instructors)
	fmt.Println("\n\nTotal Number of Instructors : ", len(instructors))
}

func TestJsonFilePersistence_GetAllRoom(t *testing.T) {
	TestPersistence.Service = &Storage.JsonFilePersistence{}
	rooms := TestPersistence.Service.GetAllRooms()
	Utils.PrettyPrint(rooms)
	fmt.Println("\n\nTotal Number of Rooms : ", len(rooms))
}
