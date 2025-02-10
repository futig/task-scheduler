package domain

import (
	"github.com/futig/task-scheduler/internal/domain/enums"
)

type TaskRemind struct {
	ChatId  int64
	Start   int // В минутах с начала дня
	End     int // В минутах с начала дня
	Comment string
	Type    enums.RemindType
}
