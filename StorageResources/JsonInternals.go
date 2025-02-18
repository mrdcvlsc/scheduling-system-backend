package StorageResources

import (
	"encoding/json"
	"os"
	"path"
	"sort"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

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

	for idx := range instructors {
		instructors[idx].InstructorID = uint16(idx + 1)
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
