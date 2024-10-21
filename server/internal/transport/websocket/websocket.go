package websocket

import (
	// "chat/internal/models/lobbyModels"
	// lobbyHandlers "chat/internal/services/lobby"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWebsocket() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		lobbyId := r.URL.Query().Get("lobbyId")

		HandlerConnection(upgrader, w, r, userId, lobbyId)
	})
}

// func HandlerConnection(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request, room *lobbyHandlers.Room, userId string) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
// 		return
// 	}
// 	defer conn.Close()

// 	user := &lobbyModels.User{
// 		Id:         userId,
// 		Connection: conn,
// 	}
// 	room.AddUserEvent(user, "join")
// }

// func HandlerSendMessageBrodcast(user *lobbyModels.User, message lobbyModels.Message) error {
// 	err := user.Connection.WriteJSON(message)
// 	if err != nil {
// 		logger.Log.Error("Error sending message to user %s: %v", user.Id, err)
// 		return err
// 	}
// 	return nil
// }
