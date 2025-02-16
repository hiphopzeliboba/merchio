package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser_Success(t *testing.T) {
	// Инициализация моков
	mockRepo := new(MockUserRepository)
	s := &serv{userRepository: mockRepo}

	username := "testuser"
	password := "password123"

	// Мокаем создание пользователя в репозитории
	mockRepo.On("CreateUser", mock.Anything, username, mock.Anything).Return(int64(1), nil)

	// Вызов функции
	id, err := s.CreateUser(context.Background(), username, password)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	// Убедимся, что мок был вызван
	mockRepo.AssertExpectations(t)
}

//
//func TestCreateUser_HashPasswordError(t *testing.T) {
//	// Переопределим функцию hashPassword, чтобы симулировать ошибку
//	originalHashPassword := hashPassword
//	hashPassword = func(password string) (string, error) {
//		return "", errors.New("ошибка хеширования пароля")
//	}
//	defer func() { hashPassword = originalHashPassword }() // Восстановим оригинальную функцию после теста
//
//	// Инициализация моков
//	mockRepo := new(MockUserRepository)
//	s := &serv{userRepository: mockRepo}
//
//	username := "testuser"
//	password := "password123"
//
//	// Вызов функции
//	id, err := s.CreateUser(context.Background(), username, password)
//
//	// Проверки
//	assert.Error(t, err)
//	assert.Equal(t, int64(0), id)
//	assert.Equal(t, "ошибка хеширования пароля", err.Error())
//}

func TestCreateUser_CreateUserError(t *testing.T) {
	// Инициализация моков
	mockRepo := new(MockUserRepository)
	s := &serv{userRepository: mockRepo}

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Мокаем ошибку при создании пользователя
	mockRepo.On("CreateUser", mock.Anything, username, string(hashedPassword)).Return(int64(0), errors.New("ошибка создания пользователя"))

	// Вызов функции
	id, err := s.CreateUser(context.Background(), username, password)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, "ошибка создания пользователя", err.Error())

	// Убедимся, что мок был вызван
	mockRepo.AssertExpectations(t)
}
