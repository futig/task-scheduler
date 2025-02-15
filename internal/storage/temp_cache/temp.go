package tempcache

import (
	"sync"

	"github.com/futig/task-scheduler/internal/domain"
)

type TempCache struct {
	Storage sync.Map
}

func (s *TempCache) GetValue(key int64) (domain.User, bool) {
	if val, ok := s.Storage.Load(key); ok {
		strVal, ok := val.(domain.User)
		return strVal, ok
	}
	return domain.User{}, false
}

func (s *TempCache) SetValue(key int64, val domain.User) error {
	s.Storage.Store(key, val)
	return nil
}
