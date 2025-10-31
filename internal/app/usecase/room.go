package usecase

import (
	"github.com/kittiphop/zombie_board_game_back/internal/domain/model"
	"github.com/kittiphop/zombie_board_game_back/internal/domain/repository"
)

type RoomUseCase struct {
	roomRepo repository.RoomRepository
}

func NewRoomUseCase(roomRepo repository.RoomRepository) *RoomUseCase {
	return &RoomUseCase{roomRepo: roomRepo}
}

func (uc *RoomUseCase) GetRooms() ([]*model.Room, error) {
	return uc.roomRepo.GetRooms()
}
