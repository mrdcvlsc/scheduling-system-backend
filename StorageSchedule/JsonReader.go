package StorageSchedule

import (
	"path"

	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

type JsonReader struct{}

func (s *JsonReader) LoadSchedules() (Schedule.UniTimeTables, error) {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	saved_file := path.Join(project_root, "scheduling-system-temporary-data", "univ.sched")

	read_bytes, err_read_from_bin_file := Utils.ReadFromBinFile(saved_file)

	if err_read_from_bin_file != nil {
		return nil, err_read_from_bin_file
	}

	deserialize_data := Schedule.DeserializeUniversitySchedule(read_bytes)

	return deserialize_data, nil
}
