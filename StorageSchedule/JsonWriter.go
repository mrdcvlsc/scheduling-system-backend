package StorageSchedule

import (
	"fmt"
	"path"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

type JsonWriter struct{}

func (s *JsonWriter) SaveSchedules(university_schedule Schedule.UniTimeTables, semester int) error {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	var saved_file string

	for semester_idx := range len(Curriculum.SEMESTER_INDEX_NAME) {
		if semester_idx == semester {
			saved_file = path.Join(project_root, "scheduling-system-temporary-data", fmt.Sprintf("univ-sem-%d.sched", semester_idx+1))
			break
		}
	}

	UniSchedPersistenceMutex.Lock()
	defer UniSchedPersistenceMutex.Unlock()

	Utils.SaveToBinFile(saved_file, Schedule.SerializeUniversitySchedule(university_schedule))

	return nil
}
