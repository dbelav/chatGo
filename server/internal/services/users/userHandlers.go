package userHandlers

import (
	"chat/internal/database"
	usermodels "chat/internal/models/userModels"
	"database/sql"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUser(c *gin.Context, db *sql.DB) {
	var request usermodels.RegisterUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, usermodels.RegisterUserResponce{
			Message: "Invalid input" + err.Error(),
		})
		return
	}
	if !validateUsername(request.Username) || !validatePassword(request.Password) {
		c.JSON(http.StatusBadRequest, usermodels.RegisterUserResponce{
			Message: "Invalid username or password",
		})
		return
	}

	idUser := uuid.New()

	err := database.RegisterUser(idUser.String(),
		request.Username,
		request.Password,
		db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, usermodels.RegisterUserResponce{
			Message: "Error register user",
		})
		return
	}

	c.JSON(http.StatusCreated, usermodels.RegisterUserResponce{
		UserId:  idUser.String(),
		Message: "Successful registered",
	})
}

func validateUsername(username string) bool {
	valid := regexp.MustCompile(`^[a-zA-Z0-9]{3,20}$`)
	return valid.MatchString(username)
}

func validatePassword(password string) bool {
	valid := regexp.MustCompile(`^[a-zA-Z0-9]{3,50}$`)
	return valid.MatchString(password)
}
