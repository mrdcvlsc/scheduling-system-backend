package subjecttest_test

import (
	"testing"

	ga "github.com/mrdcvlsc/scheduling-system-backend/GeneticAlgorithm"
)

func TestSubjectDbRecord1(t *testing.T) {
	subject := ga.SubjectDbRecord{}

	//////////////

	subject.SetLecHours(5)
	if subject.GetLecHours() != 5 {
		t.Error("subject.SetLecHours(5) - failed")
	}

	if subject.GetLabHours() != 0 {
		t.Error("subject.SetLecHours(5) + subject.GetLabHours() != 0 - failed")
	}

	//////////////

	subject.SetLecHours(0)
	if subject.GetLecHours() != 0 {
		t.Errorf("subject.SetLecHours(0) - failed = %d", subject.GetLecHours())
	}

	if subject.GetLabHours() != 0 {
		t.Error("subject.SetLecHours(0) + subject.GetLabHours() != 0 - failed")
	}

	//////////////

	subject.SetLabHours(5)
	if subject.GetLabHours() != 5 {
		t.Errorf("subject.SetLabHours(5) - failed = %d", subject.GetLabHours())
	}

	if subject.GetLecHours() != 0 {
		t.Error("subject.SetLabHours(5) + subject.GetLecHours() != 0 - failed")
	}

	//////////////

	subject.SetLabHours(0)
	if subject.GetLecHours() != 0 {
		t.Error("subject.SetLabHours(0) - failed")
	}

	if subject.GetLabHours() != 0 {
		t.Error("subject.SetLabHours(0) + subject.GetLabHours() != 0 - failed")
	}

	//////////////

	if subject.IsGymType() {
		t.Error("subject.IsGymType() - failed - this should not be a gym type yet")
	}
}

func TestSubjectDbRecord(t *testing.T) {
	subject := ga.SubjectDbRecord{}

	for i := uint8(1); i < uint8(ga.MAX_SUBJECT_HOUR_DURATION-1); i++ {
		subject.SetGymHours(i)

		if subject.GetLecHours() != subject.GetLabHours() {
			t.Errorf("subject.SetGymHours(%d) - test 1 failed", i)
		}

		if subject.GetLecHours() != i {
			t.Errorf("subject.SetGymHours(%d) - test 2 failed", i)
		}

		subject.SetLecHours(i + 1)

		if subject.GetLecHours() == subject.GetLabHours() {
			t.Errorf("subject.SetLecHours(%d) - test 1 failed", i+1)
		}

		if subject.GetLecHours() != (i + 1) {
			t.Errorf("subject.SetLecHours(%d) - test 2 failed", i+1)
		}
	}

	subject.SetLecHours(6)

	if subject.GetLecHours() != 6 {
		t.Error("subject.SetLecHours(6) - failed")
	}

	subject.SetLabHours(7)

	if subject.GetLabHours() != 7 {
		t.Error("subject.SetLabHours(7) - failed")
	}

	if subject.GetLecHours() != 6 {
		t.Errorf("subject.SetLecHours(6) = (%d) - failed 2nd ", subject.GetLecHours())
	}

	subject.SetGymHours(1)

	if subject.GetGymHours() != 1 {
		t.Error("subject.GetGymHours() != 0 : failed")
	}
}
