package lobbyHandlers

import (
	"chat/internal/database"
	errormodels "chat/internal/models/errorModels"
	"database/sql"

	"chat/internal/models/lobbyModels"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Room struct {
	Id         string
	Users      map[string]*lobbyModels.User
	Brodcast   chan lobbyModels.Message
	UserEvents chan lobbyModels.UserEvents
	Quit       chan bool
}

var Rooms = make(map[string]*Room) // all existed rooms

// get room from all list rooms
func GetRoomById(id string) *Room {
	return Rooms[id]
}

// create room and this exemplar will using with goroutines
func NewRoom(id string) *Room {
	return &Room{
		Id:         id,
		Users:      make(map[string]*lobbyModels.User),
		Brodcast:   make(chan lobbyModels.Message, 1000),
		UserEvents: make(chan lobbyModels.UserEvents),
		Quit:       make(chan bool),
	}
}

func (r *Room) SendMessage(msg lobbyModels.Message) {
	r.Brodcast <- msg
}

func (r *Room) AddUserEvent(user *lobbyModels.User, event string) {
	r.UserEvents <- lobbyModels.UserEvents{
		User:  user,
		Event: event,
	}
}

func (r *Room) AddUser(user *lobbyModels.User) {
	r.Users[user.Id] = user
}

func (r *Room) DeleteUser(user *lobbyModels.User) {
	delete(r.Users, user.Id)
}

func CreateLobby(c *gin.Context, db *sql.DB) (*Room, error) {
	idLobby := uuid.New()

	err := database.CreateLobbyInDatabase(idLobby.String(), db)
	if err != nil {
		return nil, err
	}

	room := NewRoom(idLobby.String())
	Rooms[room.Id] = room // add in map all room

	return room, nil
}

func JoinLobby(c *gin.Context, db *sql.DB) error {
	var request lobbyModels.JoinLobbyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		return errormodels.InvalidInput
	}

	err := database.JoinLobby(request.UserId, request.LobbyId, db)
	if err != nil {
		return err
	}

	return nil
}
