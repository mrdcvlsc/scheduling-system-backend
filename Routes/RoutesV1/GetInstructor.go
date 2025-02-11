package RoutesV1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

/*
GET:

	"/instructor_erd?department_id=D"
*/
func GetDepartmentInstructorsEncodingResourceDefault(ctx *gin.Context) {
	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	default_encoding_resource, err_encoding_resource := GeneticAlgorithm.ReadDefaultEncodingResource(RouteGlobals.ResourcesPersistence)

	if err_encoding_resource != nil {
		log.Print(err_encoding_resource)
		ctx.String(http.StatusInternalServerError, "we're currently unable to get the default encoding resource of your department instructors")
		return
	}

	general_instructors := default_encoding_resource.DeptIdToInstructors[0]
	department_instructors := default_encoding_resource.DeptIdToInstructors[uint16(department_id)]

	combined_instructors := make([]Instructors.Instructor, 0)

	combined_instructors = append(combined_instructors, general_instructors...)
	combined_instructors = append(combined_instructors, department_instructors...)

	ctx.JSON(http.StatusOK, combined_instructors)
}

/*
GET:

	"/instructor_era?department_id=D&semester=[0-1]"
*/
func GetDepartmentInstructorsEncodingResourceAllocation(ctx *gin.Context) {
	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	selected_semester, is_valid_semester_param := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	university_schedule, has_obtained := ObtainUniversitySchedule(ctx, nil, selected_semester)

	if !has_obtained {
		return
	}

	curriculums, err_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.GetAllCurriculum()

	if err_curriculums != nil {
		log.Print(err_curriculums)
		ctx.String(http.StatusInternalServerError, "we're unable to retrieve the curriculum for the requested semester index")
		return
	}

	university_encoding_resource, err_generating := GeneticAlgorithm.GenerateEncodingResourceFromUniTimeTable(
		university_schedule,
		curriculums, selected_semester,
		RouteGlobals.ResourcesPersistence,
	)

	if err_generating != nil {
		log.Print(err_generating)
		ctx.String(http.StatusInternalServerError, "we're unable to generate the encoding resource for that semester university schedule")
		return
	}

	general_instructors := university_encoding_resource.DeptIdToInstructors[0]
	department_instructors := university_encoding_resource.DeptIdToInstructors[uint16(department_id)]

	combined_instructors := make([]Instructors.Instructor, 0)

	combined_instructors = append(combined_instructors, general_instructors...)
	combined_instructors = append(combined_instructors, department_instructors...)

	ctx.JSON(http.StatusOK, combined_instructors)
}
