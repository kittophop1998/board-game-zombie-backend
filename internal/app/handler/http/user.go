package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kittiphop/zombie_board_game_back/internal/app/usecase"
)

type UserHandler struct {
	UserUC *usecase.UserUseCase
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.UserUC.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, users)
}
