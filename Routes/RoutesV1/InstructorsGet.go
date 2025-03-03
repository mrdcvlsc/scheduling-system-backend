package RoutesV1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Instructors"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
)

type InstructorTablePage struct {
	Instructors      []Instructors.InstructorWithTimeString `json:"Instructors"`
	TotalInstructors int                                    `json:"TotalInstructors"`
}

/*
GET:

	"/instructors/d?department_id=D&page_size=[N>0]&page[0-N>0]"
*/
func GetDepartmentInstructorsDefaults(ctx *gin.Context) {
	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)
	if !is_valid_department_id_param {
		return
	}

	page_size, is_valid_page_size_param := IsValidPageSize(ctx)
	if !is_valid_page_size_param {
		return
	}

	page, is_valid_page_param := IsValidPage(ctx)
	if !is_valid_page_param {
		return
	}

	default_encoding_resource, err_encoding_resource := GeneticAlgorithm.ReadDefaultEncodingResource(RouteGlobals.ResourcesPersistence)

	if err_encoding_resource != nil {
		log.Print(err_encoding_resource)
		ctx.String(http.StatusInternalServerError, "we're currently unable to get the default encoding resource of your department instructors")
		return
	}

	department_instructors := default_encoding_resource.DeptIdToInstructors[uint16(department_id)]
	department_instructors_stringify_time := make([]Instructors.InstructorWithTimeString, 0)

	for i, instructor := range department_instructors {
		if i < (page_size * page) {
			continue
		}

		time_stringify := make([]string, 0)

		for _, limb := range instructor.Time {
			time_stringify = append(time_stringify, strconv.FormatUint(limb, 10))
		}

		department_instructors_stringify_time = append(department_instructors_stringify_time, Instructors.InstructorWithTimeString{
			InstructorID:  instructor.InstructorID,
			DepartmentID:  instructor.DepartmentID,
			FirstName:     instructor.FirstName,
			MiddleInitial: instructor.MiddleInitial,
			LastName:      instructor.LastName,
			Time:          time_stringify,
		})

		if len(department_instructors_stringify_time) >= page_size {
			break
		}
	}

	instructor_table_page := &InstructorTablePage{
		Instructors:      department_instructors_stringify_time,
		TotalInstructors: len(department_instructors),
	}

	ctx.JSON(http.StatusOK, instructor_table_page)
}

/*
GET:

	"/instructors/a?department_id=D&semester=[0-1]&page_size=[N>0]&page[0-N>0]"
*/
func GetDepartmentInstructorsAllocated(ctx *gin.Context) {
	department_id, is_valid_department_id_param := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id_param {
		return
	}

	selected_semester, is_valid_semester_param := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_param {
		return
	}

	page_size, is_valid_page_size_param := IsValidPageSize(ctx)
	if !is_valid_page_size_param {
		return
	}

	page, is_valid_page_param := IsValidPage(ctx)
	if !is_valid_page_param {
		return
	}

	university_schedule, has_obtained := ObtainUniversitySchedule(ctx, nil, selected_semester)

	if !has_obtained {
		return
	}

	curriculums, err_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

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

	department_instructors := university_encoding_resource.DeptIdToInstructors[uint16(department_id)]

	department_instructors_stringify_time := make([]Instructors.InstructorWithTimeString, 0)

	for i, instructor := range department_instructors {
		if i < (page_size * page) {
			continue
		}

		department_instructors_stringify_time = append(department_instructors_stringify_time, Instructors.InstructorWithTimeString{
			InstructorID:  instructor.InstructorID,
			DepartmentID:  instructor.DepartmentID,
			FirstName:     instructor.FirstName,
			MiddleInitial: instructor.MiddleInitial,
			LastName:      instructor.LastName,
			Time:          instructor.Time.Stringify(),
		})

		if len(department_instructors_stringify_time) >= page_size {
			break
		}
	}

	instructor_table_page := &InstructorTablePage{
		Instructors:      department_instructors_stringify_time,
		TotalInstructors: len(department_instructors),
	}

	ctx.JSON(http.StatusOK, instructor_table_page)
}
