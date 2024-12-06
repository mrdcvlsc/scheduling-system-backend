package main

import (
	"encoding/json"
	"fmt"
	"testing"

	ga "github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
)

func TestGeneticAlgorithmBasics(t *testing.T) {
	fmt.Println("Test/Dev Program for genetic algorithm scheduling")

	uniSectionSchedules := ga.NewSectionsWeeklySchedules(1)

	uniSectionSchedules[0][0][1] = *ga.NewTimeSlot(222, 333, 444)

	fmt.Println("================================================================")
	testSectionSchedule := uniSectionSchedules.GetSectionWeeklySchedule(0)
	fmt.Printf("len(testSectionSchedule) = %d\n", len(testSectionSchedule))
	fmt.Println("================================================================")

	testDaySchedule := testSectionSchedule.GetDaySchedule(0)
	fmt.Printf("len(testDaySchedule) = %d\n", len(testDaySchedule))

	fmt.Println("================================================================")

	testTimeSlot := testDaySchedule.GetTimeSlot(1)

	fmt.Printf("testTimeSlot = %+v\n", testTimeSlot)

	testTimeSlot.ToggleCFlagSubjectID()
	testTimeSlot.ToggleCFlagInstructorID()
	testTimeSlot.ToggleCFlagRoomID()

	jsonFormat_testTimeSlot, _ := json.MarshalIndent(testTimeSlot, "", " ")
	fmt.Printf("testTimeSlot = %s\n", jsonFormat_testTimeSlot)
	fmt.Printf("testTimeSlot = %+v\n", testTimeSlot)

	fmt.Println("================================================================")

	jsonFormat_uniSectionSchedules, errJson := json.MarshalIndent(uniSectionSchedules, "", " ")
	fmt.Printf("uniSectionSchedules = %s\n", jsonFormat_uniSectionSchedules)

	fmt.Println("================================================================")

	if uniSectionSchedules[0][0][1].IsConstInstructorID() != true {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if uniSectionSchedules[0][0][1].IsConstRoomID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if uniSectionSchedules[0][0][1].IsConstSubjectID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	if uniSectionSchedules[0][0][1].GetInstructorID() != 333 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 333")
	}

	if uniSectionSchedules[0][0][1].GetRoomID() != 444 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 444")
	}

	if uniSectionSchedules[0][0][1].GetSubjectID() != 222 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 222")
	}

	fmt.Println("================================================================")

	testTimeSlot.SetInstructorID(23423)
	testTimeSlot.SetRoomID(12122)
	testTimeSlot.SetSubjectID(9898)

	if uniSectionSchedules[0][0][1].GetInstructorID() != 23423 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 23423")
	}

	if uniSectionSchedules[0][0][1].GetRoomID() != 12122 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 12122")
	}

	if uniSectionSchedules[0][0][1].GetSubjectID() != 9898 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 9898")
	}

	if uniSectionSchedules[0][0][1].IsConstInstructorID() != true {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if uniSectionSchedules[0][0][1].IsConstRoomID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if uniSectionSchedules[0][0][1].IsConstSubjectID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	fmt.Println("================================================================")

	testTimeSlot.ToggleCFlagSubjectID()
	testTimeSlot.ToggleCFlagInstructorID()
	testTimeSlot.ToggleCFlagRoomID()

	if uniSectionSchedules[0][0][1].IsConstInstructorID() != false {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if uniSectionSchedules[0][0][1].IsConstRoomID() != false {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if uniSectionSchedules[0][0][1].IsConstSubjectID() != false {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	if uniSectionSchedules[0][0][1].GetInstructorID() != 23423 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 23423")
	}

	if uniSectionSchedules[0][0][1].GetRoomID() != 12122 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 12122")
	}

	if uniSectionSchedules[0][0][1].GetSubjectID() != 9898 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 9898")
	}

	if errJson != nil {
		t.Error(errJson)
	}
}
