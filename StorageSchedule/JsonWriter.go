package StorageSchedule

import (
	"path"

	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

type JsonWriter struct{}

func (s *JsonWriter) SaveSchedules(university_schedule Schedule.UniTimeTables) error {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	saved_file := path.Join(project_root, "scheduling-system-temporary-data", "univ.sched")

	Utils.SaveToBinFile(saved_file, Schedule.SerializeUniversitySchedule(university_schedule))

	return nil
}
