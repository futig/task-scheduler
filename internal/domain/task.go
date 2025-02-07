package domain

import "time"

type Task struct {
	TaskId string
	UserId string
	Day DayOfTheWeek
	Start time.Time
	End time.Time
}
