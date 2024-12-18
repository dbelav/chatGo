package http

import (
	userHandlers "chat/internal/services/users"
	"chat/internal/transport/websocket"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Api(database *sql.DB) {
	routes := gin.Default()

	routes.GET("/ws", func(c *gin.Context) {
		userId := c.Query("userId")
		lobbyId := c.Query("lobbyId")
		useName := c.Query("userName")
		websocket.HandlerConnection(c.Writer, c.Request, userId, lobbyId, useName)
	})

	lobbyGroup := routes.Group("/lobby")
	{
		lobbyGroup.POST("/create", func(c *gin.Context) {
			CreateRoomHandler(c, database)
		})

		lobbyGroup.POST("/join", func(c *gin.Context) {
			JoinLobbyHandler(c, database)
		})

		lobbyGroup.GET("/history", func(c *gin.Context) {
			GetMessagesHistoryByRoomHandler(c, database)
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
