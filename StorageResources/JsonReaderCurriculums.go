package StorageResources

import (
	"encoding/json"
	"os"
	"path"
	"sort"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func (s *JsonReader) ReadAllCurriculum() ([]Curriculum.Curriculum, error) {
	return json_read_all_curriculums()
}

func json_read_all_curriculums() ([]Curriculum.Curriculum, error) {
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

		curriculums = append(curriculums, *course)
	}

	sort.Slice(curriculums, func(i, j int) bool {
		return curriculums[i].CurriculumID < curriculums[j].CurriculumID
	})

	return curriculums, nil
}
