package service

import "context"

type AuthorizeService interface {
	RegisterUser(ctx context.Context, id, email, password string) error
	Authorize(ctx context.Context, email, password string) error
}
