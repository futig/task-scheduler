package tempstorage

import (
	"fmt"

	"github.com/futig/task-scheduler/internal/domain"
)

func (s *TempStorageContext) CreateUserState(userState domain.User) error {
	if _, ok := s.Users.Load(userState.ChatID); ok {
		return fmt.Errorf("пользователь с таким id уже существует %d", userState.ChatID)
	}
	s.Users.Store(userState.ChatID, userState)
	return nil
}

func (s *TempStorageContext) GetUserState(chatId int64) (domain.User, bool, error) {
	if val, ok := s.Users.Load(chatId); ok {
		userState := val.(domain.User)
		return userState, ok, nil
	}
	return domain.User{}, false, fmt.Errorf("пользователь с таким id не существует: %d", chatId)
}

func (s *TempStorageContext) UpdateUserState(chatId int64, userState domain.User) error {
	if _, ok := s.Users.Load(chatId); ok {
		s.Users.Store(chatId, userState)
		return nil
	}
	return fmt.Errorf("пользователь с таким id не существует: %d", chatId)
}

func (s *TempStorageContext) DeleteUserState(chatId int64) error {
	s.Reminds.Delete(chatId)
	return nil
}