package tempcache

import (
	"sync"

	"github.com/futig/task-scheduler/internal/domain/enums"
)

type TempCache struct {
	Storage sync.Map
}

func (s *TempCache) GetValue(key string) (enums.UserState, bool) {
	if val, ok := s.Storage.Load(key); ok {
		strVal, ok := val.(enums.UserState)
		return strVal, ok
	}
	return 0, false
}

func (s *TempCache) SetValue(key string, val enums.UserState) error {
	s.Storage.Store(key, val)
	return nil
}
