package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	tlectxhttp "github.com/kolbis/corego/context/transport/http"
	tlelogger "github.com/kolbis/corego/logger"
	tlehttp "github.com/kolbis/corego/transport/http"
	"github.com/kolbis/go-kit-user-ms-example/shared"
	"github.com/kolbis/go-kit-user-ms-example/svc/transport"
)

// NewTransport will set-up router and initialize http endpoints
func NewTransport(ctx context.Context, svcEndpoints transport.Endpoints, options []kithttp.ServerOption, logger tlelogger.Logger) http.Handler {
	var (
		router = mux.NewRouter()

		// server options:
		errorLogger   = kithttp.ServerErrorLogger(logger)
		errorEncoder  = kithttp.ServerErrorEncoder(encodeErrorResponse)
		contextReader = tlectxhttp.ReadBefore()
	)

	options = append(options, errorLogger, errorEncoder, contextReader)

	getUserByIDHandler := kithttp.NewServer(
		svcEndpoints.UserByIDEndpoint,
		decodeUserByIDRequest,
		encodeUserByIDReponse,
		options...)

	router.Methods("GET").Path(shared.UserByIDServerRoute).Handler(getUserByIDHandler)

	return handlers.LoggingHandler(os.Stdout, router)
}

func encodeUserByIDReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	// EOF means that the request body is empty
	if err == io.EOF {
		// In that case we will manually construct the object by reading the URI
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		data := shared.ByIDRequestData{
			ID: id,
		}
		req = tlehttp.Request{
			Data: data,
		}
		// resetting the error
		err = nil
	}

	return req, err
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}
