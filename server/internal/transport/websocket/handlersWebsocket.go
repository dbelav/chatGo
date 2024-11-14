package websocket

import (
	"chat/internal/models/lobbyModels"
	lobbyHandlers "chat/internal/services/lobby"
	logger "chat/pkg"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandlerConnection(w http.ResponseWriter, r *http.Request, userId, lobbyId string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
		return
	}
	// defer conn.Close()

	user := &lobbyModels.User{
		Id:         userId,
		Connection: conn,
		CloseChan:  make(chan struct{}),
	}
	room := lobbyHandlers.GetRoomById(lobbyId)
	room.AddUserEvent(user, "connect")
}

func HandlerSendMessageBrodcast(user *lobbyModels.User, message lobbyModels.Message) error {
	err := user.Connection.WriteJSON(message)
	if err != nil {
		logger.Log.Error("Error sending message to user %s: %v", user.Id, err)
		return err
	}
	return nil
}

func ListenUserMessage(user *lobbyModels.User, room *lobbyHandlers.Room) {
	for {
		select {
		case <-user.CloseChan:
			fmt.Println("CloseChan")
			fmt.Println("CloseChan")
			return
		default:

			var message lobbyModels.Message
			err := user.Connection.ReadJSON(&message)
			fmt.Println("BRODCAST77777")
			fmt.Println(message)
			fmt.Println("BRODCAST77777")
			if err != nil {
				logger.Log.Error("Error reading message from user %s: %v", user.Id, err)
				return
			}

			room.Brodcast <- message
		}
		// fmt.Println("BRODCAST222")
		// fmt.Println(room.Brodcast)
		// fmt.Println("BRODCAST222")
		// select {

		// case brodcast := <-room.Brodcast:
		// 	fmt.Println("BRODCAST")
		// 	fmt.Println(brodcast)
		// 	fmt.Println("BRODCAST")
		// 	err := HandlerSendMessageBrodcast(user, brodcast)
		// 	if err != nil {
		// 		logger.Log.Error("Error sending message to user %s: %v", user.Id, err)
		// 	}
		// default:
		// }
	}
}

func ListenBrodcast(user *lobbyModels.User, room *lobbyHandlers.Room) {
	for brodcast := range room.Brodcast {
		err := HandlerSendMessageBrodcast(user, brodcast)
		if err != nil {
			logger.Log.Error("Error sending message to user %s: %v", user.Id, err)
			return
		}
	}
}
