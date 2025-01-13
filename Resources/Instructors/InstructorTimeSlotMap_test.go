package Instructors_test

import (
	"fmt"
	"testing"

	instructors "github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
)

func TestInstructorTimeSlotMap(t *testing.T) {
	mr_instructor := instructors.InstructorTimeSlotMap{}

	cnt := 0
	index := 0
	for day := 0; day < 6; day++ {
		for time_slot := 0; time_slot < 24; time_slot++ {

			available := mr_instructor.GetAvailability(day, time_slot)

			if cnt == 0 || cnt == 63 {
				fmt.Printf("Array Values (day = %d, time_slot = %d)[0]: %064b\n", day, time_slot, mr_instructor)
			}

			if available == false {
				t.Errorf(
					"GetAvailability(day = %d, time_slot = %d) : {loop test phase 1} should not be false yet",
					day, time_slot,
				)
			}

			mr_instructor.SetAvailability(false, day, time_slot)
			mr_instructor.SetAvailability(false, day, time_slot)
			available = mr_instructor.GetAvailability(day, time_slot)

			if cnt == 0 || cnt == 63 {
				fmt.Printf("Array Values (day = %d, time_slot = %d)[1]: %064b\n", day, time_slot, mr_instructor)
			}

			if available == true {
				t.Errorf(
					"GetAvailability(day = %d, time_slot = %d) : {loop test phase 2} should be false now",
					day, time_slot,
				)
			}

			if mr_instructor[index] != (uint64(1) << cnt) {
				t.Errorf(
					"[day:%d, time_slot:%d] : {loop test phase 3} (bitsetmap[%d] = %d) != ((uint64(1) << cnt) = %d)",
					day, time_slot, index, mr_instructor[index], (uint64(1) << cnt),
				)
			}

			mr_instructor.SetAvailability(true, day, time_slot)
			mr_instructor.SetAvailability(true, day, time_slot)
			available = mr_instructor.GetAvailability(day, time_slot)

			if available == false {
				t.Errorf(
					"GetAvailability(day = %d, time_slot = %d) : {loop test phase 4} should not be false again",
					day, time_slot,
				)
			}

			if mr_instructor[index] != uint64(0) {
				t.Errorf(
					"[day:%d, time_slot:%d] : {loop test phase 5} (bitsetmap[%d] = %d) != uint64(0))",
					day, time_slot, index, mr_instructor[index],
				)
			}

			if cnt == 0 || cnt == 63 {
				fmt.Printf("Array Values (day = %d, time_slot = %d)[2]: %064b\n\n", day, time_slot, mr_instructor)
			}

			cnt++

			if cnt == 64 {
				cnt = 0
				index++
			}
		}
	}
}
