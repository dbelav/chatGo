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
			fmt.Println(userEvent)
			switch userEvent.Event {
			case "connect":
				fmt.Println("userEvent Connect")
				message := lobbyModels.Message{
					From:    userEvent.User.Id,
					Content: fmt.Sprintf("User %s connected", userEvent.User.Id),
				}
				room.AddUser(userEvent.User)
				go websocket.ListenUserMessage(userEvent.User, room)
				go websocket.ListenBrodcast(userEvent.User, room)
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
