package StorageResources_test

import (
	"fmt"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func TestJsonFilePersistence_Update(t *testing.T) {
	TestPersistence.ReaderService = &StorageResources.JsonReader{}
	TestPersistence.WriterService = &StorageResources.JsonWriter{}

	instructors, err := TestPersistence.ReaderService.ReadAllInstructors()

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

	backup_instructor := instructors[2]

	instructors[2].FirstName = "Carcharodon"
	instructors[2].LastName = "Astra"

	updated_instructor := instructors[2]

	if backup_instructor == instructors[2] {
		t.Error("that should not be equal anymore")
	}

	err_create := TestPersistence.WriterService.CreateInstructor(instructors[2])

	if err_create == nil {
		t.Error("there should be an error there")
	}

	err_update := TestPersistence.WriterService.UpdateInstructor(instructors[2])

	if err_update != nil {
		t.Error(err_update)
	}

	new_instructor := Instructors.Instructor{
		DepartmentID:  1,
		FirstName:     "Mephiston",
		MiddleInitial: "E",
		LastName:      "Calistarius",
		Time:          [3]uint64{1, 3},
	}

	err_create_2 := TestPersistence.WriterService.CreateInstructor(new_instructor)

	if err_create_2 != nil {
		t.Error(err_create_2)
	}

	instructors, err = TestPersistence.ReaderService.ReadAllInstructors()

	if err != nil {
		t.Error(err)
	}

	has_found_new_instructor := false
	has_updated := false

	for _, instructor := range instructors {
		if instructor.FirstName == new_instructor.FirstName {
			new_instructor.InstructorID = instructor.InstructorID
			if instructor == new_instructor {
				has_found_new_instructor = true
			}
		}

		if updated_instructor == instructor {
			has_updated = true
		}
	}

	if !has_found_new_instructor {
		t.Error("New instructor created earlier not found")
	}

	if !has_updated {
		t.Error("there is no updated instructor")
	}

	fmt.Println("\n\nTotal Number of Instructors : ", len(instructors))
}
