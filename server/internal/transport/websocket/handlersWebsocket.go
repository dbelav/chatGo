package websocket

import (
	"chat/internal/models/lobbyModels"
	lobbyHandlers "chat/internal/services/lobby"
	logger "chat/pkg"
	"database/sql"
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

func ListenUserMessage(user *lobbyModels.User, room *lobbyHandlers.Room, db *sql.DB) {
	for {
		select {
		case <-user.CloseChan:
			fmt.Println("CloseChan")
			return

		default: // default handler message
			var message lobbyModels.Message
			err := user.Connection.ReadJSON(&message)
			lobbyHandlers.SaveMassageInHistory(message, user.Id, room.Id, db)
			if err != nil {
				logger.Log.Error("Error reading message from user %s: %v", user.Id, err)
				return
			}
			room.Brodcast <- message 
		}
	}
}

func ListenUserChanelFromBrodcast(user *lobbyModels.User, room *lobbyHandlers.Room) {
	for {
		channelInterface, ok := room.UserChannels.Load(user.Id)
		if !ok {
			return
		}
		userChannel, ok := channelInterface.(chan lobbyModels.Message)
		if !ok {
			return
		}
		select {
		case message := <-userChannel:
			err := HandlerSendMessageBrodcast(user, message)
			if err != nil {
				logger.Log.Error("Error sending message to user %s: %v", user.Id, err)
			}

		default:
		}
	}
}

func ListenBrodcast(room *lobbyHandlers.Room) {
	for brodcast := range room.Brodcast {
		room.UserChannels.Range(func(_, userChanInterface any) bool {
			userChan, ok := userChanInterface.(chan lobbyModels.Message)
			if ok {
				userChan <- brodcast
			} else {
				logger.Log.Warn("Failed to send message, user channel type mismatch")
			}
			return true
		})
		fmt.Println("message sent")
	}
}
