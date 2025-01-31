package RouteGlobals

import (
	"errors"
	"sync"

	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageSchedule"
)

var ResourcesPersistence *StorageResources.Persistence
var SchedulePersistence *StorageSchedule.Persistence
var ScheduleCache *scheduleCache

const NUM_OF_SEMESTERS int = 2

type scheduleCache struct {
	rw_mutex          sync.RWMutex
	semester_schedule [NUM_OF_SEMESTERS]Schedule.UniTimeTables
}

// returns: (Schedule.UniTimeTables, nil) if has cached schedule, (nil, error) if something goes wrong (nil, nil) if no cached schedule.
func (s *scheduleCache) GetCachedUniversitySchedule(semester int) (Schedule.UniTimeTables, bool, error) {

	if semester < 0 {
		return nil, false, errors.New("cached schedule semester index underflow")
	}

	if semester >= NUM_OF_SEMESTERS {
		return nil, false, errors.New("cached schedule semester index overflow")
	}

	s.rw_mutex.Lock()
	defer s.rw_mutex.Unlock()

	if s.semester_schedule[semester] != nil {
		return s.semester_schedule[semester], true, nil
	}

	return nil, false, nil
}

// returns: (Schedule.UniTimeTables, nil) if has cached schedule, (nil, error) if something goes wrong (nil, nil) if no cached schedule.
func (s *scheduleCache) SetCachedUniversitySchedule(semester int, university_schedule Schedule.UniTimeTables) error {

	if semester < 0 {
		return errors.New("cached schedule semester index underflow")
	}

	if semester >= NUM_OF_SEMESTERS {
		return errors.New("cached schedule semester index overflow")
	}

	s.rw_mutex.Lock()
	defer s.rw_mutex.Unlock()

	s.semester_schedule[semester] = university_schedule

	return nil
}

var Test int
