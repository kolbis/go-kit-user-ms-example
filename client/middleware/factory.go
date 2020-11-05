package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	httptransport "github.com/go-kit/kit/transport/http"
	tlecb "github.com/thelotter-enterprise/corego/circuitbreaker"
	tlectxhttp "github.com/thelotter-enterprise/corego/context/transport/http"
	tlefallback "github.com/thelotter-enterprise/corego/fallback"
	tlelimiter "github.com/thelotter-enterprise/corego/ratelimit"
	"github.com/thelotter-enterprise/corego/utils"
	"github.com/thelotter-enterprise/usergo/shared"
)

// GetUserByIDFactory will return sd.Factory which will be used to constuct the transport endpoint
func GetUserByIDFactory(ctx context.Context, id int) sd.Factory {
	// TODO: is this the way to handle path replacement?
	path := fmt.Sprintf(shared.UserByIDClientRoute, id)

	// Instance will look something like this "http://localhost:8080"
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		// instance = "http://localhost:8080" // i keep it for testing
		breakermw := tlecb.NewDefaultHystrixCommandMiddleware("get_user_by_id")
		limitermw := tlelimiter.NewDefaultErrorLimitterMiddleware()
		fallbackmw := tlefallback.NewFallbackMiddleware(getDefaultUser())

		tgt, _ := url.Parse(instance)
		tgt.Path = path

		endpoint := httptransport.NewClient(
			"GET",
			tgt,
			encodeGetUserByIDRequest,
			decodeGetUserByIDResponse,
			tlectxhttp.WriteBefore()).Endpoint()

		endpoint = breakermw(endpoint)
		endpoint = limitermw(endpoint)
		// fallback should run last. if circuit was opened, it will return the response from the fallback
		endpoint = fallbackmw(endpoint)

		return endpoint, nil, nil
	}
}

func encodeGetUserByIDRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := shared.NewByIDRequest(ctx, request.(int))
	enc := utils.EncodeRequestToJSON(ctx, r, req)
	return enc
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponseData
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func getDefaultUser() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := shared.ByIDResponseData{
			User: shared.User{
				ID: 1,
			},
		}
		return resp, nil
	}
}
