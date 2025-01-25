package Storage

import (
	"encoding/json"
	"os"
	"path"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

// this type is for development and testing only
type JsonFilePersistence struct{}

func (s *JsonFilePersistence) GetAllRooms() ([]Rooms.Room, error) {

	projectRoot, err_find_root := Utils.FindProjectRoot()
	if err_find_root != nil {
		return nil, err_find_root
	}

	rooms_json_file := path.Join(projectRoot, "scheduling-system-temporary-data", "rooms.json")
	rooms_byte_data, err := os.ReadFile(rooms_json_file)

	if err != nil {
		return nil, err
	}

	all_rooms := make([]Rooms.Room, 0)
	err = json.Unmarshal(rooms_byte_data, &all_rooms)

	if err != nil {
		return nil, err
	}

	for idx := range all_rooms {
		all_rooms[idx].RoomID = uint16(idx + 1)
	}

	return all_rooms, nil
}

func (s *JsonFilePersistence) GetAllInstructors() ([]Instructors.Instructor, error) {

	projectRoot, err_find_root := Utils.FindProjectRoot()
	if err_find_root != nil {
		return nil, err_find_root
	}

	instructors_json_file := path.Join(projectRoot, "scheduling-system-temporary-data", "instructors.json")
	instructors_byte_data, err := os.ReadFile(instructors_json_file)

	if err != nil {
		return nil, err
	}

	all_instructors := make([]Instructors.Instructor, 0)
	err = json.Unmarshal(instructors_byte_data, &all_instructors)

	if err != nil {
		return nil, err
	}

	for idx := range all_instructors {
		all_instructors[idx].InstructorID = uint16(idx + 1)
	}

	return all_instructors, nil
}

func (s *JsonFilePersistence) GetAllCurriculum() ([]Curriculum.Curriculum, error) {

	// map for getting subject id using subject code
	subject_code_map_id := make(map[string]uint16)

	subject_data, err := s.GetAllSubjects()

	if err != nil {
		return nil, err
	}

	// get the ID of each subject code

	for i, subject := range subject_data {
		subject_code_map_id[subject.Code] = uint16(i + 1)
	}

	// read all curriculum json files

	projectRoot, err_find_root := Utils.FindProjectRoot()
	if err_find_root != nil {
		return nil, err_find_root
	}

	curriculum_json_folder := path.Join(projectRoot, "scheduling-system-temporary-data", "curriculums")
	curriculum_json_files, curriculum_read_err := os.ReadDir(curriculum_json_folder)

	if curriculum_read_err != nil {
		return nil, curriculum_read_err
	}

	all_curriculums := make([]Curriculum.Curriculum, 0)

	for i, json_file := range curriculum_json_files {
		curriculum_data := &Curriculum.Curriculum{}
		curriculum_data.CurriculumID = uint16(i + 1)

		json_file_path := path.Join(curriculum_json_folder, json_file.Name())
		byte_data, err := os.ReadFile(json_file_path)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byte_data, &curriculum_data)
		if err != nil {
			return nil, err
		}

		for i, year_level := range curriculum_data.YearLevels {
			for j, semester := range year_level.Semesters {
				for k, subject := range semester.Subjects {
					curriculum_data.YearLevels[i].Semesters[j].Subjects[k].ID = subject_code_map_id[subject.Code]
				}
			}
		}

		all_curriculums = append(all_curriculums, *curriculum_data)
	}

	return all_curriculums, nil
}

func (s *JsonFilePersistence) GetAllSubjects() ([]Curriculum.Subject, error) {

	projectRoot, err_find_root := Utils.FindProjectRoot()
	if err_find_root != nil {
		return nil, err_find_root
	}

	subjects_json_file := path.Join(projectRoot, "scheduling-system-temporary-data", "all-subjects.json")
	subjects_byte_data, err := os.ReadFile(subjects_json_file)

	if err != nil {
		return nil, err
	}

	subject_data := []Curriculum.Subject{}
	err = json.Unmarshal(subjects_byte_data, &subject_data)

	if err != nil {
		return nil, err
	}

	for i := range subject_data {
		subject_data[i].ID = uint16(i + 1)
	}

	return subject_data, nil
}

func (s *JsonFilePersistence) GetAllDepartments() ([]Departments.Department, error) {

	projectRoot, err_find_root := Utils.FindProjectRoot()
	if err_find_root != nil {
		return nil, err_find_root
	}

	departments_json_file := path.Join(projectRoot, "scheduling-system-temporary-data", "departments.json")
	departments_byte_data, err := os.ReadFile(departments_json_file)

	if err != nil {
		return nil, err
	}

	all_departments := make([]Departments.Department, 0)
	err = json.Unmarshal(departments_byte_data, &all_departments)

	if err != nil {
		return nil, err
	}

	return all_departments, nil
}
