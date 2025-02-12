package StorageSchedule

import "github.com/mrdcvlsc/scheduling-system-backend/Schedule"

type saveRepository interface {
	LoadSchedules(semester int) (Schedule.UniTimeTables, error)
}

type loadRepository interface {
	SaveSchedules(university_schedule Schedule.UniTimeTables, semester int) error
}

type Persistence struct {
	LoadService loadRepository
	SaveService saveRepository
}
