package user

import (
	"context"
	"errors"
	"merchio/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) IsUserPresent(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, username, password string) (int64, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

//
//func (m *MockUserRepository) IsUserPresent(ctx context.Context, username string) (bool, error) {
//	args := m.Called(ctx, username)
//	return args.Bool(0), args.Error(1)
//}
//
//func (m *MockUserRepository) CreateUser(ctx context.Context, username, password string) (int, error) {
//	args := m.Called(ctx, username, password)
//	return args.Int(0), args.Error(1)
//}
//
//func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
//	args := m.Called(ctx, username)
//	return args.Get(0).(*model.User), args.Error(1)
//}

type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) GenerateToken(userID int) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func TestAuth_Success(t *testing.T) {
	// Инициализация моков
	mockRepo := new(MockUserRepository)
	mockUtils := new(MockUtils)
	s := &serv{userRepository: mockRepo}

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Создаем моки
	mockRepo.On("IsUserPresent", mock.Anything, username).Return(true, nil)
	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(&model.User{ID: 1, Password: string(hashedPassword)}, nil)
	mockUtils.On("GenerateToken", 1).Return("valid-token", nil)

	// Вызов функции
	token, err := s.Auth(context.Background(), username, password)

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, "valid-token", token)

	// Убедимся, что моки были вызваны
	mockRepo.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestAuth_UserNotFound(t *testing.T) {
	// Инициализация моков
	mockRepo := new(MockUserRepository)
	//mockUtils := new(MockUtils)
	s := &serv{userRepository: mockRepo}

	username := "testuser"
	password := "password123"

	// Создаем моки
	mockRepo.On("IsUserPresent", mock.Anything, username).Return(false, nil)

	// Вызов функции
	token, err := s.Auth(context.Background(), username, password)

	// Проверки
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "ошибка существования пользователя", err.Error())

	// Убедимся, что моки были вызваны
	mockRepo.AssertExpectations(t)
}

func TestAuth_IncorrectPassword(t *testing.T) {
	// Инициализация моков
	mockRepo := new(MockUserRepository)
	//mockUtils := new(MockUtils)
	s := &serv{userRepository: mockRepo}

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("wrongpassword"), bcrypt.DefaultCost)

	// Создаем моки
	mockRepo.On("IsUserPresent", mock.Anything, username).Return(true, nil)
	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(&model.User{ID: 1, Password: string(hashedPassword)}, nil)

	// Вызов функции
	token, err := s.Auth(context.Background(), username, password)

	// Проверки
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "неверный пароль", err.Error())

	// Убедимся, что моки были вызваны
	mockRepo.AssertExpectations(t)
}

func TestAuth_GenerateTokenError(t *testing.T) {
	// Инициализация моков
	mockRepo := new(MockUserRepository)
	mockUtils := new(MockUtils)
	s := &serv{userRepository: mockRepo}

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Создаем моки
	mockRepo.On("IsUserPresent", mock.Anything, username).Return(true, nil)
	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(&model.User{ID: 1, Password: string(hashedPassword)}, nil)
	mockUtils.On("GenerateToken", 1).Return("", errors.New("ошибка генерации токена"))

	// Вызов функции
	token, err := s.Auth(context.Background(), username, password)

	// Проверки
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "ошибка генерации токена", err.Error())

	// Убедимся, что моки были вызваны
	mockRepo.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestAuth_DBError(t *testing.T) {
	// Инициализация моков
	mockRepo := new(MockUserRepository)
	//mockUtils := new(MockUtils)
	s := &serv{userRepository: mockRepo}

	username := "testuser"
	password := "password123"

	// Создаем моки
	mockRepo.On("IsUserPresent", mock.Anything, username).Return(false, errors.New("ошибка БД"))

	// Вызов функции
	token, err := s.Auth(context.Background(), username, password)

	// Проверки
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "ошибка существования пользователя", err.Error())

	// Убедимся, что моки были вызваны
	mockRepo.AssertExpectations(t)
}
