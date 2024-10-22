package app

import (
	"fmt"
	"os"

	"chat/internal/database"
	"chat/internal/transport/http"
	logger "chat/pkg"
)

func StartApp() {
	logger.InitLogs()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	database := database.ConnectDatabase(psqlInfo)

	http.Api(database)
	// websocket.StartWebsocket()
}
