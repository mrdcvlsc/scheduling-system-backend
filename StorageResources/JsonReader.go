package StorageResources

import (
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

// this type is for development and testing only
type JsonReader struct{}

func (s *JsonReader) ReadAllRooms() ([]Rooms.Room, error) {
	return json_read_all_rooms()
}

func (s *JsonReader) ReadAllInstructors() ([]Instructors.Instructor, error) {

	instructors_with_time_string, err := s.ReadAllInstructorsWithTimeString()

	if err != nil {
		return nil, err
	}

	instructors := make([]Instructors.Instructor, 0)

	for _, instructor_with_time_str := range instructors_with_time_string {
		instructor := Instructors.Instructor{
			InstructorID:  instructor_with_time_str.InstructorID,
			DepartmentID:  instructor_with_time_str.DepartmentID,
			FirstName:     instructor_with_time_str.FirstName,
			MiddleInitial: instructor_with_time_str.MiddleInitial,
			LastName:      instructor_with_time_str.LastName,
		}

		instructor.Time.StringParse(instructor_with_time_str.Time)

		instructors = append(instructors, instructor)
	}

	return instructors, nil
}

func (s *JsonReader) ReadDepartmentInstructors(department_id int) ([]Instructors.Instructor, error) {

	instructors_with_time_string, err := s.ReadAllInstructorsWithTimeString()

	if err != nil {
		return nil, err
	}

	department_instructors := make([]Instructors.Instructor, 0)

	for _, instructor_with_time_str := range instructors_with_time_string {
		if instructor_with_time_str.DepartmentID != uint16(department_id) {
			continue
		}

		instructor := Instructors.Instructor{
			InstructorID:  instructor_with_time_str.InstructorID,
			DepartmentID:  instructor_with_time_str.DepartmentID,
			FirstName:     instructor_with_time_str.FirstName,
			MiddleInitial: instructor_with_time_str.MiddleInitial,
			LastName:      instructor_with_time_str.LastName,
		}

		instructor.Time.StringParse(instructor_with_time_str.Time)

		department_instructors = append(department_instructors, instructor)
	}

	return department_instructors, nil
}

func (s *JsonReader) ReadAllCurriculum() ([]Curriculum.Curriculum, error) {
	return json_read_all_curriculums()
}

func (s *JsonReader) ReadAllSubjects() ([]Curriculum.Subject, error) {
	return json_read_all_subjects()
}

func (s *JsonReader) ReadAllDepartments() ([]Departments.Department, error) {
	return json_read_all_departments()
}

func (s *JsonReader) ReadAllInstructorsWithTimeString() ([]Instructors.InstructorWithTimeString, error) {
	return json_read_all_instructors_with_time_string()
}
