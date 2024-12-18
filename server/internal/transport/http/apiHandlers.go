package http

import (
	errormodels "chat/internal/models/errorModels"
	"chat/internal/models/lobbyModels"
	lobbyHandlers "chat/internal/services/lobby"
	"chat/internal/transport"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JoinLobbyHandler(c *gin.Context, database *sql.DB) {
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
}

func CreateRoomHandler(c *gin.Context, database *sql.DB) {
	room, err := lobbyHandlers.CreateLobby(c, database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, lobbyModels.ErrorCreateLobbyResponse{
			Message: "Error create lobby",
		})
	}

	go transport.RunRoom(room, database)

	c.JSON(http.StatusCreated, lobbyModels.CreateLobbyResponse{
		LobbyId: room.Id,
		Message: "Successful created lobby",
	})
}

func GetMessagesHistoryByRoomHandler(c *gin.Context, database *sql.DB) {

}
