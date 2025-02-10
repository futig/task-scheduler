package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id      uuid.UUID
	ChatId  int64
	Day     time.Weekday
	Start   int // В минутах с начала дня
	End     int // В минутах с начала дня
	Comment string
}
