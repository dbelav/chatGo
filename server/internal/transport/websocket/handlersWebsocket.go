package websocket

import (
	"chat/internal/models/lobbyModels"
	lobbyHandlers "chat/internal/services/lobby"
	logger "chat/pkg"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandlerConnection(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request, userId, lobbyId string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	user := &lobbyModels.User{
		Id:         userId,
		Connection: conn,
		CloseChan: make(chan struct{}),
	}
	room := lobbyHandlers.GetRoomById(lobbyId)
	room.AddUserEvent(user, "join")
}

func HandlerSendMessageBrodcast(user *lobbyModels.User, message lobbyModels.Message) error {
	err := user.Connection.WriteJSON(message)
	if err != nil {
		logger.Log.Error("Error sending message to user %s: %v", user.Id, err)
		return err
	}
	return nil
}

func ListenMessageToUser(user *lobbyModels.User, room *lobbyHandlers.Room) {
	for {
		select {
		case brodcast := <-room.Brodcast:
			err := HandlerSendMessageBrodcast(user, brodcast)
			if err != nil {
				logger.Log.Error("Error sen message to %s", user.Id)
			}
		}

	}
}
