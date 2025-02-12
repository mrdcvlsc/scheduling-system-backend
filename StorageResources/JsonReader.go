package StorageResources

import (
	"encoding/json"
	"os"
	"path"
	"sort"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

// this type is for development and testing only
type JsonReader struct{}

func (s *JsonReader) ReadAllRooms() ([]Rooms.Room, error) {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	rooms_json_file := path.Join(project_root, "scheduling-system-temporary-data", "rooms.json")
	rooms_byte_data, err := os.ReadFile(rooms_json_file)

	if err != nil {
		return nil, err
	}

	rooms := make([]Rooms.Room, 0)
	err = json.Unmarshal(rooms_byte_data, &rooms)

	if err != nil {
		return nil, err
	}

	for room_idx := range rooms {
		rooms[room_idx].RoomID = uint16(room_idx + 1)
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].RoomID < rooms[j].RoomID
	})

	return rooms, nil
}

func (s *JsonReader) ReadAllInstructors() ([]Instructors.Instructor, error) {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	instructors_json_file := path.Join(project_root, "scheduling-system-temporary-data", "instructors.json")
	instructors_byte_data, err := os.ReadFile(instructors_json_file)

	if err != nil {
		return nil, err
	}

	instructors := make([]Instructors.Instructor, 0)
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

func (s *JsonReader) ReadDepartmentInstructors(department_id int) ([]Instructors.Instructor, error) {
	all_instructors, err := s.ReadAllInstructors()

	if err != nil {
		return nil, err
	}

	department_instructors := make([]Instructors.Instructor, 0)

	for _, instructor := range all_instructors {
		if instructor.DepartmentID == uint16(department_id) {
			department_instructors = append(department_instructors, instructor)
		}
	}

	return department_instructors, nil
}

func (s *JsonReader) ReadAllCurriculum() ([]Curriculum.Curriculum, error) {

	// map for getting subject id using subject code
	subject_code_to_subject_id := make(map[string]uint16)

	subjects, err_subjects := s.ReadAllSubjects()

	if err_subjects != nil {
		return nil, err_subjects
	}

	// get the ID of each subject code

	for subject_idx, subject := range subjects {
		subject_code_to_subject_id[subject.Code] = uint16(subject_idx + 1)
	}

	// read all curriculum json files

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	curriculum_json_folder := path.Join(project_root, "scheduling-system-temporary-data", "curriculums")
	curriculum_json_files, curriculum_read_err := os.ReadDir(curriculum_json_folder)

	if curriculum_read_err != nil {
		return nil, curriculum_read_err
	}

	curriculums := make([]Curriculum.Curriculum, 0)

	for i, json_file := range curriculum_json_files {
		course := &Curriculum.Curriculum{}
		course.CurriculumID = uint16(i + 1)

		json_file_path := path.Join(curriculum_json_folder, json_file.Name())

		byte_data, err_byte_data := os.ReadFile(json_file_path)

		if err_byte_data != nil {
			return nil, err_byte_data
		}

		if err := json.Unmarshal(byte_data, &course); err != nil {
			return nil, err
		}

		for year_level_idx, year_level := range course.YearLevels {
			for semester_idx, semester := range year_level.Semesters {
				for subject_idx, subject := range semester.Subjects {
					course.YearLevels[year_level_idx].Semesters[semester_idx].Subjects[subject_idx].ID = subject_code_to_subject_id[subject.Code]
				}
			}
		}

		curriculums = append(curriculums, *course)
	}

	sort.Slice(curriculums, func(i, j int) bool {
		return curriculums[i].CurriculumID < curriculums[j].CurriculumID
	})

	return curriculums, nil
}

func (s *JsonReader) ReadAllSubjects() ([]Curriculum.Subject, error) {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	subjects_json_file := path.Join(project_root, "scheduling-system-temporary-data", "all-subjects.json")
	subjects_byte_data, err := os.ReadFile(subjects_json_file)

	if err != nil {
		return nil, err
	}

	subjects := []Curriculum.Subject{}

	err = json.Unmarshal(subjects_byte_data, &subjects)

	if err != nil {
		return nil, err
	}

	for i := range subjects {
		subjects[i].ID = uint16(i + 1)
	}

	sort.Slice(subjects, func(i, j int) bool {
		return subjects[i].ID < subjects[j].ID
	})

	return subjects, nil
}

func (s *JsonReader) ReadAllDepartments() ([]Departments.Department, error) {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	departments_json_file := path.Join(project_root, "scheduling-system-temporary-data", "departments.json")
	departments_byte_data, err := os.ReadFile(departments_json_file)

	if err != nil {
		return nil, err
	}

	departments := make([]Departments.Department, 0)

	if err = json.Unmarshal(departments_byte_data, &departments); err != nil {
		return nil, err
	}

	sort.Slice(departments, func(i, j int) bool {
		return departments[i].DepartmentID < departments[j].DepartmentID
	})

	return departments, nil
}
