package StorageReader_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/StorageReader"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

var TestPersistence StorageReader.Persistence

func TestJsonFilePersistence_GetAllSubjects(t *testing.T) {
	TestPersistence.Service = &StorageReader.JsonReader{}
	subjects, err := TestPersistence.Service.GetAllSubjects()

	if err != nil {
		t.Error(err)
	}

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

	empty_strings := Utils.CheckForEmptyStrings(subjects, "subjects")

	if len(empty_strings) > 0 {
		for _, err := range empty_strings {
			t.Error("empty string :", err)
			fmt.Println()
		}
	}
}

func TestJsonFilePersistence_GetAllCurriculum(t *testing.T) {
	TestPersistence.Service = &StorageReader.JsonReader{}
	curriculums, err := TestPersistence.Service.GetAllCurriculum()

	if err != nil {
		t.Error(err)
	}

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

	empty_strings := Utils.CheckForEmptyStrings(curriculums, "curriculums")

	if len(empty_strings) > 0 {
		for _, err := range empty_strings {
			t.Error("empty string :", err)
			fmt.Println()
		}
	}
}

func TestJsonFilePersistence_GetAllInstructor(t *testing.T) {
	TestPersistence.Service = &StorageReader.JsonReader{}
	instructors, err := TestPersistence.Service.GetAllInstructors()

	if err != nil {
		t.Error(err)
	}

	empty_strings := Utils.CheckForEmptyStrings(instructors, "instructors")

	if len(empty_strings) > 0 {
		for _, err := range empty_strings {
			t.Error("empty string :", err)
			fmt.Println()
		}
	}

	fmt.Println("\n\nTotal Number of Instructors : ", len(instructors))
}

func TestJsonFilePersistence_GetAllRoom(t *testing.T) {
	TestPersistence.Service = &StorageReader.JsonReader{}
	rooms, err := TestPersistence.Service.GetAllRooms()

	if err != nil {
		t.Error(err)
	}

	empty_strings := Utils.CheckForEmptyStrings(rooms, "rooms")

	if len(empty_strings) > 0 {
		for _, err := range empty_strings {
			t.Error("empty string :", err)
			fmt.Println()
		}
	}

	fmt.Println("\n\nTotal Number of Rooms : ", len(rooms))
}

func TestJsonFilePersistence_GetDepartments(t *testing.T) {
	TestPersistence.Service = &StorageReader.JsonReader{}
	departments, err := TestPersistence.Service.GetAllDepartments()

	if err != nil {
		t.Error(err)
	}

	empty_strings := Utils.CheckForEmptyStrings(departments, "departments")

	if len(empty_strings) > 0 {
		for _, err := range empty_strings {
			t.Error("empty string :", err)
			fmt.Println()
		}
	}

	fmt.Println("\n\nTotal Number of Rooms : ", len(departments))
}
