package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Const"
	"github.com/mrdcvlsc/scheduling-system-backend/Resources/Curriculum"
	"github.com/mrdcvlsc/scheduling-system-backend/RouteGlobals"
	"github.com/mrdcvlsc/scheduling-system-backend/Routes/RoutesV1"
	"github.com/mrdcvlsc/scheduling-system-backend/Routes/RoutesV2"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageResources"
	"github.com/mrdcvlsc/scheduling-system-backend/StorageSchedule"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

var SessionStore = cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))

func TestIntegrationEditCurriculumSection(t *testing.T) {

	// initialize the router

	// setup_router()
	router := setup_router()

	// generate university schedules for all semesters

	departments, err_read_departments := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllDepartments()

	if err_read_departments != nil {
		t.Fatalf("Failed to read departments: %v", err_read_departments)
	}

	for semester := range Curriculum.SUPPORTED_SEMESTERS {
		for _, department := range departments {
			request := httptest.NewRequest(
				http.MethodPost, fmt.Sprintf(
					"/v1/generate_schedule?semester=%d&department_id=%d",
					semester, department.DepartmentID,
				),
				http.NoBody,
			)

			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code < http.StatusOK || response.Code >= http.StatusMultipleChoices {
				t.Fatalf("Failed to generate schedules: status code %d, body: %s", response.Code, response.Body.String())
			}
		}
	}

	// wait for the schedules to be generated

	for {
		request := httptest.NewRequest(http.MethodGet, "/v1/gen_status", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code < http.StatusOK || response.Code >= http.StatusMultipleChoices {
			t.Fatalf("Failed to check generation status: status code %d, body: %s", response.Code, response.Body.String())
		}

		var get_body struct {
			IsGenerating bool `json:"status"`
		}

		if err := json.Unmarshal(response.Body.Bytes(), &get_body); err != nil {
			t.Fatalf("Failed to parse generation status response: %v", err)
		}

		if !get_body.IsGenerating {
			break
		}

		time.Sleep(3 * time.Second)
		t.Log("Waiting for schedule generation to finish...")
	}

	// edit curriculum section counts

	curriculums, err_read_curriculums := RouteGlobals.ResourcesPersistence.ReaderService.ReadAllCurriculum()

	if err_read_curriculums != nil {
		t.Fatalf("Failed to read curriculums: %v", err_read_curriculums)
	}

	for c, curriculum := range curriculums {
		for y, year_level := range curriculum.YearLevels {
			for s := range year_level.Semesters {
				curriculums[c].YearLevels[y].Semesters[s].Sections = Utils.RandomInRange(0, 7)
			}
		}
	}

	for _, curriculum := range curriculums {

		json_curriculum, err := json.Marshal(curriculum)

		if err != nil {
			panic(err)
		}

		request := httptest.NewRequest(http.MethodPatch, "/v1/curriculum_update", bytes.NewBuffer(json_curriculum))

		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code < http.StatusOK || response.Code >= http.StatusMultipleChoices {
			t.Fatalf("Failed to edit curriculum sections: status code %d, body: %s", response.Code, response.Body.String())
		}
	}

	// validate the schedules

	t.Log("Validating...")

	for semester := range Curriculum.SUPPORTED_SEMESTERS {
		university_schedule, err_load_sched := RouteGlobals.SchedulePersistence.LoadService.LoadSchedules(semester)

		if err_load_sched != nil {
			t.Fatalf("Failed to load schedules: %v", err_load_sched)
		}

		GeneticAlgorithm.IterateSectionsWeekSchedule(university_schedule, curriculums, semester, nil, nil,
			func(indicies GeneticAlgorithm.IterIndices, values GeneticAlgorithm.IterValues) GeneticAlgorithm.IterReturnType {

				subject_from_schedule := make(map[uint16]int)
				for _, subject := range values.WeekSched.GetWeekSubjectsJSON() {
					_, has_id := subject_from_schedule[subject.SubjectID]
					if !has_id {
						subject_from_schedule[subject.SubjectID] = subject.TimeSlotSize
					} else {
						subject_from_schedule[subject.SubjectID] += subject.TimeSlotSize
					}
				}

				subject_from_curriculum := make(map[uint16]int)
				for _, subject := range values.Semester.Subjects {
					subject_from_curriculum[subject.ID] = int(subject.LecHours+subject.LabHours) * Const.N_HOUR_TIME_SLOTS
				}

				if (len(subject_from_schedule) != len(subject_from_curriculum)) && indicies.Section < 4 {
					t.Fatalf(
						"Mismatch in number of subjects, %s, %s, %s, section %s : schedule has %d, curriculum has %d",
						values.Curriculum.CurriculumCode,
						values.YearLevel.Name, Curriculum.SEMESTER_INDEX_NAME[indicies.Semester],
						Curriculum.SECTION[indicies.Section],
						len(subject_from_schedule),
						len(subject_from_curriculum),
					)
				}

				for k, v := range subject_from_schedule {
					if subject_from_schedule[k] != subject_from_curriculum[k] && indicies.Section < 4 {
						t.Fatalf(
							"Mismatch in %s, %s, %s, section %s, subject id %d: schedule has %d time slots, while curriculum has %d time slots",
							values.Curriculum.CurriculumCode,
							values.YearLevel.Name, Curriculum.SEMESTER_INDEX_NAME[indicies.Semester],
							Curriculum.SEMESTER_INDEX_NAME[indicies.Section],
							k, v, subject_from_curriculum[k],
						)
					}
				}

				if indicies.Section >= 4 && len(subject_from_schedule) != 0 {
					t.Fatalf(
						"%s, %s, %s, section %s has %d subjects even though it should not contain any because it is a new section",
						values.Curriculum.CurriculumCode,
						values.YearLevel.Name, Curriculum.SEMESTER_INDEX_NAME[indicies.Semester],
						Curriculum.SECTION[indicies.Section],
						subject_from_schedule,
					)
				} else if indicies.Section < 4 && len(subject_from_schedule) == 0 {
					t.Fatalf(
						"%s, %s, %s, section %s has 0 subjects even though it should contain at least one",
						values.Curriculum.CurriculumCode,
						values.YearLevel.Name, Curriculum.SEMESTER_INDEX_NAME[indicies.Semester],
						Curriculum.SECTION[indicies.Section],
					)
				}

				return GeneticAlgorithm.IterProceed
			},
		)
	}
}

func setup_router() *gin.Engine {

	switch os.Getenv("USE_DATABASE") {

	case "MongoDB":

		//////////////////////////////////////////////////////////////////////////
		//                         MongoDB Persistence
		//////////////////////////////////////////////////////////////////////////

		mongo_client := StorageResources.NewMongodbClient()
		defer StorageResources.CloseMongodbClient(mongo_client)

		RouteGlobals.ResourcesPersistence = &StorageResources.Persistence{
			ReaderService: &StorageResources.MongodbReader{
				Mongo: &StorageResources.MongoDB{Client: mongo_client},
			},

			WriterService: &StorageResources.MongodbWriter{
				Mongo: &StorageResources.MongoDB{Client: mongo_client},
			},
		}

		RouteGlobals.SchedulePersistence = &StorageSchedule.Persistence{
			LoadService: &StorageSchedule.MongodbReader{
				Mongo: &StorageSchedule.MongoDB{Client: mongo_client},
			},

			SaveService: &StorageSchedule.MongodbWriter{
				Mongo: &StorageSchedule.MongoDB{Client: mongo_client},
			},
		}

	default:

		//////////////////////////////////////////////////////////////////////////
		//                          JSON Persistence
		//////////////////////////////////////////////////////////////////////////

		RouteGlobals.ResourcesPersistence = &StorageResources.Persistence{
			ReaderService: &StorageResources.JsonReader{},
			WriterService: &StorageResources.JsonWriter{},
		}

		RouteGlobals.SchedulePersistence = &StorageSchedule.Persistence{
			LoadService: &StorageSchedule.JsonReader{},
			SaveService: &StorageSchedule.JsonWriter{},
		}
	}

	//////////////////////////////////////////////////////////////////////////
	//                         Initialize Globals
	//////////////////////////////////////////////////////////////////////////

	RouteGlobals.InitializeCachedUniversitySchedule()
	RouteGlobals.InitDeptSchedGenQueue()

	//////////////////////////////////////////////////////////////////////////

	var router *gin.Engine

	if gin.Mode() == gin.ReleaseMode {
		router = gin.Default()
	} else {
		gin.SetMode(gin.TestMode)
		router = gin.Default()
	}

	// maximum memory limit for multipart form file uploads
	router.MaxMultipartMemory = 5 << 20 // 5 MiB

	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	//////////////////////////////////////////////////////////////////////////
	//                              API-v1
	//////////////////////////////////////////////////////////////////////////

	v1 := router.Group("/v1")
	v2 := router.Group("/v2")

	v1.GET("/const", RoutesV1.GetConst)

	// ============= department routes and handlers =============

	v1.GET("/all_departments", RoutesV1.GetAllDepartments)
	v1.GET("/departments", RoutesV1.GetDepartmentsPaginated)
	v1.GET("/department_data", RoutesV1.GetCurriculumsDataInDepartment)
	v1.POST("/department_add", RoutesV1.PostDepartment)
	v1.PATCH("/department_update", RoutesV1.PatchDepartment)
	v1.DELETE("/department_remove", RoutesV1.DeleteDepartment)

	// ============= instructor routes and handlers =============

	v1.GET("/instructors/d", RoutesV1.GetDepartmentInstructorsDefaults)
	v1.GET("/instructors/a", RoutesV1.GetDepartmentInstructorsAllocated)
	v1.POST("/instructor_add", RoutesV1.PostInstructor)
	v1.PATCH("/instructor_update", RoutesV1.PatchInstructor)
	v1.DELETE("/instructor_remove", RoutesV1.DeleteInstructor)

	v2.GET("instructor_basic", RoutesV2.GetInstructorBasic)
	v2.GET("instructors", RoutesV2.GetDepartmentInstructors)
	v2.GET("instructor_resources", RoutesV2.GetInstructorResource)

	// ============= room routes and handlers =============

	v1.GET("/rooms", RoutesV1.GetDepartmentRooms)
	v1.POST("/room_add", RoutesV1.PostRoom)
	v1.PATCH("/room_update", RoutesV1.PatchRoom)
	v1.DELETE("/room_remove", RoutesV1.DeleteRoom)

	// ============= subject routes and handlers =============

	v1.GET("/subjects", RoutesV1.GetSubjects)
	v1.POST("/subject_add", RoutesV1.PostSubject)
	v1.PATCH("/subject_update", RoutesV1.PatchSubject)
	v1.DELETE("/subject_remove", RoutesV1.DeleteSubject)

	// ============= curriculum routes and handlers =============

	v1.GET("/curriculum_list", RoutesV1.GetDepartmentCurriculumList)
	v1.GET("/curriculum_load", RoutesV1.GetCurriculum)
	v1.POST("/curriculum_add", RoutesV1.PostCurriculum)
	v1.PATCH("/curriculum_update", RoutesV1.PatchCurriculum)
	v1.DELETE("/curriculum_remove", RoutesV1.DeleteCurriculum)

	// ============= schedule routes and handlers =============

	v1.GET("/gen_status", RoutesV1.GetGenStatus)
	v1.GET("/dept_gen_result", RoutesV1.GetDeptartmentGenerationResult)

	v1.GET("/university_schedule", RoutesV1.GetUniversitySchedule)
	v1.POST("/university_schedule", RoutesV1.PostUniversitySchedule)

	v1.GET("/class_schedule", RoutesV1.GetClassSchedule)
	v2.GET("/class_json_schedule", RoutesV2.GetJsonClassSchedule)
	v2.DELETE("/clear_class_schedule", RoutesV2.DeleteClearClassSchedule)
	v1.DELETE("/clear_department_schedules", RoutesV1.DeleteClearDepartmentSchedule)
	v1.POST("/generate_schedule", RoutesV1.RequestGenerateSchedule)

	v2.POST("/available_subject_moves", RoutesV2.GetSubjectAvailableTimeSlotMoves)
	v2.POST("/subject_move", RoutesV2.PostSubjectTimeSlotMove)

	v2.GET("/validate_schedules", RoutesV2.GetValidateSchedules)

	if os.Getenv("GIN_MODE") != "release" {
		v1.GET("/generate_schedule", RoutesV1.RequestGenerateSchedule) // for dev only
	}

	// ============= survery routes and handlers =============

	v2.POST("/add_schedule_preference", RoutesV2.PostWeekTimeTableSurvery)

	return router
}
