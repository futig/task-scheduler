package domain

import (
	"fmt"
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

func (t *Task) String() string {
	startH := t.Start / 60
	startM := t.Start % 60
	endH := t.End / 60
	endM := t.End % 60
	return fmt.Sprintf("%d:%d-%d:%d ~ %s", startH, startM, endH, endM, t.Comment)
}