package database

import (
	"database/sql"
	"os"

	errormodels "chat/internal/models/errorModels"
	logger "chat/pkg"

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

func CreateLobbyInDatabase(id string, database *sql.DB) error {
	query := "INSERT INTO rooms(id) VALUES($1)"
	_, err := database.Exec(query, id)
	if err != nil {
		logger.Log.Error("Error create lobby in database")
		return err
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
