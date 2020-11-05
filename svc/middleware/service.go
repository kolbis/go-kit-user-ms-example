package middleware

import "github.com/kolbis/go-kit-user-ms-example/svc"

// ServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type ServiceMiddleware func(svc.Service) svc.Service
