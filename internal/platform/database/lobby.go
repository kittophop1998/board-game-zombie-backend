package database

import "gorm.io/gorm"

type LobbyPostGres struct {
	db *gorm.DB
}

func NewLobbyPostGres(db *gorm.DB) *LobbyPostGres {
	return &LobbyPostGres{db: db}
}
