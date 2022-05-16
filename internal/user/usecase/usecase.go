package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/hasher"

	"github.com/jackc/pgx"
)

type UserUsecase interface {
	UserLogin(context.Context, models.UserCredentials) (models.User, int, error)
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

func (uu *userUsecase) UserLogin(ctx context.Context, creds models.UserCredentials) (models.User, int, error) {
	user, err := uu.UserRepository.GetUserByEmail(ctx, creds.Email)
	if err == pgx.ErrNoRows {
		return models.User{}, http.StatusNotFound, fmt.Errorf("user with email %s is not found", creds.Email)
	} else if err != nil {
		return models.User{}, http.StatusInternalServerError, err
	}

	isVerify, err := hasher.ComparePasswords(user.Password, creds.Password)
	if err != nil {
		return models.User{}, http.StatusInternalServerError, err
	}

	if !isVerify {
		return models.User{}, http.StatusForbidden, err
	}

	return user, http.StatusOK, nil
}
