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

/*
POST:

	"/instructor_update"
*/
func PostInstructor(ctx *gin.Context) {
	update_instructor_with_time_str := Instructors.InstructorWithTimeString{}

	if err := ctx.BindJSON(&update_instructor_with_time_str); err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read instructor update data")
		return
	}

	update_instructor := Instructors.Instructor{
		InstructorID:  update_instructor_with_time_str.InstructorID,
		DepartmentID:  update_instructor_with_time_str.DepartmentID,
		FirstName:     update_instructor_with_time_str.FirstName,
		MiddleInitial: update_instructor_with_time_str.MiddleInitial,
		LastName:      update_instructor_with_time_str.LastName,
	}

	update_instructor.Time.StringParse(update_instructor_with_time_str.Time)

	err := RouteGlobals.ResourcesPersistence.WriterService.UpdateInstructor(update_instructor)

	if err != nil {
		ctx.String(http.StatusBadRequest, "we are unable to properly read instructor update data")
		return
	}

	ctx.String(http.StatusOK, "instructor update successful")
}

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

	combined_instructor_stringify_time := make([]Instructors.InstructorWithTimeString, 0)

	for _, instructor := range combined_instructors {
		time_stringify := make([]string, 0)

		for _, limb := range instructor.Time {
			time_stringify = append(time_stringify, strconv.FormatUint(limb, 10))
		}

		combined_instructor_stringify_time = append(combined_instructor_stringify_time, Instructors.InstructorWithTimeString{
			InstructorID:  instructor.InstructorID,
			DepartmentID:  instructor.DepartmentID,
			FirstName:     instructor.FirstName,
			MiddleInitial: instructor.MiddleInitial,
			LastName:      instructor.LastName,
			Time:          time_stringify,
		})
	}

	ctx.JSON(http.StatusOK, combined_instructor_stringify_time)
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

	general_instructors := university_encoding_resource.DeptIdToInstructors[0]

	department_instructors := university_encoding_resource.DeptIdToInstructors[uint16(department_id)]

	combined_instructors := make([]Instructors.Instructor, 0)

	combined_instructors = append(combined_instructors, general_instructors...)
	combined_instructors = append(combined_instructors, department_instructors...)

	combined_instructor_stringify_time := make([]Instructors.InstructorWithTimeString, 0)

	for _, instructor := range combined_instructors {
		combined_instructor_stringify_time = append(combined_instructor_stringify_time, Instructors.InstructorWithTimeString{
			InstructorID:  instructor.InstructorID,
			DepartmentID:  instructor.DepartmentID,
			FirstName:     instructor.FirstName,
			MiddleInitial: instructor.MiddleInitial,
			LastName:      instructor.LastName,
			Time:          instructor.Time.Stringify(),
		})
	}

	ctx.JSON(http.StatusOK, combined_instructor_stringify_time)
}
