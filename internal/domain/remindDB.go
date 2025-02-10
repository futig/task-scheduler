package domain

import (
	"github.com/futig/task-scheduler/internal/domain/enums"
	"github.com/google/uuid"
)

type Remind struct {
	Id     uuid.UUID
	TaskId uuid.UUID
	Type   enums.RemindType
	Time   int
}
