package domain

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain/enums"
)

type Remind struct {
	TimingId string
	TaskId   string
	Type     enums.RemindType
	Time     time.Time
}
