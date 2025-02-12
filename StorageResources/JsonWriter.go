package StorageResources

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"sort"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

// this type is for development and testing only
type JsonWriter struct{}

func (s *JsonWriter) CreateInstructor(instructor Instructors.Instructor) error {
	if instructor.InstructorID != 0 {
		return errors.New("cannot create a new instructor with a non zero instructor ID")
	}

	instructor_with_time_str, err_read := s.readAllInstructorsWithTimeString()

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

	s.saveAllInstructorsWithTimeString(instructor_with_time_str)

	return nil
}

func (s *JsonWriter) UpdateInstructor(instructor Instructors.Instructor) error {
	if instructor.InstructorID == 0 {
		return errors.New("parameter argument missing invalid instructor ID")
	}

	instructors_with_time_string, err_read := s.readAllInstructorsWithTimeString()

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

	save_err := s.saveAllInstructorsWithTimeString(instructors_with_time_string)

	if save_err != nil {
		return save_err
	}

	return nil
}

func (s *JsonWriter) readAllInstructorsWithTimeString() ([]Instructors.InstructorWithTimeString, error) {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	instructors_json_file := path.Join(project_root, "scheduling-system-temporary-data", "instructors.json")
	instructors_byte_data, err := os.ReadFile(instructors_json_file)

	if err != nil {
		return nil, err
	}

	instructors := make([]Instructors.InstructorWithTimeString, 0)
	err = json.Unmarshal(instructors_byte_data, &instructors)

	if err != nil {
		return nil, err
	}

	for idx := range instructors {
		instructors[idx].InstructorID = uint16(idx + 1)
	}

	sort.Slice(instructors, func(i, j int) bool {
		return instructors[i].InstructorID < instructors[j].InstructorID
	})

	return instructors, nil
}

func (s *JsonWriter) saveAllInstructorsWithTimeString(instructors []Instructors.InstructorWithTimeString) error {
	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	instructors_json_file := path.Join(project_root, "scheduling-system-temporary-data", "instructors.json")

	sort.Slice(instructors, func(i, j int) bool {
		return instructors[i].InstructorID < instructors[j].InstructorID
	})

	instructors_byte_data, err := json.MarshalIndent(instructors, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(instructors_json_file, instructors_byte_data, 0644); err != nil {
		return err
	}

	return nil
}
