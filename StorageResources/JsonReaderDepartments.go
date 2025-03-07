package StorageResources

import (
	"encoding/json"
	"os"
	"path"
	"sort"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func (s *JsonReader) ReadAllDepartments() ([]Departments.Department, error) {
	return json_read_all_departments()
}

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
