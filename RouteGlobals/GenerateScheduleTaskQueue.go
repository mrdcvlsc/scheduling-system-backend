package RouteGlobals

import (
	"log"
	"sync"
)

var rw_mutex sync.RWMutex
var departments_schedule_generation_queue []uint16

// initialize department generation request count.
func InitDeptSchedGenQueue() {
	if departments_schedule_generation_queue == nil {
		log.Print("department to encode requests initialized")
		departments_schedule_generation_queue = make([]uint16, 0)
	} else {
		log.Print("department to encode requests is already initialized")
	}
}

// return true if the department was added to the schedule generation task queue.
//
// return false if the department was already in the schedule generation task queue.
func PushNewDeptToDeptSchedGenQueue(department_id uint16) bool {
	rw_mutex.RLock()
	defer rw_mutex.RUnlock()

	for _, dept_id := range departments_schedule_generation_queue {
		if dept_id == department_id {
			return false
		}
	}

	departments_schedule_generation_queue = append(departments_schedule_generation_queue, department_id)
	return true
}

func PopDepartmentToEncodeFromSchedGenQueue() map[uint16]bool {
	rw_mutex.Lock()
	defer rw_mutex.Unlock()

	if departments_schedule_generation_queue == nil {
		return nil
	}

	if len(departments_schedule_generation_queue) == 0 {
		return nil
	}

	department_id := departments_schedule_generation_queue[0]

	if len(departments_schedule_generation_queue) > 1 {
		departments_schedule_generation_queue = departments_schedule_generation_queue[1:]
	} else {
		departments_schedule_generation_queue = make([]uint16, 0)
	}

	department_to_encode := make(map[uint16]bool)
	department_to_encode[department_id] = true

	return department_to_encode
}
