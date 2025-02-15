package service

import (
	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/storage"
)

func GetUserState(cache storage.Cache, storage storage.Storage, chatID int64) (domain.User, bool, error) {
	if state, ok := cache.GetValue(chatID); ok {
		return state, true, nil;
	}
	return storage.GetUserState(chatID)
}

func SetUserState(cache storage.Cache,  storage storage.Storage, chatID int64, user domain.User) error {
	_ = cache.SetValue(chatID, user)

	err := storage.UpdateUserState(chatID, user)
	if err != nil {
		return storage.CreateUserState(user)
	}
	return nil
}