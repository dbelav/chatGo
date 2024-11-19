package app

import (
	"fmt"
	"os"

	"chat/internal/database"
	lobbyHandlers "chat/internal/services/lobby"
	"chat/internal/transport"
	"chat/internal/transport/http"
	logger "chat/pkg"
)

func StartApp() {
	fmt.Println("Test before logger initialization")
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

	err := lobbyHandlers.InitAllRooms(database) // get from DB and create exemplars
	if err != nil {
		logger.Log.Error("Error start all rooms")
		os.Exit(1)
	}

	for _, room := range lobbyHandlers.Rooms { // start all existed room at start
		go transport.RunRoom(room)
	}

	http.Api(database)

	// websocket.StartWebsocket()
}
