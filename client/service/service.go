package service

import (
	"context"

	tlehttp "github.com/thelotter-enterprise/corego/transport/http"
)

// ServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type ServiceMiddleware func(Service) Service

// Service defines all the APIs available for the service
type Service interface {
	// Gets the user by an ID
	GetUserByID(ctx context.Context, id int) tlehttp.Response

	// Gets the user by email
	GetUserByEmail(ctx context.Context, email string) tlehttp.Response
}
