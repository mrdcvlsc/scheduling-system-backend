package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Departments"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

type CurriculumMicroData struct {
	CurriculumID   uint16               `json:"CurriculumID"`
	CurriculumName string               `json:"CurriculumName"`
	CurriculumCode string               `json:"CurriculumCode"`
	YearLevels     []YearLevelMicroData `json:"YearLevels"`
}

type YearLevelMicroData struct {
	Name     string `json:"Name"`
	Sections int    `json:"Sections"`
}

/*
GET:

	"/all_departments"
*/
func GetAllDepartments(ctx *gin.Context) {
	all_departments, err_read_all_departments := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllDepartments()

	if err_read_all_departments != nil {
		log.Printf("error - %s", err_read_all_departments.Error())
		ctx.String(http.StatusInternalServerError, "we're currently unable to get the list of all departments")
		return
	}

	departments_json := make([]Departments.Department, 0)

	for _, department := range all_departments {
		departments_json = append(departments_json, Departments.Department{
			DepartmentID: department.DepartmentID,
			Name:         department.Name,
			Code:         department.Code,
		})
	}

	ctx.JSON(http.StatusOK, departments_json)
}

/*
GET:

	"/department_data?department_id=D&semester=[0-1]"
*/
func GetCurriculumsDataInDepartment(ctx *gin.Context) {

	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	selected_semester, is_valid_semester_param := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	all_curriculums, err_read_all_curriculum := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_read_all_curriculum != nil {
		ctx.String(http.StatusInternalServerError, "unable to read curriculums for that department")
		return
	}

	slice_of_curriculum_micro_data := make([]CurriculumMicroData, 0)

	for _, curriculum := range all_curriculums {

		curriculum_micro_data := &CurriculumMicroData{
			CurriculumID:   curriculum.CurriculumID,
			CurriculumName: curriculum.CurriculumName,
			CurriculumCode: curriculum.CurriculumCode,
			YearLevels:     make([]YearLevelMicroData, len(curriculum.YearLevels)),
		}

		for yl_idx, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue
			}

			for semester_idx, semester := range year_level.Semesters {

				if semester_idx != selected_semester {
					continue
				}

				curriculum_micro_data.YearLevels[yl_idx].Name = year_level.Name
				curriculum_micro_data.YearLevels[yl_idx].Sections = semester.Sections
			}
		}

		if department_id == int(curriculum.DepartmentID) {
			slice_of_curriculum_micro_data = append(slice_of_curriculum_micro_data, *curriculum_micro_data)
		}
	}

	ctx.JSON(http.StatusOK, slice_of_curriculum_micro_data)
}
