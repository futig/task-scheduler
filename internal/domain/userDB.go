package domain

import (
	"github.com/futig/task-scheduler/internal/domain/enums"
)

type User struct {
	ChatID int64
	State  enums.UserState
}
