package http

import (
	errormodels "chat/internal/models/errorModels"
	"chat/internal/models/lobbyModels"
	lobbyHandlers "chat/internal/services/lobby"
	userHandlers "chat/internal/services/users"
	"chat/internal/transport/websocket"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Api(database *sql.DB) {
	routes := gin.Default()
	lobbyGroup := routes.Group("/lobby")
	{
		lobbyGroup.POST("/create", func(c *gin.Context) {
			room, err := lobbyHandlers.CreateLobby(c, database)
			if err != nil {
				c.JSON(http.StatusInternalServerError, lobbyModels.ErrorCreateLobbyResponse{
					Message: "Error create lobby",
				})
			}

			go runRoom(room)

			c.JSON(http.StatusCreated, lobbyModels.CreateLobbyResponse{
				LobbyId: room.Id,
				Message: "Successful created lobby",
			})
		})

		lobbyGroup.POST("/join", func(c *gin.Context) {
			err := lobbyHandlers.JoinLobby(c, database)
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
				if errors.Is(err, errormodels.InvalidInput) {
					c.JSON(http.StatusBadRequest, lobbyModels.JoinLobbyResponce{
						Message: "Invalid input" + err.Error(),
					})
					return
				}
			}

			c.JSON(http.StatusOK, lobbyModels.JoinLobbyResponce{
				Message: "Successful joined lobby",
			})
		})

		lobbyGroup.GET("/history", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Комната создана"})
		})
	}

	userGroup := routes.Group("/user")
	{
		userGroup.POST("/register", func(c *gin.Context) {
			userHandlers.RegisterUser(c, database)
		})
	}

	routes.Run()
}

func runRoom(room *lobbyHandlers.Room) {
	for {
		select {
		case userEvent := <-room.UserEvents:
			switch userEvent.Event {
			case "join":
				message := lobbyModels.Message{
					From:    userEvent.User.Id,
					Content: fmt.Sprintf("User %s joined", userEvent.User.Id),
				}
				room.AddUser(userEvent.User)
				go websocket.ListenMessageToUser(userEvent.User, room)
				room.Brodcast <- message
			case "leave":
				message := lobbyModels.Message{
					From:    userEvent.User.Id,
					Content: fmt.Sprintf("User %s leaved", userEvent.User.Id),
				}
				room.DeleteUser(userEvent.User)
				close(userEvent.User.CloseChan)
				room.Brodcast <- message
			}
		}
	}
}
