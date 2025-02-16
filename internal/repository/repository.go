package repository

import (
	"context"
	"merchio/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, password string) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	IsUserPresent(ctx context.Context, username string) (bool, error)
}
