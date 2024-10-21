package errormodels

import "errors"

var ErrUserAlreadyJoined = errors.New("user already joined lobby")
var ErrNoLobbyExist = errors.New("Lobby is no exist")
var InvalidInput = errors.New("Invalid input")
