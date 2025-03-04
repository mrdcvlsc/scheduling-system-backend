package StorageResources

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func json_read_all_departments() ([]Departments.Department, error) {

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

func json_save_all_departments(departments []Departments.Department) error {
	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	departments_json_file := path.Join(project_root, "scheduling-system-temporary-data", "departments.json")

	sort.Slice(departments, func(i, j int) bool {
		return departments[i].DepartmentID < departments[j].DepartmentID
	})

	departments_byte_data, err := json.MarshalIndent(departments, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(departments_json_file, departments_byte_data, 0644); err != nil {
		return err
	}

	return nil
}

func json_read_all_curriculums() ([]Curriculum.Curriculum, error) {

	// map for getting subject id using subject code
	// subject_code_to_subject_id := make(map[string]uint16)

	// subjects, err_subjects := json_read_all_subjects()

	// if err_subjects != nil {
	// 	return nil, err_subjects
	// }

	// get the ID of each subject code

	// for subject_idx, subject := range subjects {
	// 	subject_code_to_subject_id[subject.Code] = uint16(subject_idx + 1)
	// }

	// read all curriculum json files

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return nil, err_project_root
	}

	curriculum_json_folder := path.Join(project_root, "scheduling-system-temporary-data", "curriculums")
	curriculum_json_files, err_read_dir := os.ReadDir(curriculum_json_folder)

	if err_read_dir != nil {
		return nil, err_read_dir
	}

	curriculums := make([]Curriculum.Curriculum, 0)

	for _, json_file := range curriculum_json_files {
		course := &Curriculum.Curriculum{}
		json_file_path := path.Join(curriculum_json_folder, json_file.Name())

		byte_data, err_byte_data := os.ReadFile(json_file_path)

		if err_byte_data != nil {
			return nil, err_byte_data
		}

		if err := json.Unmarshal(byte_data, &course); err != nil {
			return nil, err
		}

		// for year_level_idx, year_level := range course.YearLevels {
		// 	for semester_idx, semester := range year_level.Semesters {
		// 		for subject_idx, subject := range semester.Subjects {

		// 			_, has_subject_code := subject_code_to_subject_id[subject.Code]

		// 			if !has_subject_code {
		// 				return nil, fmt.Errorf("the subject code `%s` was not found in the subjects json file", subject.Code)
		// 			}

		// 			course.YearLevels[year_level_idx].Semesters[semester_idx].Subjects[subject_idx].ID = subject_code_to_subject_id[subject.Code]
		// 		}
		// 	}
		// }

		curriculums = append(curriculums, *course)
	}

	sort.Slice(curriculums, func(i, j int) bool {
		return curriculums[i].CurriculumID < curriculums[j].CurriculumID
	})

	return curriculums, nil
}

// no op if curriculum slice is empty
func json_save_all_curriculums(curriculums []Curriculum.Curriculum) error {
	for _, curriculum := range curriculums {
		err_save := json_save_curriculum(
			fmt.Sprintf("%s.json", Utils.RemoveWhiteSpace(curriculum.CurriculumCode)), curriculum,
		)

		if err_save != nil {
			return err_save
		}
	}

	return nil
}

func json_save_curriculum(filename string, curriculum Curriculum.Curriculum) error {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	new_curriculum_json_file := path.Join(project_root, "scheduling-system-temporary-data", "curriculums", filename)

	curriculum_byte_data, err := json.MarshalIndent(curriculum, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(new_curriculum_json_file, curriculum_byte_data, 0644); err != nil {
		return err
	}

	return nil
}

func json_edit_curriculum(filename_old, filename_new string, curriculum_new Curriculum.Curriculum) error {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	old_curriculum_json_file := path.Join(project_root, "scheduling-system-temporary-data", "curriculums", filename_old)
	new_curriculum_json_file := path.Join(project_root, "scheduling-system-temporary-data", "curriculums", filename_new)

	curriculum_byte_data, err_marshal_indent := json.MarshalIndent(curriculum_new, "", "  ")

	if err_marshal_indent != nil {
		return err_marshal_indent
	}

	if err_write_file := os.WriteFile(new_curriculum_json_file, curriculum_byte_data, 0644); err_write_file != nil {
		return err_write_file
	}

	if err_remove := os.Remove(old_curriculum_json_file); err_remove != nil {
		return err_remove
	}

	return nil
}

func json_read_all_subjects() ([]Curriculum.Subject, error) {

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

	sort.Slice(subjects, func(i, j int) bool {
		return subjects[i].ID < subjects[j].ID
	})

	return subjects, nil
}

func json_save_all_subjects(subjects []Curriculum.Subject) error {

	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	subjects_json_file := path.Join(project_root, "scheduling-system-temporary-data", "all-subjects.json")

	sort.Slice(subjects, func(i, j int) bool {
		return subjects[i].ID < subjects[j].ID
	})

	rooms_byte_data, err := json.MarshalIndent(subjects, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(subjects_json_file, rooms_byte_data, 0644); err != nil {
		return err
	}

	return nil
}

func json_read_all_rooms() ([]Rooms.Room, error) {

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

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].RoomID < rooms[j].RoomID
	})

	return rooms, nil
}

func json_save_all_rooms(rooms []Rooms.Room) error {
	project_root, err_project_root := Utils.FindProjectRoot()

	if err_project_root != nil {
		return err_project_root
	}

	rooms_json_file := path.Join(project_root, "scheduling-system-temporary-data", "rooms.json")

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].RoomID < rooms[j].RoomID
	})

	rooms_byte_data, err := json.MarshalIndent(rooms, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(rooms_json_file, rooms_byte_data, 0644); err != nil {
		return err
	}

	return nil
}

func json_read_all_instructors_with_time_string() ([]Instructors.InstructorWithTimeString, error) {

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

	sort.Slice(instructors, func(i, j int) bool {
		return instructors[i].InstructorID < instructors[j].InstructorID
	})

	return instructors, nil
}

func json_save_all_instructors_with_time_string(instructors []Instructors.InstructorWithTimeString) error {
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
