package geneticalgorithm

// specifies the number of days per week the university holds classes.
const SCHOOL_DAYS_PER_WEEK = 6

// specifies the total number of teaching hours the university is open per day.
const DAILY_TEACHING_HOURS = 12 //

// indicates the number of time slots available within a single teaching hour.
const TIME_SLOTS_PER_HOUR = 2

// represents the total number of time slots available in a single teaching day.
const DAILY_TIME_SLOTS = DAILY_TEACHING_HOURS * TIME_SLOTS_PER_HOUR

// calculates the total number of time slots available for scheduling in a week.
const WEEKLY_TIME_SLOTS = SCHOOL_DAYS_PER_WEEK * DAILY_TIME_SLOTS
