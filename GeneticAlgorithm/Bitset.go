package geneticalgorithm

import (
	"fmt"
)

// ====================================================================================
// convert an unsigned integer into 1 if it has any bit that is set to 1,
// or into 0 if all bits are set to 0, using only bitwise operations
// to compute value during compile time.

const BITSET_LIMB_WIDENESS = 64
const u8BitReduce uint8 = WEEKLY_TIME_SLOTS % BITSET_LIMB_WIDENESS

// since modding by 64 will always produce results less than 64 which
// can fit in 8 bit wide unsigned integer, we can then right away start
// reducing at the first 8 bits instead from from all 64 bits.

const u4BitReduce = (u8BitReduce >> 4) | (u8BitReduce & uint8(0b1111))
const u2BitReduce = (u4BitReduce >> 2) | (u4BitReduce & uint8(0b11))
const u1BitReduce = (u2BitReduce >> 1) | (u2BitReduce & uint8(0b1))

// ====================================================================================

const BITSET_LIMBS = (WEEKLY_TIME_SLOTS / 64) + u1BitReduce

type InstructorWeekTimeSlotAvailabilityMap [BITSET_LIMBS]uint64

// setting a bit to 1 means the instructor is available,
// and 0 if not for that corresponding bit time slot.
func (bitset *InstructorWeekTimeSlotAvailabilityMap) SetAvailability(available bool, day, time_slot int) {
	if day < 0 || day >= SCHOOL_DAYS_PER_WEEK {
		panic(fmt.Sprintf(
			"SetAvailability(..., day int = %d,...) : invalid argument `day`, accepted values are only 0 to %d",
			day, SCHOOL_DAYS_PER_WEEK-1,
		))
	}

	if time_slot < 0 || time_slot >= DAILY_TIME_SLOTS {
		panic(fmt.Sprintf(
			"SetAvailability(..., time_slot int = %d) : invalid argument `time_slot`, accepted values are only 0 to %d",
			time_slot, DAILY_TIME_SLOTS-1,
		))
	}

	bit_idx := day*DAILY_TIME_SLOTS + time_slot
	limb_idx := bit_idx / BITSET_LIMB_WIDENESS
	limb_bit_idx := bit_idx % BITSET_LIMB_WIDENESS

	if available {
		bitset[limb_idx] &= ^(uint64(1) << limb_bit_idx)
	} else {
		bitset[limb_idx] |= (uint64(1) << limb_bit_idx)
	}
}

func (bitset *InstructorWeekTimeSlotAvailabilityMap) GetAvailability(day, time_slot int) bool {
	if day < 0 || day >= SCHOOL_DAYS_PER_WEEK {
		panic(fmt.Sprintf(
			"GetAvailability(..., day int = %d,...) : invalid argument `day`, accepted values are only 0 to %d",
			day, SCHOOL_DAYS_PER_WEEK-1,
		))
	}

	if time_slot < 0 || time_slot >= DAILY_TIME_SLOTS {
		panic(fmt.Sprintf(
			"GetAvailability(..., time_slot int = %d) : invalid argument `time_slot`, accepted values are only 0 to %d",
			time_slot, DAILY_TIME_SLOTS-1,
		))
	}

	bit_idx := day*DAILY_TIME_SLOTS + time_slot
	limb_idx := bit_idx / BITSET_LIMB_WIDENESS
	limb_bit_idx := bit_idx % BITSET_LIMB_WIDENESS

	return ((bitset[limb_idx] >> limb_bit_idx) & uint64(1)) == 0
}
