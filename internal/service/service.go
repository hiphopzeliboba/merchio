package service

import (
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, username, password string) (int64, error)
	Auth(ctx context.Context, username, password string) (string, error)
}
