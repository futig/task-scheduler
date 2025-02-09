package domain

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain/enums"
)

type TaskRemind struct {
	ChatId  int64
	Day     enums.DayOfTheWeek
	Start   time.Time
	End     time.Time
	Comment string
	Type    enums.RemindType
}
