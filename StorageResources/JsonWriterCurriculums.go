package StorageResources

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func (s *JsonWriter) CreateCurriculum(new_curriculum Curriculum.Curriculum) error {

	if new_curriculum.CurriculumID != 0 {
		return errors.New("cannot create a new curriculum with a non zero curriculum ID because that would overwrite a curriculum")
	}

	all_curriculums, err_read := json_read_all_curriculums()

	if err_read != nil {
		return err_read
	}

	for _, curriculum := range all_curriculums {
		if strings.EqualFold(Utils.RemoveWhiteSpace(curriculum.CurriculumCode), Utils.RemoveWhiteSpace(new_curriculum.CurriculumCode)) {
			return errors.New("cannot create a new curriculum with that curriculum code")
		}

		if strings.EqualFold(curriculum.CurriculumName, new_curriculum.CurriculumName) {
			return errors.New("cannot create a new curriculum with that curriculum name")
		}
	}

	err_save_curriculums := json_save_curriculum(
		fmt.Sprintf("%s.json", Utils.RemoveWhiteSpace(new_curriculum.CurriculumCode)),
		Curriculum.Curriculum{
			CurriculumID:   all_curriculums[len(all_curriculums)-1].CurriculumID + 1,
			CurriculumName: new_curriculum.CurriculumName,
			CurriculumCode: new_curriculum.CurriculumCode,
			DepartmentID:   new_curriculum.DepartmentID,
			YearLevels:     new_curriculum.YearLevels,
		},
	)

	if err_save_curriculums != nil {
		return err_save_curriculums
	}

	return nil
}

func (s *JsonWriter) UpdateCurriculum(curriculum_old_name string, curriculum_new Curriculum.Curriculum) error {

	if curriculum_new.CurriculumID == 0 {
		return errors.New("parameter argument missing invalid CurriculumID")
	}

	all_curriculums, err_read := json_read_all_curriculums()

	if err_read != nil {
		return err_read
	}

	has_id := false

	for _, curriculum := range all_curriculums {
		if curriculum.CurriculumID == curriculum_new.CurriculumID {
			has_id = true
			break
		}
	}

	if !has_id {
		return errors.New("instructor to update does not exist in the json file")
	}

	err_save_curriculums := json_edit_curriculum(
		curriculum_old_name,
		fmt.Sprintf("%s.json", Utils.RemoveWhiteSpace(curriculum_new.CurriculumCode)),
		curriculum_new,
	)

	if err_save_curriculums != nil {
		return err_save_curriculums
	}

	return nil
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
