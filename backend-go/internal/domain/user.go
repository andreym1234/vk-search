package domain

import "context"

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	RoleID       int
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type AuthUseCase interface {
	Login(ctx context.Context, username, password string) (string, error)
}
