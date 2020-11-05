package svc

import (
	"context"

	"github.com/thelotter-enterprise/usergo/shared"
)

// Repository ..
type Repository interface {
	GetUserByID(ctx context.Context, userID int) (shared.User, error)
}
