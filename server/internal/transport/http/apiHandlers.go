package http

import (
	errormodels "chat/internal/models/errorModels"
	"chat/internal/models/lobbyModels"
	lobbyHandlers "chat/internal/services/lobby"
	"chat/internal/transport"
	"database/sql"
	"errors"
	"fmt"
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
				Message: "Error join Lobby.",
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
		if errors.Is(err, errormodels.NoAccessCreateLobby) {
			c.JSON(http.StatusInternalServerError, lobbyModels.ErrorResponce{
				Message: "Error, No access to create lobby",
			})
		} else {
			c.JSON(http.StatusInternalServerError, lobbyModels.ErrorResponce{
				Message: "Error create lobby",
			})
		}
		return
	}

	go transport.RunRoom(room, database)

	c.JSON(http.StatusCreated, lobbyModels.CreateLobbyResponse{
		LobbyId: room.Id,
		Message: "Successful created lobby",
	})
}

func GetMessagesHistoryByRoomHandler(c *gin.Context, database *sql.DB) {
	result, err := lobbyHandlers.GetRoomHistory(c, database)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, errormodels.RequiredQueryParams) {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, lobbyModels.ErrorResponce{
				Message: "Error, no query params",
			})
		} else {
			c.JSON(http.StatusInternalServerError, lobbyModels.ErrorResponce{
				Message: "Error get history",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, lobbyModels.HistoryRoomResponce{
		MessagesHistory: result,
	})

}
