package implementation

import (
	"context"
	"fmt"

	tlectx "github.com/kolbis/corego/context"
	tletracer "github.com/kolbis/corego/tracer"
	"github.com/kolbis/go-kit-user-ms-example/shared"
	"github.com/kolbis/go-kit-user-ms-example/svc"

	tlelogger "github.com/kolbis/corego/logger"
)

type service struct {
	repo   svc.Repository
	tracer tletracer.Tracer
	logger tlelogger.Logger
}

// NewService creates a new instance of service
// service is where we define all the business logic.
func NewService(logger tlelogger.Logger, tracer tletracer.Tracer, repo svc.Repository) svc.Service {
	return &service{
		repo:   repo,
		tracer: tracer,
		logger: logger,
	}
}

// GetUserByID will execute business logic for getting user information by id
func (s *service) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	return user, err
}

func (s *service) ConsumeLoginCommand(ctx context.Context, userID int) error {
	corrid := tlectx.GetCorrelation(ctx)
	duration, deadline := tlectx.GetTimeout(ctx)
	fmt.Printf("consumed LoggedInCommand, user: %d, correation %s, duration %s, deadline %s", userID, corrid, duration, deadline)
	return nil
}
