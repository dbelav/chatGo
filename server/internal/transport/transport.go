package transport

import (
	"chat/internal/models/lobbyModels"
	lobbyHandlers "chat/internal/services/lobby"
	"chat/internal/transport/websocket"
	"fmt"
)

func RunRoom(room *lobbyHandlers.Room) {
	for {
		select {
		case userEvent := <-room.UserEvents:
			switch userEvent.Event {
			case "connect":
				message := lobbyModels.Message{
					From:    userEvent.User.Id,
					Content: fmt.Sprintf("User %s connected", userEvent.User.Id),
				}
				room.AddUser(userEvent.User)
				fmt.Println(message)
				go websocket.ListenMessageToUser(userEvent.User, room)
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
