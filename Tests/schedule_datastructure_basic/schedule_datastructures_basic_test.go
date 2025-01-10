package schedule_datastructures_basic_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/schedule"
)

func TestScheduleTypes1(t *testing.T) {
	fmt.Println("Test/Dev Program for genetic algorithm scheduling")

	universitySchedules := schedule.NewUniTimeTables(1)

	universitySchedules[0][0][1].Set(222, 333, 444)

	fmt.Println("================================================================")
	testSectionSchedule := universitySchedules.Get(0)
	fmt.Printf("len(testSectionSchedule) = %d\n", len(testSectionSchedule))
	fmt.Println("================================================================")

	testDaySchedule := testSectionSchedule.Get(0)
	fmt.Printf("len(testDaySchedule) = %d\n", len(testDaySchedule))

	fmt.Println("================================================================")

	testTimeSlot := testDaySchedule.Get(1)

	fmt.Printf("testTimeSlot = %+v\n", testTimeSlot)

	testTimeSlot.ToggleCFlagSubjectID()
	testTimeSlot.ToggleCFlagInstructorID()
	testTimeSlot.ToggleCFlagRoomID()

	jsonFormat_testTimeSlot, _ := json.MarshalIndent(testTimeSlot, "", " ")
	fmt.Printf("testTimeSlot = %s\n", jsonFormat_testTimeSlot)
	fmt.Printf("testTimeSlot = %+v\n", testTimeSlot)

	fmt.Println("================================================================")

	jsonFormat_uniSectionSchedules, errJson := json.MarshalIndent(universitySchedules, "", " ")
	fmt.Printf("uniSectionSchedules = %s\n", jsonFormat_uniSectionSchedules)

	fmt.Println("================================================================")

	if universitySchedules[0][0][1].IsConstInstructorID() != true {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstRoomID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstSubjectID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	if universitySchedules[0][0][1].GetInstructorID() != 333 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 333")
	}

	if universitySchedules[0][0][1].GetRoomID() != 444 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 444")
	}

	if universitySchedules[0][0][1].GetSubjectID() != 222 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 222")
	}

	fmt.Println("================================================================")

	testTimeSlot.SetInstructorID(23423)
	testTimeSlot.SetRoomID(12122)
	testTimeSlot.SetSubjectID(9898)

	if universitySchedules[0][0][1].GetInstructorID() != 23423 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 23423")
	}

	if universitySchedules[0][0][1].GetRoomID() != 12122 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 12122")
	}

	if universitySchedules[0][0][1].GetSubjectID() != 9898 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 9898")
	}

	if universitySchedules[0][0][1].IsConstInstructorID() != true {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstRoomID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstSubjectID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	fmt.Println("================================================================")

	testTimeSlot.ToggleCFlagSubjectID()
	testTimeSlot.ToggleCFlagInstructorID()
	testTimeSlot.ToggleCFlagRoomID()

	if universitySchedules[0][0][1].IsConstInstructorID() != false {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstRoomID() != false {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstSubjectID() != false {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	if universitySchedules[0][0][1].GetInstructorID() != 23423 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 23423")
	}

	if universitySchedules[0][0][1].GetRoomID() != 12122 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 12122")
	}

	if universitySchedules[0][0][1].GetSubjectID() != 9898 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 9898")
	}

	if errJson != nil {
		t.Error(errJson)
	}
}

func TestScheduleTypes1_5(t *testing.T) {
	fmt.Println("Test/Dev Program for genetic algorithm scheduling")

	universitySchedules := schedule.NewUniTimeTables(1)

	universitySchedules[0][0][1].SetSubjectID(222)
	universitySchedules[0][0][1].SetInstructorID(333)
	universitySchedules[0][0][1].SetRoomID(444)

	fmt.Println("================================================================")
	testSectionSchedule := universitySchedules.Get(0)
	fmt.Printf("len(testSectionSchedule) = %d\n", len(testSectionSchedule))
	fmt.Println("================================================================")

	testDaySchedule := testSectionSchedule.Get(0)
	fmt.Printf("len(testDaySchedule) = %d\n", len(testDaySchedule))

	fmt.Println("================================================================")

	testTimeSlot := testDaySchedule.Get(1)

	fmt.Printf("testTimeSlot = %+v\n", testTimeSlot)

	testTimeSlot.ToggleCFlagSubjectID()
	testTimeSlot.ToggleCFlagInstructorID()
	testTimeSlot.ToggleCFlagRoomID()

	jsonFormat_testTimeSlot, _ := json.MarshalIndent(testTimeSlot, "", " ")
	fmt.Printf("testTimeSlot = %s\n", jsonFormat_testTimeSlot)
	fmt.Printf("testTimeSlot = %+v\n", testTimeSlot)

	fmt.Println("================================================================")

	jsonFormat_uniSectionSchedules, errJson := json.MarshalIndent(universitySchedules, "", " ")
	fmt.Printf("uniSectionSchedules = %s\n", jsonFormat_uniSectionSchedules)

	fmt.Println("================================================================")

	if universitySchedules[0][0][1].IsConstInstructorID() != true {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstRoomID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstSubjectID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	if universitySchedules[0][0][1].GetInstructorID() != 333 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 333")
	}

	if universitySchedules[0][0][1].GetRoomID() != 444 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 444")
	}

	if universitySchedules[0][0][1].GetSubjectID() != 222 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 222")
	}

	fmt.Println("================================================================")

	testTimeSlot.SetInstructorID(23423)
	testTimeSlot.SetRoomID(12122)
	testTimeSlot.SetSubjectID(9898)

	if universitySchedules[0][0][1].GetInstructorID() != 23423 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 23423")
	}

	if universitySchedules[0][0][1].GetRoomID() != 12122 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 12122")
	}

	if universitySchedules[0][0][1].GetSubjectID() != 9898 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 9898")
	}

	if universitySchedules[0][0][1].IsConstInstructorID() != true {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstRoomID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstSubjectID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	fmt.Println("================================================================")

	testTimeSlot.ToggleCFlagSubjectID()
	testTimeSlot.ToggleCFlagInstructorID()
	testTimeSlot.ToggleCFlagRoomID()

	if universitySchedules[0][0][1].IsConstInstructorID() != false {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstRoomID() != false {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules[0][0][1].IsConstSubjectID() != false {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	if universitySchedules[0][0][1].GetInstructorID() != 23423 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 23423")
	}

	if universitySchedules[0][0][1].GetRoomID() != 12122 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 12122")
	}

	if universitySchedules[0][0][1].GetSubjectID() != 9898 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 9898")
	}

	if errJson != nil {
		t.Error(errJson)
	}
}

func TestScheduleTypes2(t *testing.T) {
	fmt.Println("Test/Dev Program for genetic algorithm scheduling")

	universitySchedules := schedule.NewUniTimeTables(1)

	universitySchedules.Get(0).Get(0).Get(1).Set(222, 333, 444)

	fmt.Println("================================================================")
	testSectionSchedule := universitySchedules.Get(0)
	fmt.Printf("len(testSectionSchedule) = %d\n", len(testSectionSchedule))
	fmt.Println("================================================================")

	testDaySchedule := testSectionSchedule.Get(0)
	fmt.Printf("len(testDaySchedule) = %d\n", len(testDaySchedule))

	fmt.Println("================================================================")

	testTimeSlot := testDaySchedule.Get(1)

	fmt.Printf("testTimeSlot = %+v\n", testTimeSlot)

	testTimeSlot.ToggleCFlagSubjectID()
	testTimeSlot.ToggleCFlagInstructorID()
	testTimeSlot.ToggleCFlagRoomID()

	jsonFormat_testTimeSlot, _ := json.MarshalIndent(testTimeSlot, "", " ")
	fmt.Printf("testTimeSlot = %s\n", jsonFormat_testTimeSlot)
	fmt.Printf("testTimeSlot = %+v\n", testTimeSlot)

	fmt.Println("================================================================")

	jsonFormat_uniSectionSchedules, errJson := json.MarshalIndent(universitySchedules, "", " ")
	fmt.Printf("uniSectionSchedules = %s\n", jsonFormat_uniSectionSchedules)

	fmt.Println("================================================================")

	if universitySchedules.Get(0).Get(0).Get(1).IsConstInstructorID() != true {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules.Get(0).Get(0).Get(1).IsConstRoomID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules.Get(0).Get(0).Get(1).IsConstSubjectID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	if universitySchedules.Get(0).Get(0).Get(1).GetInstructorID() != 333 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 333")
	}

	if universitySchedules.Get(0).Get(0).Get(1).GetRoomID() != 444 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 444")
	}

	if universitySchedules.Get(0).Get(0).Get(1).GetSubjectID() != 222 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 222")
	}

	fmt.Println("================================================================")

	testTimeSlot.SetInstructorID(23423)
	testTimeSlot.SetRoomID(12122)
	testTimeSlot.SetSubjectID(9898)

	if universitySchedules.Get(0).Get(0).Get(1).GetInstructorID() != 23423 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 23423")
	}

	if universitySchedules.Get(0).Get(0).Get(1).GetRoomID() != 12122 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 12122")
	}

	if universitySchedules.Get(0).Get(0).Get(1).GetSubjectID() != 9898 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 9898")
	}

	if universitySchedules.Get(0).Get(0).Get(1).IsConstInstructorID() != true {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules.Get(0).Get(0).Get(1).IsConstRoomID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules.Get(0).Get(0).Get(1).IsConstSubjectID() != true {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	fmt.Println("================================================================")

	testTimeSlot.ToggleCFlagSubjectID()
	testTimeSlot.ToggleCFlagInstructorID()
	testTimeSlot.ToggleCFlagRoomID()

	if universitySchedules.Get(0).Get(0).Get(1).IsConstInstructorID() != false {
		t.Errorf("uniSectionSchedules[0][0][1].IsFixedInstructorID() - Not Fixed")
	}

	if universitySchedules.Get(0).Get(0).Get(1).IsConstRoomID() != false {
		t.Error("uniSectionSchedules[0][0][1].IsFixedRoomID() - Not Fixed")
	}

	if universitySchedules.Get(0).Get(0).Get(1).IsConstSubjectID() != false {
		t.Error("uniSectionSchedules[0][0][1].IsFixedSubjectID() - Not Fixed")
	}

	if universitySchedules.Get(0).Get(0).Get(1).GetInstructorID() != 23423 {
		t.Error("uniSectionSchedules[0][0][1].GetInstructorID() != 23423")
	}

	if universitySchedules.Get(0).Get(0).Get(1).GetRoomID() != 12122 {
		t.Error("uniSectionSchedules[0][0][1].GetRoomID() != 12122")
	}

	if universitySchedules.Get(0).Get(0).Get(1).GetSubjectID() != 9898 {
		t.Error("uniSectionSchedules[0][0][1].GetSubjectID() != 9898")
	}

	if errJson != nil {
		t.Error(errJson)
	}
}
