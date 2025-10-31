package repository

import "github.com/kittiphop/zombie_board_game_back/internal/domain/model"

type RoomRepository interface {
	GetRooms() ([]*model.Room, error)
}
