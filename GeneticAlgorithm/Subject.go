package geneticalgorithm

const CLASS_TYPE_LEC_HOURS_MASK uint8 = 0b00000111
const CLASS_TYPE_LAB_HOURS_MASK uint8 = 0b00111000
const CLASS_TYPE_GYM_BIT_F_MASK uint8 = 0b10000000

type SubjectDbRecord struct {
	// this should never be zero, zero means empty, none or nothing.
	SubjectID uint16

	// Determines the types of classes [lab, lec, gym] available and the number
	// of hours for each types.
	//
	// This attribute can only store up to 7 hours in each class types, since
	// hours of each class types of a subject is only stored in 3 bits.
	//
	// Lecture Hours - stored in the first least significant 3 bits.
	//
	// Lab Hours - stored in the following 3 bits after the previous least significant 3 bits.
	//
	// Gym Hours - This is a special case for FITT/PE subject, where both
	// the least significant 3 bits and the following 3 bits are populated
	// by the same hour number value and the most significant bit is 1.
	//
	// BIT FORMAT: [1-bit gym flag][1-bit unused][3-bit lab hours][3-bit lec hours]
	lec_lab_gym_hours uint8

	SubjectName string
	SubjectCode string
}

func (s *SubjectDbRecord) GetLecHours() uint8 {
	return s.lec_lab_gym_hours & CLASS_TYPE_LEC_HOURS_MASK
}

func (s *SubjectDbRecord) GetLabHours() uint8 {
	return (s.lec_lab_gym_hours & CLASS_TYPE_LAB_HOURS_MASK) >> 3
}

func (s *SubjectDbRecord) GetGymHours() uint8 {
	gym_bit_flag := (s.lec_lab_gym_hours & CLASS_TYPE_GYM_BIT_F_MASK) >> 7
	lecHours := s.GetLecHours()

	if lecHours == s.GetLabHours() && gym_bit_flag == 1 {
		return lecHours
	}

	return 0
}

func (s *SubjectDbRecord) SetLecHours(hours uint8) {
	s.lec_lab_gym_hours = ^CLASS_TYPE_LEC_HOURS_MASK | hours
}

func (s *SubjectDbRecord) SetLabHours(hours uint8) {
	s.lec_lab_gym_hours = ^CLASS_TYPE_LAB_HOURS_MASK | (hours << 3)
}

func (s *SubjectDbRecord) SetGymHours(hours uint8) {
	s.lec_lab_gym_hours = CLASS_TYPE_GYM_BIT_F_MASK | (hours << 3) | hours
}

func (s *SubjectDbRecord) IsGymType() bool {
	return ((s.lec_lab_gym_hours & CLASS_TYPE_GYM_BIT_F_MASK) >> 7) == 1
}
