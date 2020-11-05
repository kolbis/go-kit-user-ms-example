package svc

import (
	"context"

	"github.com/kolbis/go-kit-user-ms-example/shared"
)

// Repository ..
type Repository interface {
	GetUserByID(ctx context.Context, userID int) (shared.User, error)
}
