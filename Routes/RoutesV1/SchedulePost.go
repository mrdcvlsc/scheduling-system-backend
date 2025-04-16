package RoutesV1

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

const DEFAULT_INITIAL_REQUEST_COUNT uint = 30

/*
POST:

	"/generate_schedule?semester=[0-1]&department_id=[N>0]"
*/
func RequestGenerateSchedule(ctx *gin.Context) {

	semester, is_valid_semester_idx := IsValidParameterSemesterIndex(ctx)

	if !is_valid_semester_idx {
		return
	}

	department_id, is_valid_department_id := IsValidParameterDepartmentID(ctx)

	if !is_valid_department_id {
		return
	}

	response_msg := ""
	response_status := http.StatusAccepted

	if RouteGlobals.PushNewDeptToDeptSchedGenQueue(RouteGlobals.DeptSchedGenKey{
		DepartmentID: uint16(department_id),
		Semester:     semester,
	}) {
		response_msg += fmt.Sprintf("department with id %d was added to the schedule generation queue,", department_id)
	} else {
		response_msg += fmt.Sprintf("the department with id %d is already in schedule generation queue,", department_id)
		response_status = http.StatusContinue
	}

	if !RouteGlobals.IsGeneratingSchedule.Load() {
		response_msg += " the schedule generation function has started"
		go encode_schedule()
	} else {
		response_msg += " the schedule generation function is already running"
	}

	log.Printf("GenerateSchedule: generating schedule")
	ctx.String(response_status, response_msg)
}

func encode_schedule() {
	RouteGlobals.IsGeneratingSchedule.Store(true)
	defer RouteGlobals.IsGeneratingSchedule.Store(false)

	RouteGlobals.ReindexUniSchedMutex.Lock()
	defer RouteGlobals.ReindexUniSchedMutex.Unlock()

	log.Println("encode_schedule [0]: generating schedule...")

	////////////////////////////////////////////////////////////////////////////////////////

	curriculums, err_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_curriculums != nil {
		log.Fatal("encode_schedule [1]:", err_curriculums)
	}

	dept_id_to_department, err_dept_id_to_department := GeneticAlgorithm.GenerateMapDeptIdToDepartment(RouteGlobals.ResourcesPersistence)

	if err_dept_id_to_department != nil {
		log.Fatal("encode_schedule [2]:", err_dept_id_to_department)
	}

	////////////////////////////////////////////////////////////////////////////////////////

	var generated_encoding_resource *GeneticAlgorithm.EncodingResource
	var err_gen_encoding_resource error

	for {
		// get the first department in the queue that requested to generate a schedule for a specific semester

		department_to_encode, semester_to_encode, err_pop_from_queue := RouteGlobals.PopDepartmentToEncodeFromSchedGenQueue()

		if err_pop_from_queue != nil || !has_department_to_encode(department_to_encode) {
			break // no more department and semester in the schedule generation queue to be encoded in the schedules
		}

		department_id := uint16(0)

		for k := range department_to_encode {
			department_id = k
		}

		if department_id <= 0 {
			continue // a sanity check, don't generate schedules for department id that is less than 1
		}

		log.Println("encode_schedule [2.1]: pop latest task from queue")

		// get the current university schedules for the specific semester requested by the first department in the queue

		university_schedule, err_obtain_uni_sched_no_ctx := ObtainUniversityScheduleNoContextNoHorizontalValidation(semester_to_encode)

		if err_obtain_uni_sched_no_ctx != nil {
			RouteGlobals.SetDepartmentsLastScheduleGenerationResult(
				RouteGlobals.DeptSchedGenKey{DepartmentID: department_id, Semester: semester_to_encode},
				RouteGlobals.ScheduleGenerationLastResult{
					Status:  false,
					Message: err_obtain_uni_sched_no_ctx.Error(),
				},
			)

			continue
		}

		log.Println("encode_schedule [2.2]: obtain schedule for the specified semester")

		// generate the encoding resource for the obtained university schedule

		generated_encoding_resource, err_gen_encoding_resource = GeneticAlgorithm.GenerateEncodingResourceFromUniTimeTable(
			university_schedule, curriculums, semester_to_encode, RouteGlobals.ResourcesPersistence,
		)

		if err_gen_encoding_resource != nil {
			log.Printf(
				"encode_schedule [2.2.2]: error generating encoding resource for %s %s",
				dept_id_to_department[department_id].Name, Curriculum.SEMESTER_INDEX_NAME[semester_to_encode],
			)

			RouteGlobals.SetDepartmentsLastScheduleGenerationResult(
				RouteGlobals.DeptSchedGenKey{DepartmentID: department_id, Semester: semester_to_encode},
				RouteGlobals.ScheduleGenerationLastResult{
					Status:  false,
					Message: err_gen_encoding_resource.Error(),
				},
			)

			continue
		}

		log.Println("encode_schedule [2.3]: generated new schedule encoding for the department specific semester")

		// encode a new schedule in the obtained university schedule for the specific department

		max_retries := 50
		var retry int

		log.Println("encode_schedule [2.4]: trying to generate the schedule")

		for retry = 0; retry < max_retries; retry++ {
			new_encoded_university_schedule, new_encoding_resource, err_encode_individual_genome := GeneticAlgorithm.EncodeIndividualGenome(
				university_schedule,
				curriculums, dept_id_to_department,
				generated_encoding_resource, department_to_encode,
				semester_to_encode, 0,
			)

			if err_encode_individual_genome != nil {
				if retry == max_retries-1 {
					RouteGlobals.SetDepartmentsLastScheduleGenerationResult(
						RouteGlobals.DeptSchedGenKey{DepartmentID: department_id, Semester: semester_to_encode},
						RouteGlobals.ScheduleGenerationLastResult{
							Status:  false,
							Message: err_encode_individual_genome.Error(),
						},
					)

					log.Print("encode_schedule [2.5]: max retires - unable to generate schedules")
					break
				}

				continue
			}

			log.Print("encode_schedule [2.6]: new encoded university schedules generated")

			// perform checks and validation to the new encoded schedule for the specific department

			if new_encoded_university_schedule == nil {
				log.Print("encode_schedule [3]: unable to generate schedules")

				RouteGlobals.SetDepartmentsLastScheduleGenerationResult(
					RouteGlobals.DeptSchedGenKey{DepartmentID: department_id, Semester: semester_to_encode},
					RouteGlobals.ScheduleGenerationLastResult{
						Status: false,
						Message: fmt.Sprintf(
							"unable to generate schedules for the department with id %d %s",
							department_id, Curriculum.SEMESTER_INDEX_NAME[semester_to_encode],
						),
					},
				)

				break
			}

			if reflect.DeepEqual(university_schedule, new_encoded_university_schedule) {
				log.Print("encode_schedule [3.-1]: (equal) result no changes made after generating new schedule encoding")
			} else {
				log.Print("encode_schedule [3.-1]: (not-equal) new changes are made after generating new schedule encoding")
			}

			if GeneticAlgorithm.IsEqualEncodingResource(generated_encoding_resource, new_encoding_resource) {
				log.Print("encode_schedule [3.-2]: (equal) result no changes made after generating new encoding resources")
			} else {
				log.Print("encode_schedule [3.-2]: (not-equal) new changes are made after generating new encoding resources")
			}

			vertical_overlaps := false
			for _, e := range new_encoded_university_schedule.VerticalValidation(RouteGlobals.ResourcesPersistence) {
				RouteGlobals.SetDepartmentsLastScheduleGenerationResult(
					RouteGlobals.DeptSchedGenKey{DepartmentID: department_id, Semester: semester_to_encode},
					RouteGlobals.ScheduleGenerationLastResult{
						Status: false,
						Message: fmt.Sprintf(
							"error vertical overlaps detected: %s", e.Error(),
						),
					},
				)

				vertical_overlaps = true
				log.Print("encode_schedule [3.1]: error vertical overlaps")
				break
			}

			if vertical_overlaps {
				break
			}

			horizontal_overlaps := false

			errs_horizontal_overlaps := new_encoded_university_schedule.HorizontalValidation(RouteGlobals.ResourcesPersistence, department_to_encode, semester_to_encode)

			for _, e := range new_encoded_university_schedule.HorizontalValidation(RouteGlobals.ResourcesPersistence, department_to_encode, semester_to_encode) {
				RouteGlobals.SetDepartmentsLastScheduleGenerationResult(
					RouteGlobals.DeptSchedGenKey{DepartmentID: department_id, Semester: semester_to_encode},
					RouteGlobals.ScheduleGenerationLastResult{
						Status: false,
						Message: fmt.Sprintf(
							"error horizontal overlaps detected: %s", e.Error(),
						),
					},
				)

				horizontal_overlaps = true
				log.Printf("encode_schedule [3.2]: error horizontal overlaps \n\n%s\n\n", e.Error())

				Utils.PrettyPrint(errs_horizontal_overlaps)

				fmt.Print("\n\n")

				break
			}

			if horizontal_overlaps {
				break
			}

			err_save_schedules := RouteGlobals.SchedulePersistence.SaveService.SaveSchedules(new_encoded_university_schedule, semester_to_encode)

			if err_save_schedules != nil {
				log.Print("encode_schedule [4]:", err_save_schedules.Error())

				RouteGlobals.SetDepartmentsLastScheduleGenerationResult(
					RouteGlobals.DeptSchedGenKey{DepartmentID: department_id, Semester: semester_to_encode},
					RouteGlobals.ScheduleGenerationLastResult{
						Status: false,
						Message: fmt.Sprintf(
							"error saving schedule: %s", err_save_schedules.Error(),
						),
					},
				)
				break
			}

			log.Print("encode_schedule [4]: university schedule saved")

			err_set_cache := RouteGlobals.SetCachedUniversitySchedule(semester_to_encode, new_encoded_university_schedule)

			if err_set_cache != nil {
				log.Print("encode_schedule [5]:", err_set_cache.Error())

				RouteGlobals.SetDepartmentsLastScheduleGenerationResult(
					RouteGlobals.DeptSchedGenKey{DepartmentID: department_id, Semester: semester_to_encode},
					RouteGlobals.ScheduleGenerationLastResult{
						Status: false,
						Message: fmt.Sprintf(
							"error caching schedule: %s", err_set_cache.Error(),
						),
					},
				)

				break
			}

			// specific department schedule generation done

			break
		}

		log.Println("encode_schedule [5.1]: schedule generation loop done")
	}

	log.Println("encode_schedule [6]: all department schedule generation requests are done...")
}

func has_department_to_encode(department_to_encode map[uint16]bool) bool {
	has_department_to_encode_result := false

	if department_to_encode == nil {
		return has_department_to_encode_result
	}

	if len(department_to_encode) == 0 {
		return has_department_to_encode_result
	}

	for _, v := range department_to_encode {
		has_department_to_encode_result = has_department_to_encode_result || v
	}

	return has_department_to_encode_result
}
