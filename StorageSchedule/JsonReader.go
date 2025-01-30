package StorageSchedule

import (
	"path"

	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

type JsonReader struct{}

const first_semester int = 0
const second_semester int = 1

func (s *JsonReader) LoadSchedules(semester int) (Schedule.UniTimeTables, error) {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	var saved_file string

	if semester == first_semester {
		saved_file = path.Join(project_root, "scheduling-system-temporary-data", "univ-1st-sem.sched")
	} else if semester == second_semester {
		saved_file = path.Join(project_root, "scheduling-system-temporary-data", "univ-2nd-sem.sched")
	}

	read_bytes, err_read_from_bin_file := Utils.ReadFromBinFile(saved_file)

	if err_read_from_bin_file != nil {
		return nil, err_read_from_bin_file
	}

	deserialize_data := Schedule.DeserializeUniversitySchedule(read_bytes)

	return deserialize_data, nil
}
