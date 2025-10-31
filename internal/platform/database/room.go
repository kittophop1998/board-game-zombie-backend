package database

import (
	"github.com/kittiphop/zombie_board_game_back/internal/domain/model"
	"gorm.io/gorm"
)

type RoomPostGres struct {
	db *gorm.DB
}

func NewRoomPostGres(db *gorm.DB) *RoomPostGres {
	return &RoomPostGres{db: db}
}

func (r *RoomPostGres) GetRooms() ([]*model.Room, error) {
	var rooms []*model.Room
	if err := r.db.Find(&rooms).Error; err != nil {
		return nil, err
	}

	return rooms, nil
}
