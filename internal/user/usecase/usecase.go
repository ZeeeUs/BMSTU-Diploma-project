package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/hasher"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx"
)

type UserUsecase interface {
	UserLogin(context.Context, models.UserCredentials) (models.User, int, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserById(ctx context.Context, id int) (models.User, error)
}

type userUsecase struct {
	UserRepository repository.UserRepository
	Timeout        time.Duration
	logger         *logrus.Logger
}

func NewUserUsecase(ur repository.UserRepository, timeout time.Duration, log *logrus.Logger) UserUsecase {
	return &userUsecase{
		UserRepository: ur,
		Timeout:        timeout,
		logger:         log,
	}
}

func (uu *userUsecase) UserLogin(ctx context.Context, creds models.UserCredentials) (models.User, int, error) {
	user, err := uu.UserRepository.GetUserByEmail(ctx, creds.Email)
	if err == pgx.ErrNoRows {
		return models.User{}, http.StatusNotFound, fmt.Errorf("user with email %s is not found", creds.Email)
	} else if err != nil {
		return models.User{}, http.StatusInternalServerError, err
	}

	if user.PassStatus == 0 {
		return user, http.StatusOK, nil
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

func (uu *userUsecase) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	//user, err := uu.UserRepository.UpdateUser(ctx)
	return models.User{}, nil
}

func (uu *userUsecase) GetUserById(ctx context.Context, id int) (models.User, error) {
	return uu.UserRepository.GetUserById(ctx, id)
}
