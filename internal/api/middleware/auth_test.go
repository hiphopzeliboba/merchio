package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/stretchr/testify/assert"
)

type Claims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// Мок для utils.VerifyToken
func mockVerifyToken(token string) (*Claims, error) {
	if token == "validtoken" {
		return &Claims{UserId: 12345}, nil
	}
	return nil, errors.New("invalid token")
}

func TestAuthMiddleware(t *testing.T) {

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No token provided",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "отсутствует токен",
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "недействительный токен",
		},
		{
			name:           "Valid token",
			authHeader:     "Bearer validtoken",
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Хендлер для проверки контекста
			handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
				userID := r.Context().Value("user_id")
				assert.Equal(t, "12345", userID)
				w.Write([]byte("ok"))
			})

			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
