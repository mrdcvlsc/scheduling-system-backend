package Schedule_test

import (
	"os"
	"testing"

	geneticalgorithm "github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func Test_UniTimeTablesSerializationAndDeserialization(t *testing.T) {

	{
		uni_sched, err := geneticalgorithm.NewIndividual(0, 0)

		if uni_sched == nil && err != nil {
			t.Fatalf("there was an error reading data: %v\n", err)
		}

		if uni_sched.IsEmpty() {
			t.Fatal("there was no schedule generated to be tested")
		}

		serialized_data := Schedule.SerializeUniversitySchedule(&uni_sched)

		if err := Utils.SaveToBinFile("tmp-university.schedule", serialized_data); err != nil {
			t.Fatalf("failed to save serialized data to binary file: %v\n", err)
		}
	}

	{
		read_bytes, err := Utils.ReadFromBinFile("tmp-university.schedule")

		if err != nil {
			t.Fatalf("failed to read the serialized data from the binary file: %v\n", err)
		}

		deserialize_data := Schedule.DeserializeUniversitySchedule(read_bytes)

		for section_idx, original := range *deserialize_data {
			for day := 0; day < Const.N_WEEKLY_SCHOOL_DAYS; day++ {
				for time_slot := 0; time_slot < Const.N_DAILY_TIME_SLOTS; time_slot++ {
					if (*deserialize_data)[section_idx][day][time_slot] != original[day][time_slot] {
						t.Errorf("Data missmatch at time slot (day: %d, time_slot:%d)\n", day, time_slot)
					}
				}
			}
		}
	}

	err := os.Remove("tmp-university.schedule")

	if err != nil {
		t.Fatalf("error deleting file: %v\n", err)
	}
}
