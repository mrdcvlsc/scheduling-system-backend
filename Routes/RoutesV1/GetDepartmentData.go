package RoutesV1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

type DepartmentData struct {
	Curriculums []CurriculumMicroData `json:"Curriculums"`
}

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

func GetDepartmentData(ctx *gin.Context) {
	param_department_id := ctx.Query("department_id")

	if param_department_id == "" {
		ctx.String(http.StatusBadRequest, "mising 'department_id' parameter or parameter value")
		return
	}

	department_id, department_id_atoi_err := strconv.Atoi(param_department_id)

	if department_id_atoi_err != nil {
		ctx.String(http.StatusBadRequest, "invalid 'department_id' parameter value")
		return
	}

	param_semester := ctx.Query("semester")

	if param_semester == "" {
		ctx.String(http.StatusBadRequest, "mising 'semester' parameter or parameter value")
		return
	}

	selected_semester, semester_atoi_err := strconv.Atoi(param_semester)

	if semester_atoi_err != nil {
		ctx.String(http.StatusBadRequest, "invalid 'semester' parameter value")
		return
	}

	if selected_semester < 0 || selected_semester >= 2 {
		ctx.String(http.StatusBadRequest, "invalid 'semester' index value")
		return
	}

	curriculums, curriculum_err := RouteGlobals.ResourcesPersistence.ReaderService.GetAllCurriculum()

	if curriculum_err != nil {
		ctx.String(http.StatusInternalServerError, "unable to read curriculums for that department")
		return
	}

	department_data := &DepartmentData{}

	department_data.Curriculums = make([]CurriculumMicroData, 0)

	for _, curriculum := range curriculums {

		if curriculum.DepartmentID != uint16(department_id) {
			continue
		}

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
					Name:     year_level.Name,
					Sections: semester.Sections,
				}

				curriculum_micro_data.YearLevels = append(curriculum_micro_data.YearLevels, *year_level_micro_data)
			}
		}

		department_data.Curriculums = append(department_data.Curriculums, *curriculum_micro_data)
	}

	ctx.JSON(http.StatusOK, department_data)
}
