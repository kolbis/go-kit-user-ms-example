package amqp

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"

	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/thelotter-enterprise/corego/errors"
	tlelogger "github.com/thelotter-enterprise/corego/logger"
	tlerabbitmq "github.com/thelotter-enterprise/corego/transport/rabbitmq"
	"github.com/thelotter-enterprise/corego/utils"
	"github.com/thelotter-enterprise/usergo/shared"
	"github.com/thelotter-enterprise/usergo/svc/transport"
)

// NewTransport will create all the rabbitMQ consumers information
// it will not run them.
func NewTransport(svcEndpoints transport.Endpoints, logger tlelogger.Logger, connMgr *tlerabbitmq.ConnectionManager) *[]tlerabbitmq.Subscriber {
	subscribers := make([]tlerabbitmq.Subscriber, 0)
	subMgr := tlerabbitmq.NewSubscriberManager(connMgr)

	// Important: each queue can only get a single type of message!!
	// This is due to the opinionated nature of gokit
	loggedInSubscriber := subMgr.NewCommandSubscriber(
		"UserLoggedIn",
		"UserLoggedIn",
		svcEndpoints.UserLoggedInConsumerEndpoint,
		decodeLoggedInUserCommand,
		amqptransport.EncodeJSONResponse,
	)

	// here we can have additional private subscribers
	subscribers = append(subscribers, loggedInSubscriber)
	return &subscribers
}

func decodeLoggedInUserCommand(_ context.Context, msg *amqp.Delivery) (interface{}, error) {
	var data shared.LoggedInCommandData
	decoder := utils.NewDecoder()

	m := tlerabbitmq.Message{
		Payload: &tlerabbitmq.MessagePayload{},
	}
	err := json.Unmarshal(msg.Body, &m)
	if err != nil {
		return m, errors.NewApplicationError(err, "failed to decode loggedInUserCommand")
	}
	err = decoder.MapDecode(m.Payload.Data, &data)
	m.Payload.Data = data

	return m, err
}
