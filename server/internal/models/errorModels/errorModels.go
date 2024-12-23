package errormodels

import "errors"

var ErrUserAlreadyJoined = errors.New("user already joined lobby")
var ErrNoLobbyExist = errors.New("Lobby is no exist")
var InvalidInput = errors.New("Invalid input")
var NoAccessCreateLobby = errors.New("No access to create new lobby")
var RequiredQueryParams = errors.New("Required query params")
