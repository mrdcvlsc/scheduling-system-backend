package RouteGlobals

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"sync"
)

var rw_map_mutex sync.RWMutex
var department_id_to_sched_gen_last_result map[DeptSchedGenKey]ScheduleGenerationLastResult

type DeptSchedGenKey struct {
	DepartmentID uint16
	Semester     int
}

type ScheduleGenerationLastResult struct {
	Status  bool
	Message string
}

func SetDepartmentsLastScheduleGenerationResult(key DeptSchedGenKey, value ScheduleGenerationLastResult) {
	rw_map_mutex.Lock()
	defer rw_map_mutex.Unlock()

	if department_id_to_sched_gen_last_result == nil {
		department_id_to_sched_gen_last_result = make(map[DeptSchedGenKey]ScheduleGenerationLastResult)
	}

	department_id_to_sched_gen_last_result[key] = value
}

func GetDepartmentsLastScheduleGenerationResult(key DeptSchedGenKey) ScheduleGenerationLastResult {
	rw_map_mutex.RLock()
	defer rw_map_mutex.RUnlock()

	last_sched_gen_result, has_key := department_id_to_sched_gen_last_result[key]

	if !has_key {
		return ScheduleGenerationLastResult{
			Status:  false,
			Message: fmt.Sprintf("schedule is not generated yet for the department with id %d", key),
		}
	}

	return last_sched_gen_result
}

var rw_queue_mutex sync.RWMutex
var departments_schedule_generation_queue []DeptSchedGenKey

// initialize department generation request count.
func InitDeptSchedGenQueue() {
	if departments_schedule_generation_queue == nil {
		log.Print("department to encode requests initialized")
		departments_schedule_generation_queue = make([]DeptSchedGenKey, 0)
	} else {
		log.Print("department to encode requests is already initialized")
	}
}

// return true if the department was added to the schedule generation task queue.
//
// return false if the department was already in the schedule generation task queue.
func PushNewDeptToDeptSchedGenQueue(department_id_and_semester DeptSchedGenKey) bool {
	rw_queue_mutex.RLock()
	defer rw_queue_mutex.RUnlock()

	if slices.Contains(departments_schedule_generation_queue, department_id_and_semester) {
		return false
	}

	departments_schedule_generation_queue = append(departments_schedule_generation_queue, department_id_and_semester)
	return true
}

// returns only one department id mapped to a true boolean value.
//
// returns an error if the queue is empty or uninitialized.
func PopDepartmentToEncodeFromSchedGenQueue() (map[uint16]bool, int, error) {
	rw_queue_mutex.Lock()
	defer rw_queue_mutex.Unlock()

	if departments_schedule_generation_queue == nil {
		return nil, -1, errors.New("the current schedule generation queue is uninitialized")
	}

	if len(departments_schedule_generation_queue) == 0 {
		return nil, -1, errors.New("the current schedule generation queue is empty")
	}

	first_dept_sem := departments_schedule_generation_queue[0]

	if len(departments_schedule_generation_queue) > 1 {
		departments_schedule_generation_queue = departments_schedule_generation_queue[1:]
	} else {
		departments_schedule_generation_queue = make([]DeptSchedGenKey, 0)
	}

	department_to_encode := make(map[uint16]bool)
	department_to_encode[first_dept_sem.DepartmentID] = true

	return department_to_encode, first_dept_sem.Semester, nil
}
