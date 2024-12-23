package database

import (
	errormodels "chat/internal/models/errorModels"
	"chat/internal/models/lobbyModels"
	logger "chat/pkg"
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
)

func ConnectDatabase(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("Error connect database")
		os.Exit(1)
	}

	return db
}

func CreateLobbyInDatabase(roomId, userID string, database *sql.DB) error {
	// query := "INSERT INTO rooms(id) VALUES($1)"
	query := `
		INSERT INTO rooms (id)
		SELECT $1
		FROM users
		WHERE users.id = $2`

	result, err := database.Exec(query, roomId, userID)
	if err != nil {
		fmt.Println(err)
		logger.Log.Error("Error create lobby in database")
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		logger.Log.Error("No access to create new lobby")
		return errormodels.NoAccessCreateLobby
	}

	return nil
}

func RegisterUser(id, username, password string, database *sql.DB) error {
	query := `INSERT INTO users(id, username, password) VALUES($1, $2, $3)`
	_, err := database.Exec(query, id, username, password)
	if err != nil {
		logger.Log.Error("Error register user")
		return err
	}
	return nil
}

func JoinLobby(userId, roomId string, database *sql.DB) error {
	query := `INSERT INTO room_users (user_id, room_id)
			SELECT users.id, rooms.id
			FROM users, rooms
			WHERE users.id = $1 AND rooms.id = $2`

	result, err := database.Exec(query, userId, roomId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			logger.Log.Error("Error join in room, user already joined lobby")
			return errormodels.ErrUserAlreadyJoined
		}
		logger.Log.Error("Error join in room")
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil || affectedRows == 0 {
		logger.Log.Error("No rows were inserted")
		return errormodels.ErrNoLobbyExist
	}

	return nil
}

func GetAllRoomFromDB(database *sql.DB) ([]string, error) {
	var roomsId []string

	query := `SELECT id FROM rooms`
	result, err := database.Query(query)
	if err != nil {
		logger.Log.Error("Error get rooms")
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		var roomId string

		err := result.Scan(&roomId)
		if err != nil {
			return nil, err
		}
		roomsId = append(roomsId, roomId)
	}

	if err := result.Err(); err != nil {
		logger.Log.Error("Error during result iteration")
		return nil, err
	}

	return roomsId, nil
}

func AddMessageInDataBase(message lobbyModels.Message, userId, roomId string, database *sql.DB) {
	query := `INSERT INTO messages (room_id, user_id, message) VALUES($1, $2, $3)`
	database.Exec(query, roomId, userId, message.Content)
}

func GetHistoryRoomFromDB(roomID, userId string, database *sql.DB) ([]lobbyModels.Message, error) {
	var historyMassage []lobbyModels.Message
	query := `SELECT message, user_id FROM messages
	 WHERE room_id = $1 AND EXISTS (SELECT 1 FROM users WHERE id = $2)`
	result, err := database.Query(query, roomID, userId)
	if err != nil {
		logger.Log.Error("Error get history room")
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var message lobbyModels.Message

		err := result.Scan(&message.Content, &message.From)
		if err != nil {
			return nil, err
		}
		historyMassage = append(historyMassage, message)
	}

	if err := result.Err(); err != nil {
		logger.Log.Error("Error during result iteration")
		return nil, err
	}

	return historyMassage, nil
}

func CheckDataForConnectionWebsocketDB(userId, roomId string, database *sql.DB) bool { // if true - OK, false - not OK
	var userExists, roomExists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1) AS user_exists, 
					 EXISTS(SELECT 1 FROM rooms WHERE id = $2) AS room_exists`

	row := database.QueryRow(query, userId, roomId)

	err := row.Scan(&userExists, &roomExists)
	if err != nil {
		return false
	}

	if !userExists || !roomExists {
		return false
	}

	return true
}
