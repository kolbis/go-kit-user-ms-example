package middleware

import "github.com/thelotter-enterprise/usergo/svc"

// ServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type ServiceMiddleware func(svc.Service) svc.Service
