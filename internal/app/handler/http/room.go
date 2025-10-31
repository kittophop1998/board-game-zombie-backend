package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kittiphop/zombie_board_game_back/infrastructure/lib"
	"github.com/kittiphop/zombie_board_game_back/internal/app/usecase"
)

type RoomHandler struct {
	RoomUC *usecase.RoomUseCase
}

func (h *RoomHandler) GetRooms(c *gin.Context) {
	rooms, err := h.RoomUC.GetRooms()
	if err != nil {
		lib.ResponseInternalServerError(c, "Cannot get rooms", err.Error())
		return
	}

	lib.ResponseSuccess(c, rooms)
}
