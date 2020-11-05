package svc

import (
	"context"

	"github.com/thelotter-enterprise/usergo/shared"
)

// Service is the API exposed by the microservice
type Service interface {

	// GetUserByID will return the user based on the inpit ID
	GetUserByID(ctx context.Context, userID int) (shared.User, error)

	// ConsumeLoginCommand will consume a rabbit message when a user is logged in
	ConsumeLoginCommand(ctx context.Context, userID int) error
}
