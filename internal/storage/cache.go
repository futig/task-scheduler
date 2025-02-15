package storage

import (
	"github.com/futig/task-scheduler/internal/domain"
)

type Cache interface {
	SetValue(key int64, user domain.User) error
	GetValue(key int64) (domain.User, bool)
}
