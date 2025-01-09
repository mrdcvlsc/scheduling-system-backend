package storage_test

import (
	"testing"

	p "github.com/mrdcvlsc/scheduling-system-backend/Storage"
	"github.com/mrdcvlsc/scheduling-system-backend/Utils"
)

func TestJsonFilePersistence(t *testing.T) {
	json_persistence := p.JsonFilePersistence{}
	curriculums := json_persistence.GetAllCurriculum()
	Utils.PrettyPrint(curriculums)
}
