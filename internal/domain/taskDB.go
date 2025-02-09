package domain

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain/enums"
)

type Task struct {
	TaskId  string
	ChatId  string
	Day     enums.DayOfTheWeek
	Start   time.Time
	End     time.Time
	Comment string
}
