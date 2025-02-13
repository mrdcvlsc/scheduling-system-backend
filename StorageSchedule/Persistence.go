package StorageSchedule

import "github.com/mrdcvlsc/scheduling-system-backend/Schedule"

type loadRepository interface {
	LoadSchedules(semester int) (Schedule.UniTimeTables, error)
}

type saveRepository interface {
	SaveSchedules(university_schedule Schedule.UniTimeTables, semester int) error
}

type Persistence struct {
	SaveService saveRepository
	LoadService loadRepository
}
