package StorageResources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Rooms"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

// this type is for development and testing only
type JsonWriter struct{}

func (s *JsonWriter) CreateDepartment(new_department Departments.Department) error {

	if new_department.DepartmentID != 0 {
		return errors.New("cannot create a newdepartment with a non zero department ID because that would overwrite a department")
	}

	all_departments, err_read := json_read_all_departments()

	if err_read != nil {
		return err_read
	}

	all_departments = append(all_departments, Departments.Department{
		DepartmentID: all_departments[len(all_departments)-1].DepartmentID + 1,
		Code:         new_department.Code,
		Name:         new_department.Name,
	})

	err_save_departments := json_save_all_departments(all_departments)

	if err_save_departments != nil {
		return err_save_departments
	}

	return nil
}

func (s *JsonWriter) UpdateDepartment(department_to_update Departments.Department) error {

	if department_to_update.DepartmentID == 0 {
		return errors.New("parameter argument missing invalid department ID")
	}

	all_departments, err_read := json_read_all_departments()

	if err_read != nil {
		return err_read
	}

	has_id := false
	to_update_idx := -1

	for idx, department := range all_departments {
		if department.DepartmentID == department_to_update.DepartmentID {
			has_id = true
			to_update_idx = idx
			break
		}
	}

	if !has_id {
		return errors.New("department to update does not exist in the json file")
	}

	all_departments[to_update_idx] = department_to_update

	err_save_departments := json_save_all_departments(all_departments)

	if err_save_departments != nil {
		return err_save_departments
	}

	return nil
}

func (s *JsonWriter) CreateRoom(new_room Rooms.Room) error {

	if new_room.RoomID != 0 {
		return errors.New("cannot create a new room with a non zero room ID because that would overwrite a room item")
	}

	all_rooms, err_read := json_read_all_rooms()

	if err_read != nil {
		return err_read
	}

	all_rooms = append(all_rooms, Rooms.Room{
		RoomID:       all_rooms[len(all_rooms)-1].RoomID + 1,
		DepartmentID: new_room.DepartmentID,
		Capacity:     new_room.Capacity,
		RoomType:     new_room.RoomType,
		Name:         new_room.Name,
	})

	err_save_instructors := json_save_all_rooms(all_rooms)

	if err_save_instructors != nil {
		return err_save_instructors
	}

	return nil
}

func (s *JsonWriter) UpdateRoom(room_to_update Rooms.Room) error {

	if room_to_update.RoomID == 0 {
		return errors.New("parameter argument missing invalid room ID")
	}

	all_rooms, err_read := json_read_all_rooms()

	if err_read != nil {
		return err_read
	}

	has_id := false
	to_update_idx := -1

	for idx, room := range all_rooms {
		if room.RoomID == room_to_update.RoomID {
			has_id = true
			to_update_idx = idx
			break
		}
	}

	if !has_id {
		return errors.New("room to update does not exist in the json file")
	}

	all_rooms[to_update_idx] = Rooms.Room{
		RoomID:       room_to_update.RoomID,
		DepartmentID: room_to_update.DepartmentID,
		Capacity:     room_to_update.Capacity,
		RoomType:     room_to_update.RoomType,
		Name:         room_to_update.Name,
	}

	err_save_rooms := json_save_all_rooms(all_rooms)

	if err_save_rooms != nil {
		return err_save_rooms
	}

	return nil
}

func (s *JsonWriter) CreateSubject(new_subject Curriculum.Subject) error {

	if new_subject.ID != 0 {
		return errors.New("cannot create a new room with a non zero room ID because that would overwrite a room item")
	}

	all_subject, err_read := json_read_all_subjects()

	if err_read != nil {
		return err_read
	}

	all_subject = append(all_subject, Curriculum.Subject{
		ID:                    all_subject[len(all_subject)-1].ID + 1,
		Code:                  new_subject.Code,
		Name:                  new_subject.Name,
		LecHours:              new_subject.LecHours,
		LabHours:              new_subject.LabHours,
		BitFlags:              new_subject.BitFlags,
		DesignatedInstructors: new_subject.DesignatedInstructors,
	})

	err_save_subjects := json_save_all_subjects(all_subject)

	if err_save_subjects != nil {
		return err_save_subjects
	}

	return nil
}

func (s *JsonWriter) UpdateSubject(subject_to_update Curriculum.Subject) error {

	if subject_to_update.ID == 0 {
		return errors.New("parameter argument missing invalid subject ID")
	}

	all_subjects, err_read := json_read_all_subjects()

	if err_read != nil {
		return err_read
	}

	all_curriculums, err_read_all_curriculums := json_read_all_curriculums()

	if err_read_all_curriculums != nil {
		return err_read_all_curriculums
	}

	has_id := false
	to_update_idx := -1

	for idx, subject := range all_subjects {
		if subject.ID == subject_to_update.ID {
			has_id = true
			to_update_idx = idx
			break
		}
	}

	if !has_id {
		return errors.New("room to update does not exist in the json file")
	}

	updated_subject := Curriculum.Subject{
		ID:                    subject_to_update.ID,
		Code:                  subject_to_update.Code,
		Name:                  subject_to_update.Name,
		LecHours:              subject_to_update.LecHours,
		LabHours:              subject_to_update.LabHours,
		BitFlags:              subject_to_update.BitFlags,
		DesignatedInstructors: subject_to_update.DesignatedInstructors,
	}

	all_subjects[to_update_idx] = updated_subject

	for curriculum_idx := range all_curriculums {
		for yrlvl_idx := range all_curriculums[curriculum_idx].YearLevels {
			for semester_idx := range all_curriculums[curriculum_idx].YearLevels[yrlvl_idx].Semesters {
				for subject_idx := range all_curriculums[curriculum_idx].YearLevels[yrlvl_idx].Semesters[semester_idx].Subjects {
					curriculum_subject_id := all_curriculums[curriculum_idx].YearLevels[yrlvl_idx].Semesters[semester_idx].Subjects[subject_idx].ID
					if subject_to_update.ID == curriculum_subject_id {
						all_curriculums[curriculum_idx].YearLevels[yrlvl_idx].Semesters[semester_idx].Subjects[subject_idx] = updated_subject
					}
				}
			}
		}
	}

	err_save_rooms := json_save_all_subjects(all_subjects)

	if err_save_rooms != nil {
		return err_save_rooms
	}

	err_save_curriculums := json_save_all_curriculums(all_curriculums)

	if err_save_curriculums != nil {
		return err_save_curriculums
	}

	return nil
}

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

func (s *JsonWriter) CreateInstructor(new_instructor Instructors.Instructor) error {

	if new_instructor.InstructorID != 0 {
		return errors.New("cannot create a new instructor with a non zero instructor ID because it will overwrite an instructor")
	}

	instructors_with_time_str, err_read := json_read_all_instructors_with_time_string()

	if err_read != nil {
		return err_read
	}

	instructors_with_time_str = append(instructors_with_time_str, Instructors.InstructorWithTimeString{
		InstructorID:  instructors_with_time_str[len(instructors_with_time_str)-1].InstructorID + 1,
		DepartmentID:  new_instructor.DepartmentID,
		FirstName:     new_instructor.FirstName,
		MiddleInitial: new_instructor.MiddleInitial,
		LastName:      new_instructor.LastName,
		Time:          new_instructor.Time.Stringify(),
	})

	err_save_instructors := json_save_all_instructors_with_time_string(instructors_with_time_str)

	if err_save_instructors != nil {
		return err_save_instructors
	}

	return nil
}

func (s *JsonWriter) UpdateInstructor(instructor_to_update Instructors.Instructor) error {

	if instructor_to_update.InstructorID == 0 {
		return errors.New("parameter argument missing invalid instructor ID")
	}

	instructors_with_time_str, err_read := json_read_all_instructors_with_time_string()

	if err_read != nil {
		return err_read
	}

	has_id := false
	to_update_idx := -1

	for idx, instructor_w_t_str := range instructors_with_time_str {
		if instructor_w_t_str.InstructorID == instructor_to_update.InstructorID {
			has_id = true
			to_update_idx = idx
			break
		}
	}

	if !has_id {
		return errors.New("instructor to update does not exist in the json file")
	}

	instructors_with_time_str[to_update_idx] = Instructors.InstructorWithTimeString{
		InstructorID:  instructor_to_update.InstructorID,
		DepartmentID:  instructor_to_update.DepartmentID,
		FirstName:     instructor_to_update.FirstName,
		MiddleInitial: instructor_to_update.MiddleInitial,
		LastName:      instructor_to_update.LastName,
		Time:          instructor_to_update.Time.Stringify(),
	}

	err_save_instructors := json_save_all_instructors_with_time_string(instructors_with_time_str)

	if err_save_instructors != nil {
		return err_save_instructors
	}

	return nil
}

func (s *JsonWriter) DeleteInstructor(instructor_id uint16) error {

	if instructor_id == 0 {
		return errors.New("parameter argument missing invalid instructor ID")
	}

	instructors_with_time_str, err_read := json_read_all_instructors_with_time_string()

	if err_read != nil {
		return err_read
	}

	instructors_with_time_string_deleted := make([]Instructors.InstructorWithTimeString, 0)

	has_id := false

	for _, instructor_w_t_str := range instructors_with_time_str {
		if instructor_w_t_str.InstructorID == instructor_id {
			has_id = true
		} else {
			instructors_with_time_string_deleted = append(instructors_with_time_string_deleted, instructor_w_t_str)
		}
	}

	if !has_id {
		return errors.New("instructor to update does not exist in the json file")
	}

	err_save_instructors := json_save_all_instructors_with_time_string(instructors_with_time_string_deleted)

	if err_save_instructors != nil {
		return err_save_instructors
	}

	return nil
}

func (s *JsonWriter) DeleteRoom(room_id uint16) error {

	if room_id == 0 {
		return errors.New("parameter argument missing invalid room ID")
	}

	all_rooms, err_read := json_read_all_rooms()

	if err_read != nil {
		return err_read
	}

	rooms_deleted := make([]Rooms.Room, 0)

	has_id := false

	for _, room := range all_rooms {
		if room.RoomID == room_id {
			has_id = true
		} else {
			rooms_deleted = append(rooms_deleted, room)
		}
	}

	if !has_id {
		return errors.New("room to update does not exist in the json file")
	}

	err_save_rooms := json_save_all_rooms(rooms_deleted)

	if err_save_rooms != nil {
		return err_save_rooms
	}

	return nil
}
