package implementation

import (
	"context"

	"github.com/kolbis/go-kit-user-ms-example/shared"
	"github.com/kolbis/go-kit-user-ms-example/svc"
)

type repo struct {
	// database for example
}

// NewRepository ...
func NewRepository() svc.Repository {
	return &repo{}
}

// GetUserByID ...
func (r repo) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	user := shared.User{
		ID:        userID,
		Email:     "guyk@net-bet.net",
		FirstName: "guy",
		LastName:  "kolbis",
	}

	return user, nil
}
