package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kittiphop/zombie_board_game_back/infrastructure/lib"
	"github.com/kittiphop/zombie_board_game_back/internal/app/usecase"
)

type UserHandler struct {
	UserUC *usecase.UserUseCase
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.UserUC.GetUsers()
	if err != nil {
		lib.ResponseInternalServerError(c, "Cannot get users", err.Error())
		return
	}

	lib.ResponseSuccess(c, users)
}
