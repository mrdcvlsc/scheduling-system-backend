package StorageResources

import (
	"errors"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
)

// this type is for development and testing only
type JsonWriter struct{}

func (s *JsonWriter) CreateInstructor(instructor Instructors.Instructor) error {
	if instructor.InstructorID != 0 {
		return errors.New("cannot create a new instructor with a non zero instructor ID")
	}

	instructor_with_time_str, err_read := json_read_all_instructors_with_time_string()

	if err_read != nil {
		return err_read
	}

	if instructor.InstructorID != 0 {
		for _, instructor_with_time_str := range instructor_with_time_str {
			if instructor_with_time_str.DepartmentID == instructor.InstructorID {
				return errors.New("that instructor already exist")
			}
		}
	}

	instructor_with_time_str = append(instructor_with_time_str, Instructors.InstructorWithTimeString{
		InstructorID:  uint16(len(instructor_with_time_str) + 1),
		DepartmentID:  instructor.DepartmentID,
		FirstName:     instructor.FirstName,
		MiddleInitial: instructor.MiddleInitial,
		LastName:      instructor.LastName,
		Time:          instructor.Time.Stringify(),
	})

	json_save_all_instructors_with_time_string(instructor_with_time_str)

	return nil
}

func (s *JsonWriter) UpdateInstructor(instructor Instructors.Instructor) error {
	if instructor.InstructorID == 0 {
		return errors.New("parameter argument missing invalid instructor ID")
	}

	instructors_with_time_string, err_read := json_read_all_instructors_with_time_string()

	if err_read != nil {
		return err_read
	}

	has_id := false
	to_update_idx := -1

	for idx, instructor_w_t_str := range instructors_with_time_string {
		if instructor_w_t_str.InstructorID == instructor.InstructorID {
			has_id = true
			to_update_idx = idx
			break
		}
	}

	if !has_id {
		return errors.New("instructor to update does not exist in the json file")
	}

	instructors_with_time_string[to_update_idx] = Instructors.InstructorWithTimeString{
		InstructorID:  instructor.InstructorID,
		DepartmentID:  instructor.DepartmentID,
		FirstName:     instructor.FirstName,
		MiddleInitial: instructor.MiddleInitial,
		LastName:      instructor.LastName,
		Time:          instructor.Time.Stringify(),
	}

	err_save_instructors := json_save_all_instructors_with_time_string(instructors_with_time_string)

	if err_save_instructors != nil {
		return err_save_instructors
	}

	return nil
}
