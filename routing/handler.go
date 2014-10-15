package routing

import (
	"encoding/json"
	"time"

	"code.google.com/p/go.net/context"
	pbuf "code.google.com/p/gogoprotobuf/proto"

	"github.com/opentarock/service-api/go/client"
	"github.com/opentarock/service-api/go/proto"
	"github.com/opentarock/service-api/go/proto_errors"
	"github.com/opentarock/service-api/go/proto_lobby"
	"github.com/opentarock/service-api/go/proto_msghandler"
	"github.com/opentarock/service-api/go/reqcontext"
	"github.com/opentarock/service-api/go/user"
	"github.com/opentarock/service-msghandler/messages"
)

const defaultRequestTimeout = 1 * time.Minute

type RouteMessageHandler struct {
	lobbyClient client.LobbyClient
}

func NewRouteMessageHandler(lobbyClient client.LobbyClient) *RouteMessageHandler {
	return &RouteMessageHandler{
		lobbyClient: lobbyClient,
	}
}

func (h *RouteMessageHandler) HandleMessage(msg *proto.Message) proto.CompositeMessage {
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
		response.Data = pbuf.String(messages.Marshal(messages.NewInvalidRequestMalformed()))
		return proto.CompositeMessage{Message: &response}
	}

	command, _ := jsonMessage[messages.ParamCommand].(string)

	logger.Info("Routing message", "command", command)

	var r string
	switch commandHead(command) {
	case messages.CmdLobby:
		r = h.routeLobbyMessage(ctx, command, jsonMessage)
	default:
		r = messages.Marshal(messages.NewUnknownCommandError(command))
	}

	response.Data = pbuf.String(r)

	return proto.CompositeMessage{Message: &response}
}

func (h *RouteMessageHandler) routeLobbyMessage(ctx context.Context, c string, m map[string]interface{}) string {
	logger := reqcontext.ContextLogger(ctx)
	switch c {
	case "lobby.room.create":
		roomName, _ := m["name"].(string)
		result, err := h.lobbyClient.CreateRoom(ctx, roomName, nil)
		if err != nil {
			logger.Error("Failed to create room", "error", err)
			return messages.Marshal(messages.NewServerError())
		}
		m := make(map[string]interface{})
		m[messages.ParamResponse] = "lobby.room.create"
		room, err := toJsonRoom(ctx, result.GetRoom())
		if err != nil {
			logger.Error("Problem getting room data", "error", err)
			return messages.Marshal(messages.NewServerError())
		}
		m["room"] = room
		return messages.Marshal(m)
	case "lobby.room.list":
		result, err := h.lobbyClient.ListRooms(ctx)
		if err != nil {
			logger.Error("Failed to list rooms", "error", err)
			return messages.Marshal(messages.NewServerError())
		}
		m := make(map[string]interface{})
		m[messages.ParamResponse] = "lobby.room.list"
		roomList := make([]map[string]interface{}, 0, len(result.GetRooms()))
		for _, room := range result.GetRooms() {
			r, err := toJsonRoom(ctx, room)
			if err != nil {
				logger.Error("Problem getting room data", "error", err)
				return messages.Marshal(messages.NewServerError())
			}
			roomList = append(roomList, r)
		}
		m["rooms"] = roomList
		return messages.Marshal(m)
	default:
		return messages.Marshal(messages.NewUnknownCommandError(c))
	}
}

func toJsonRoom(ctx context.Context, room *proto_lobby.Room) (map[string]interface{}, error) {
	r := make(map[string]interface{})
	r["id"] = room.GetId()
	r["name"] = room.GetName()
	owner, err := fetchPlayerData(ctx, user.Id(room.GetOwner()))
	if err != nil {
		return nil, err
	}
	r["owner"] = owner
	players := make([]map[string]interface{}, 0, len(room.GetPlayers()))
	for _, p := range room.GetPlayers() {
		player, err := fetchPlayerData(ctx, user.Id(p))
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	r["players"] = players
	return r, err
}

func fetchPlayerData(ctx context.Context, userId user.Id) (map[string]interface{}, error) {
	p := make(map[string]interface{})
	p["id"] = userId
	p["nickname"] = "user nickname"
	return p, nil
}

func commandHead(c string) string {
	parts := messages.ParseCommand(c)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
