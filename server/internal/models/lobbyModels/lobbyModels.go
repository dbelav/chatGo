package lobbyModels

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
}

type User struct {
	Id         string
	Connection *websocket.Conn
	CloseChan  chan struct{}
}

type UserEvents struct {
	User  *User
	Event string
}

type CreateLobbyResponse struct {
	LobbyId string `json:"lobbyId"`
	Message string `json:"message"`
}

type ErrorCreateLobbyResponse struct {
	Message string `json:"message"`
}

type JoinLobbyRequest struct {
	LobbyId string `json:"lobbyId"`
	UserId  string `json:"userId"`
}

type JoinLobbyResponce struct {
	Message string `json:"message"`
}
