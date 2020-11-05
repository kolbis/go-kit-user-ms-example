package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	tlehttp "github.com/thelotter-enterprise/corego/transport/http"
	tlerabbitmq "github.com/thelotter-enterprise/corego/transport/rabbitmq"
	"github.com/thelotter-enterprise/corego/utils"
	"github.com/thelotter-enterprise/usergo/shared"
	"github.com/thelotter-enterprise/usergo/svc"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	UserByIDEndpoint             endpoint.Endpoint
	UserLoggedInConsumerEndpoint endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s svc.Service) Endpoints {
	return Endpoints{
		UserByIDEndpoint:             makeUserByIDEndpoint(s),
		UserLoggedInConsumerEndpoint: makeUserLoggedInConsumerEndpoint(s),
	}
}

func makeUserByIDEndpoint(service svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		var data shared.ByIDRequestData

		// we either getting a built in reuest or a map. We need to defrenciate between cases.
		// a concrete request is send when there is no body and we read from uri
		if req, ok := request.(tlehttp.Request); ok == false {
			decoder := utils.NewDecoder()
			err = decoder.MapDecode(request, &req)
			err = decoder.MapDecode(req.Data, &data)
			req.Data = data
		} else {
			data = req.Data.(shared.ByIDRequestData)
		}

		user, err := service.GetUserByID(ctx, data.ID)
		return shared.NewUserResponse(user), err
	}
}

func makeUserLoggedInConsumerEndpoint(service svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		message := request.(tlerabbitmq.Message)
		data := message.Payload.Data.(shared.LoggedInCommandData)
		err := service.ConsumeLoginCommand(ctx, data.ID)
		return true, err
	}
}
