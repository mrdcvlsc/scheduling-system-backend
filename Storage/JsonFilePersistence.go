package storage

import (
	"encoding/json"
	"os"
	"path"

	c "github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
)

type JsonFilePersistence struct{}

func (s *JsonFilePersistence) GetAllCurriculum() []c.Curriculum {

	// map for getting subject id using subject code
	subject_code_map_id := make(map[string]uint16)

	// read all subject from all-subjects.json file

	subjects_json_file := path.Join("..", "scheduling-system-temporary-data", "all-subjects.json")
	subjects_byte_data, err := os.ReadFile(subjects_json_file)

	if err != nil {
		panic(err)
	}

	subject_data := []c.Subject{}
	err = json.Unmarshal(subjects_byte_data, &subject_data)

	// get the ID of each subject code

	for i, subject := range subject_data {
		subject_code_map_id[subject.Code] = uint16(i + 1)
	}

	if err != nil {
		panic(err)
	}

	// read all curriculum json files

	curriculum_json_folder := path.Join("..", "scheduling-system-temporary-data", "curriculums")
	curriculum_json_files, curriculum_read_err := os.ReadDir(curriculum_json_folder)

	if curriculum_read_err != nil {
		panic(curriculum_read_err)
	}

	all_curriculums := make([]c.Curriculum, 0)

	for i, json_file := range curriculum_json_files {
		curriculum_data := &c.Curriculum{}
		curriculum_data.CurriculumID = uint16(i + 1)

		json_file_path := path.Join(curriculum_json_folder, json_file.Name())
		byte_data, err := os.ReadFile(json_file_path)

		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(byte_data, &curriculum_data)
		if err != nil {
			panic(err)
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

	return all_curriculums
}
