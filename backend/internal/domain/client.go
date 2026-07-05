package domain

import "context"

type Client struct {
	ID           int64  `db:"id" json:"id"`
	Login        string `db:"login" json:"login"`
	PasswordHash string `db:"password_hash" json:"-"`
	CreatedAt    string `db:"created_at" json:"createdAt"`
}

type ClientRepository interface {
	InsertClient(ctx context.Context, login, passwordHash string) (Client, error)
	GetClientByLogin(ctx context.Context, login string) (Client, error)
	GetClientByID(ctx context.Context, id int64) (Client, error)
}
