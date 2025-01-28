package StorageReader

import "github.com/mrdcvlsc/scheduling-system-backend/Resources/Schedule"

type writerRepository interface {
	SaveSchedules(university_schedule *Schedule.WeekTimeTable) error
}

type Persistence struct {
	Service writerRepository
}
