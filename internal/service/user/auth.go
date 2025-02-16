package user

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"merchio/internal/utils"
)

func (s *serv) Auth(ctx context.Context, username, password string) (string, error) {
	// Получаем пользователя из БД
	present, err := s.userRepository.IsUserPresent(ctx, username)
	if err != nil {
		return "", errors.New("ошибка существования пользователя")
	}
	if !present {
		_, err = s.userRepository.CreateUser(ctx, username, password)
	}
	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("неверный пароль")
	}

	// Генерируем JWT токен
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
