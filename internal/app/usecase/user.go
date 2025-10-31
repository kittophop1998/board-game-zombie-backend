package usecase

import (
	"github.com/kittiphop/zombie_board_game_back/internal/domain/model"
	"github.com/kittiphop/zombie_board_game_back/internal/domain/repository"
)

type UserUseCase struct {
	UserRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepo: userRepo,
	}
}

func (uc *UserUseCase) GetUsers() ([]*model.User, error) {
	return uc.UserRepo.GetUsers()
}
