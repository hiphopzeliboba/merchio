package user

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

func (s *serv) CreateUser(ctx context.Context, username, password string) (int64, error) {

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return 0, err
	}
	id, err := s.userRepository.CreateUser(ctx, username, hashedPassword)
	return id, err
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
