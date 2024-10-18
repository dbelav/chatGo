package http

import (
	lobbyHandlers "chat/internal/services/lobby"
	userHandlers "chat/internal/services/users"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Api(database *sql.DB) {
	routes := gin.Default()
	lobbyGroup := routes.Group("/lobby")
	{
		lobbyGroup.POST("/create", func(c *gin.Context) {
			lobbyHandlers.CreateLobby(c, database)
		})

		lobbyGroup.POST("/join", func(c *gin.Context) {
			lobbyHandlers.JoinLobby(c, database)
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
