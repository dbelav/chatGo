package lobbyHandlers

import (
	"database/sql"
	"errors"
	"net/http"

	"chat/internal/database"
	errormodels "chat/internal/models/errorModels"
	"chat/internal/models/lobbyModels"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateLobby(c *gin.Context, db *sql.DB) {
	idLobby := uuid.New()
	err := database.CreateLobbyInDatabase(idLobby.String(), db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, lobbyModels.ErrorCreateLobbyResponse{
			Message: "Error create lobby",
		})
		return
	}
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
			c.JSON(http.StatusInternalServerError, lobbyModels.JoinLobbyResponce{
				Message: "Error join Lobby. Lobby is not exist",
			})
			return
		}
	}
	c.JSON(http.StatusOK, lobbyModels.JoinLobbyResponce{
		Message: "Successful joined lobby",
	})
}

// 106fb31c-2af2-4c66-b163-d02815fbf95d

// fa827ac1-a55f-4d24-8b58-0fa869ef9066
// d7bd2e10-d38c-4961-add8-3e3e57bd5134
