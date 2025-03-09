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

	curriculums_json_file := path.Join(project_root, "scheduling-system-temporary-data", "curriculums.json")
	curriculums_byte_data, err := os.ReadFile(curriculums_json_file)

	if err != nil {
		return nil, err
	}

	curriculums := make([]Curriculum.Curriculum, 0)
	err = json.Unmarshal(curriculums_byte_data, &curriculums)

	if err != nil {
		return nil, err
	}

	sort.Slice(curriculums, func(i, j int) bool {
		return curriculums[i].CurriculumID < curriculums[j].CurriculumID
	})

	return curriculums, nil
}
