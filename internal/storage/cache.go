package storage

import "github.com/futig/task-scheduler/internal/domain/enums"

type Cache interface {
	SetValue(key string, val enums.UserState) error
	GetValue(key string) (enums.UserState, bool)
}
