package Rooms_test

import (
	"math/rand"
	"testing"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Schedule"
)

func Test_Rooms_timeSlotClassCount_methods(t *testing.T) {
	room := &Rooms.Room{}
	array := [Schedule.N_WEEKLY_SCHOOL_DAYS][Schedule.N_DAILY_TIME_SLOTS]uint8{}

	for day := 0; day < Schedule.N_WEEKLY_SCHOOL_DAYS; day++ {
		for time_slot := 0; time_slot < Schedule.N_DAILY_TIME_SLOTS; time_slot++ {
			array[day][time_slot] = uint8(rand.Intn(16))
			room.SetTimeSlotClassCount(day, time_slot, array[day][time_slot])
		}
	}

	for day := 0; day < Schedule.N_WEEKLY_SCHOOL_DAYS; day++ {
		for time_slot := 0; time_slot < Schedule.N_DAILY_TIME_SLOTS; time_slot++ {
			if room.GetTimeSlotClassCount(day, time_slot) != array[day][time_slot] {
				t.Errorf("set then get method result to wrong value at : [%d][%d]", day, time_slot)
			}
		}
	}

	for day := 0; day < Schedule.N_WEEKLY_SCHOOL_DAYS; day++ {
		for time_slot := 0; time_slot < Schedule.N_DAILY_TIME_SLOTS; time_slot++ {
			num_of_inc := 1 + rand.Intn(2)
			for i := 0; i < num_of_inc; i++ {
				if array[day][time_slot] < 15 {
					array[day][time_slot]++
				}

				if room.GetTimeSlotClassCount(day, time_slot) < 15 {
					room.IncTimeSlotClassCount(day, time_slot)
				}
			}
		}
	}

	for day := 0; day < Schedule.N_WEEKLY_SCHOOL_DAYS; day++ {
		for time_slot := 0; time_slot < Schedule.N_DAILY_TIME_SLOTS; time_slot++ {
			if room.GetTimeSlotClassCount(day, time_slot) != array[day][time_slot] {
				t.Errorf("increment result to wrong value at : [%d][%d]", day, time_slot)
			}
		}
	}
}
