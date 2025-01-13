package Storage

import (
	"encoding/json"
	"os"
	"path"

	curriculum "github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	instructors "github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	room "github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
)

// this type is for development and testing only
type JsonFilePersistence struct{}

func (s *JsonFilePersistence) GetAllRooms() []room.Room {
	rooms_json_file := path.Join("..", "scheduling-system-temporary-data", "rooms.json")
	rooms_byte_data, err := os.ReadFile(rooms_json_file)

	if err != nil {
		panic(err)
	}

	all_rooms := make([]room.Room, 0)
	err = json.Unmarshal(rooms_byte_data, &all_rooms)

	if err != nil {
		panic(err)
	}

	for idx := range all_rooms {
		all_rooms[idx].RoomID = uint16(idx + 1)
	}

	return all_rooms
}

func (s *JsonFilePersistence) GetAllInstructors() []instructors.Instructor {
	instructors_json_file := path.Join("..", "scheduling-system-temporary-data", "instructors.json")
	instructors_byte_data, err := os.ReadFile(instructors_json_file)

	if err != nil {
		panic(err)
	}

	all_instructors := make([]instructors.Instructor, 0)
	err = json.Unmarshal(instructors_byte_data, &all_instructors)

	if err != nil {
		panic(err)
	}

	for idx := range all_instructors {
		all_instructors[idx].InstructorID = uint16(idx + 1)
	}

	return all_instructors
}

func (s *JsonFilePersistence) GetAllCurriculum() []curriculum.Curriculum {

	// map for getting subject id using subject code
	subject_code_map_id := make(map[string]uint16)

	// read all subject from all-subjects.json file

	subjects_json_file := path.Join("..", "scheduling-system-temporary-data", "all-subjects.json")
	subjects_byte_data, err := os.ReadFile(subjects_json_file)

	if err != nil {
		panic(err)
	}

	subject_data := []curriculum.Subject{}
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

	all_curriculums := make([]curriculum.Curriculum, 0)

	for i, json_file := range curriculum_json_files {
		curriculum_data := &curriculum.Curriculum{}
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
