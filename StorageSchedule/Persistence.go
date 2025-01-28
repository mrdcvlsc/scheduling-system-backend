package StorageSchedule

import "github.com/mrdcvlsc/scheduling-system-backend/Schedule"

type writerRepository interface {
	SaveSchedules(university_schedule Schedule.UniTimeTables) error
}

type Persistence struct {
	WriterService writerRepository
}
