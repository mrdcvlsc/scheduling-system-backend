package StorageResources

import (
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

type readerRepository interface {
	// the order of curriculums returned by this method is always sorted by curriculum ID.
	ReadAllCurriculum() ([]Curriculum.Curriculum, error)
	ReadAllSubjects() ([]Curriculum.Subject, error)
	ReadAllDepartments() ([]Departments.Department, error)
	ReadAllInstructors() ([]Instructors.Instructor, error)
	ReadDepartmentInstructors(department_id int) ([]Instructors.Instructor, error)
	ReadAllRooms() ([]Rooms.Room, error)
}

type writerRepository interface {
	CreateDepartment(new_department Departments.Department) error
	UpdateDepartment(department_to_update Departments.Department) error

	CreateSubject(new_subject Curriculum.Subject) error
	UpdateSubject(subject_to_update Curriculum.Subject) error

	CreateCurriculum(new_curriculum Curriculum.Curriculum) error
	UpdateCurriculum(curriculum_old_name string, curriculum_new Curriculum.Curriculum) error

	CreateInstructor(new_instructor Instructors.Instructor) error
	UpdateInstructor(instructor_to_update Instructors.Instructor) error
	DeleteInstructor(instructor_id uint16) error

	CreateRoom(new_room Rooms.Room) error
	UpdateRoom(room_to_update Rooms.Room) error
	DeleteRoom(room_id uint16) error
}

type Persistence struct {
	ReaderService readerRepository
	WriterService writerRepository
}
