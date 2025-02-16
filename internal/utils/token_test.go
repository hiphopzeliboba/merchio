package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		shouldError bool
	}{
		{
			name:        "Успешная генерация и проверка токена",
			userID:      1,
			shouldError: false,
		},
		{
			name:        "Генерация токена с нулевым ID",
			userID:      0,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Генерация токена
			token, err := GenerateToken(tt.userID)
			if tt.shouldError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Проверка токена
			claims, err := VerifyToken(token)
			assert.NoError(t, err)
			assert.NotNil(t, claims)
			assert.Equal(t, tt.userID, claims.UserId)
		})
	}
}

func TestVerifyToken_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		shouldError bool
	}{
		{
			name:        "Неверный формат токена",
			token:       "invalid.token.format",
			shouldError: true,
		},
		{
			name:        "Пустой токен",
			token:       "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := VerifyToken(tt.token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}
