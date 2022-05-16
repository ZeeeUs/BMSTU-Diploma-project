package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/hasher"
)

type UserUsecase interface {
	Signup(context.Context, models.User) (models.User, int, error)
}

type userUsecase struct {
	UserRepository repository.UserRepository
	Timeout        time.Duration
}

func NewUserUsecase(ur repository.UserRepository, timeout time.Duration) UserUsecase {
	return &userUsecase{
		UserRepository: ur,
		Timeout:        timeout,
	}
}

func (uu *userUsecase) Signup(ctx context.Context, user models.User) (models.User, int, error) {
	hashedPass, err := hasher.HashAndSalt(user.Password)
	if err != nil {
		return models.User{}, http.StatusInternalServerError, err
	}

	user.Password = hashedPass
	createdUser, err := uu.UserRepository.CreateUser(ctx, user)

	if err != nil {
		return models.User{}, http.StatusInternalServerError, err
	}

	return createdUser, http.StatusOK, nil
}
