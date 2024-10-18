package usermodels

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserResponce struct {
	UserId  string `json:"userId"`
	Message string `json:"message"`
}
