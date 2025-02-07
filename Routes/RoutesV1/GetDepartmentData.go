package RoutesV1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

type CurriculumMicroData struct {
	CurriculumID   uint16               `json:"CurriculumID"`
	CurriculumName string               `json:"CurriculumName"`
	CurriculumCode string               `json:"CurriculumCode"`
	YearLevels     []YearLevelMicroData `json:"YearLevels"`
}

type YearLevelMicroData struct {
	Name                            string `json:"Name"`
	SectionsUniversityScheduleIndex []int  `json:"Sections"`
}

// GET : /v1/department_data?department_id=D&semester=S
func GetDepartmentData(ctx *gin.Context) {

	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	selected_semester, is_valid_semester_param := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	curriculums, curriculum_err := RouteGlobals.ResourcesPersistence.ReaderService.GetAllCurriculum()

	if curriculum_err != nil {
		ctx.String(http.StatusInternalServerError, "unable to read curriculums for that department")
		return
	}

	slice_of_curriculum_micro_data := make([]CurriculumMicroData, 0)

	schedule_idx := 0

	for _, curriculum := range curriculums {

		curriculum_micro_data := &CurriculumMicroData{
			CurriculumID:   curriculum.CurriculumID,
			CurriculumName: curriculum.CurriculumName,
			CurriculumCode: curriculum.CurriculumCode,
			YearLevels:     make([]YearLevelMicroData, 0),
		}

		for _, year_level := range curriculum.YearLevels {

			if !year_level.IsActive {
				continue
			}

			for semester_idx, semester := range year_level.Semesters {

				if semester_idx != selected_semester {
					continue
				}

				year_level_micro_data := &YearLevelMicroData{
					Name:                            year_level.Name,
					SectionsUniversityScheduleIndex: make([]int, 0),
				}

				for section_idx := 0; section_idx < semester.Sections; section_idx++ {

					if department_id == int(curriculum.DepartmentID) {
						year_level_micro_data.SectionsUniversityScheduleIndex = append(
							year_level_micro_data.SectionsUniversityScheduleIndex,
							schedule_idx,
						)
					}

					schedule_idx++
				}

				if department_id == int(curriculum.DepartmentID) {
					curriculum_micro_data.YearLevels = append(curriculum_micro_data.YearLevels, *year_level_micro_data)
				}
			}
		}

		if department_id == int(curriculum.DepartmentID) {
			slice_of_curriculum_micro_data = append(slice_of_curriculum_micro_data, *curriculum_micro_data)
		}
	}

	ctx.JSON(http.StatusOK, slice_of_curriculum_micro_data)
}
