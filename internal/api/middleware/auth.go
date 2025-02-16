package middleware

import (
	"context"
	"merchio/internal/utils"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			// тут должен быть редирект наверно
			http.Error(w, "отсутствует токен", http.StatusUnauthorized)
			return
		}

		// Удаляем префикс "Bearer " если он есть
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "недействительный токен", http.StatusUnauthorized)
			return
		}

		// Добавляем данные пользователя в контекст
		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		//ctx = context.WithValue(ctx, "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
