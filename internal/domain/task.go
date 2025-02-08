package domain

import "time"

type Task struct {
	TaskId  string
	ChatId  string
	Day     DayOfTheWeek
	Start   time.Time
	End     time.Time
	Comment string
}
