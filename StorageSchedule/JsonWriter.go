package StorageSchedule

import (
	"path"
	"sync"

	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

type JsonWriter struct {
	Mutex *sync.Mutex
}

func (s *JsonWriter) SaveSchedules(university_schedule Schedule.UniTimeTables, semester int) error {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	var saved_file string

	if semester == first_semester {
		saved_file = path.Join(project_root, "scheduling-system-temporary-data", "univ-1st-sem.sched")
	} else if semester == second_semester {
		saved_file = path.Join(project_root, "scheduling-system-temporary-data", "univ-2nd-sem.sched")
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	Utils.SaveToBinFile(saved_file, Schedule.SerializeUniversitySchedule(university_schedule))

	return nil
}
