package lobbyModels

type CreateLobbyResponse struct {
	LobbyId string `json:"lobbyId"`
	Message string `json:"message"`
}

type ErrorCreateLobbyResponse struct {
	Message string `json:"message"`
}

type JoinLobbyRequest struct {
	LobbyId string `json:"lobbyId"`
	UserId  string `json:"userId"`
}

type JoinLobbyResponce struct {
	Message string `json:"message"`
}
