package StorageResources

import "sync"

// this type is for development and testing only
type JsonReader struct {
	SubjectMutex     *sync.Mutex
	CurriculumsMutex *sync.Mutex
	DepartmentMutex  *sync.Mutex
	InstructorMutex  *sync.Mutex
	RoomMutex        *sync.Mutex
}
