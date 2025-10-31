package http

import "github.com/kittiphop/zombie_board_game_back/internal/app/usecase"

type Handlers struct {
	User *UserHandler
	Room *RoomHandler
}

type HandlerDeps struct {
	UserUC *usecase.UserUseCase
	RoomUC *usecase.RoomUseCase
}

var H *Handlers

func InitializeHandlers(deps *HandlerDeps) {
	H = &Handlers{
		User: &UserHandler{UserUC: deps.UserUC},
		Room: &RoomHandler{RoomUC: deps.RoomUC},
	}
}
