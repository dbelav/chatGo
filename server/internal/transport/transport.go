package transport

import (
	"chat/internal/models/lobbyModels"
	lobbyHandlers "chat/internal/services/lobby"
	"chat/internal/transport/websocket"
	"database/sql"
	"fmt"
)

func RunRoom(room *lobbyHandlers.Room, db *sql.DB) {
	for {
		select {
		case userEvent := <-room.UserEvents:
			go websocket.ListenBrodcast(room)

			switch userEvent.Event {
			case "connect":
				message := lobbyModels.Message{
					From:    userEvent.User.Id,
					Content: fmt.Sprintf("User %s connected", userEvent.User.Id),
				}
				room.AddUser(userEvent.User)
				go websocket.ListenUserMessage(userEvent.User, room, db)
				go websocket.ListenUserChanelFromBrodcast(userEvent.User, room)
				room.Brodcast <- message
			case "disconnect":
				message := lobbyModels.Message{
					From:    userEvent.User.Id,
					Content: fmt.Sprintf("User %s disconnected", userEvent.User.Id),
				}
				room.DeleteUser(userEvent.User)
				close(userEvent.User.CloseChan)
				room.Brodcast <- message
			}
		}
	}
}
