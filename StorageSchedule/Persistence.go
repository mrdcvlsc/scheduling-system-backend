package StorageSchedule

import "github.com/mrdcvlsc/scheduling-system-backend/Schedule"

type readerRepository interface {
	LoadSchedules(semester int) (Schedule.UniTimeTables, error)
}

type writerRepository interface {
	SaveSchedules(university_schedule Schedule.UniTimeTables, semester int) error
}

type Persistence struct {
	ReaderService readerRepository
	WriterService writerRepository
}
