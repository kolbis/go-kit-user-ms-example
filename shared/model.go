package shared

import (
	"context"

	tlehttp "github.com/thelotter-enterprise/corego/transport/http"
)

// User ...
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ByIDRequestData ...
type ByIDRequestData struct {
	ID int `json:"id"`
}

// ByIDResponseData ...
type ByIDResponseData struct {
	User User
}

// NewByIDRequest will create a Request with ByIDRequestData
func NewByIDRequest(ctx context.Context, id int) tlehttp.Request {
	data := ByIDRequestData{
		ID: id,
	}
	req := tlehttp.Request{}.Wrap(ctx, data)
	return req
}

// NewUserResponse ...
func NewUserResponse(user User) ByIDResponseData {
	return ByIDResponseData{User: user}
}

// LoggedInCommandData ...
type LoggedInCommandData struct {
	ID   int
	Name string
}
