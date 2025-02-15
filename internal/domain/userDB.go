package domain

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain/enums"
	"github.com/google/uuid"
)

type User struct {
	ChatID         int64
	State          enums.UserState
	ChoosenWeekday time.Weekday
	ChoosenItem    uuid.UUID
	LastMessageId  int
}
