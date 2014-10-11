package routing

import (
	"encoding/json"
	"time"

	"code.google.com/p/go.net/context"
	pbuf "code.google.com/p/gogoprotobuf/proto"

	"github.com/opentarock/service-api/go/proto"
	"github.com/opentarock/service-api/go/proto_errors"
	"github.com/opentarock/service-api/go/proto_msghandler"
	"github.com/opentarock/service-api/go/reqcontext"
	"github.com/opentarock/service-api/go/service"
	"github.com/opentarock/service-msghandler/messages"
)

const defaultRequestTimeout = 1 * time.Minute

func NewRouteMessageHandler() service.MessageHandler {
	return service.MessageHandlerFunc(func(msg *proto.Message) proto.CompositeMessage {
		ctx, cancel := reqcontext.WithRequest(context.Background(), msg, defaultRequestTimeout)
		defer cancel()

		logger := reqcontext.ContextLogger(ctx)

		var request proto_msghandler.RouteMessageRequest
		err := msg.Unmarshal(&request)
		if err != nil {
			logger.Warn("Error unmarshalling request", "error", err)
			return proto.CompositeMessage{
				Message: proto_errors.NewMalformedMessage(request.GetMessageType()),
			}
		}

		var response proto_msghandler.RouteMessageResponse

		var jsonMessage map[string]interface{}
		err = json.Unmarshal([]byte(request.GetData()), &jsonMessage)
		if err != nil {
			//TODO
			response.Data = pbuf.String(messages.Marshal(&messages.ErrorMessage{
				Error: "todo",
			}))
			return proto.CompositeMessage{Message: &response}
		}

		cmdString, _ := jsonMessage[messages.ParamCommand].(string)
		command := messages.ParseCommand(cmdString)

		var r string
		switch command[0] {
		case messages.CmdLobby:
			r = routeLobbyMessage(command[1:], jsonMessage)
		default:
			// TODO
			r = messages.Marshal(&messages.ErrorMessage{
				Error: "todo2",
			})
		}

		response.Data = pbuf.String(r)

		return proto.CompositeMessage{Message: &response}
	})
}

func routeLobbyMessage(c []string, m map[string]interface{}) string {
	return "{}"
}
