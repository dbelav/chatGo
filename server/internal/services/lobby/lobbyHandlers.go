package lobbyHandlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"chat/internal/database"
	errormodels "chat/internal/models/errorModels"
	"chat/internal/transport/websocket"
	logger "chat/pkg"

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
func newRoom(id string) *Room {
	return &Room{
		Id:         id,
		Users:      make(map[string]*lobbyModels.User),
		Brodcast:   make(chan lobbyModels.Message, 1000),
		UserEvents: make(chan lobbyModels.UserEvents),
		Quit:       make(chan bool),
	}
}

func (r *Room) runRoom() {
	for {
		select {
		case userEvent := <-r.UserEvents:
			switch userEvent.Event {
			case "join":
				message := lobbyModels.Message{
					From:    userEvent.User.Id,
					Content: fmt.Sprintf("User %s joined", userEvent.User.Id),
				}
				r.addUser(userEvent.User)
				go r.listenMessageToUser(userEvent.User)
				r.Brodcast <- message
			case "leave":
				message := lobbyModels.Message{
					From:    userEvent.User.Id,
					Content: fmt.Sprintf("User %s leaved", userEvent.User.Id),
				}
				r.deleteUser(userEvent.User)
				go r.listenMessageToUser(userEvent.User)
				r.Brodcast <- message
			}
		}
	}
}

func (r *Room) listenMessageToUser(user *lobbyModels.User) {
	for {
		select {
		case brodcast := <-r.Brodcast:
			err := websocket.HandlerSendMessageBrodcast(user, brodcast)
			if err != nil {
				logger.Log.Error("Error sen message to %s", user.Id)
			}
		}

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

func (r *Room) addUser(user *lobbyModels.User) {
	r.Users[user.Id] = user
}

func (r *Room) deleteUser(user *lobbyModels.User) {
	delete(r.Users, user.Id)
}

func CreateLobby(c *gin.Context, db *sql.DB) {
	idLobby := uuid.New()
	err := database.CreateLobbyInDatabase(idLobby.String(), db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, lobbyModels.ErrorCreateLobbyResponse{
			Message: "Error create lobby",
		})
		return
	}
	room := newRoom(idLobby.String())
	Rooms[room.Id] = room // add in map all room
	// go runRoom()

	c.JSON(http.StatusCreated, lobbyModels.CreateLobbyResponse{
		LobbyId: idLobby.String(),
		Message: "Successful created lobby",
	})
}

func JoinLobby(c *gin.Context, db *sql.DB) {
	var request lobbyModels.JoinLobbyRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, lobbyModels.JoinLobbyResponce{
			Message: "Invalid input" + err.Error(),
		})
		return
	}

	err := database.JoinLobby(request.UserId, request.LobbyId, db)
	if err != nil {
		if errors.Is(err, errormodels.ErrUserAlreadyJoined) {
			c.JSON(http.StatusConflict, lobbyModels.JoinLobbyResponce{
				Message: "Error join Lobby. User have already joined",
			})
			return
		}
		if errors.Is(err, errormodels.ErrNoLobbyExist) {
			c.JSON(http.StatusConflict, lobbyModels.JoinLobbyResponce{
				Message: "Error join Lobby. Lobby is not exist",
			})
			return
		}
	}
	c.JSON(http.StatusOK, lobbyModels.JoinLobbyResponce{
		Message: "Successful joined lobby",
	})
}

// func (r *lobbyModels.Rooms) ()

// 106fb31c-2af2-4c66-b163-d02815fbf95d

// fa827ac1-a55f-4d24-8b58-0fa869ef9066
// d7bd2e10-d38c-4961-add8-3e3e57bd5134
