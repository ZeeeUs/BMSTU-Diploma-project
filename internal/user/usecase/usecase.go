package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/storage"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/hasher"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx"
)

type UserUsecase interface {
	UserLogin(context.Context, models.UserCredentials) (models.User, error)
	UpdateUser(ctx context.Context, newPass string, email string) (int, error)
	GetUserById(ctx context.Context, id int) (models.User, error)
}

type userUsecase struct {
	UserRepository storage.UserStorage
	Timeout        time.Duration
	logger         *logrus.Logger
}

func NewUserUsecase(ur storage.UserStorage, timeout time.Duration, log *logrus.Logger) UserUsecase {
	return &userUsecase{
		UserRepository: ur,
		Timeout:        timeout,
		logger:         log,
	}
}

func (uu *userUsecase) UserLogin(ctx context.Context, creds models.UserCredentials) (models.User, error) {
	user, err := uu.UserRepository.GetUserByEmail(ctx, creds.Email)
	if err == pgx.ErrNoRows {
		return models.User{}, fmt.Errorf("user with email %s is not found", creds.Email)
	} else if err != nil {
		return models.User{}, err
	}

	if !user.PassStatus {
		if strings.Compare(user.Password, creds.Password) == 0 {
			return user, nil
		}

		return models.User{}, errors.New("401")
	}

	isVerify, err := hasher.ComparePasswords(user.Password, creds.Password)
	if err != nil {
		return models.User{}, err
	}

	if !isVerify {
		return models.User{}, err
	}

	return user, nil
}

func (uu *userUsecase) UpdateUser(ctx context.Context, newPass string, email string) (int, error) {
	id, err := uu.UserRepository.UpdateUser(ctx, newPass, email)
	if err != nil {
		uu.logger.Errorf("User use case: faile to UpdateUser: %s", err)
		return 0, err
	}

	if id <= 0 {
		return 0, errors.New("update table in db is not correct")
	}

	return id, nil
}

func (uu *userUsecase) GetUserById(ctx context.Context, id int) (models.User, error) {
	return uu.UserRepository.GetUserById(ctx, id)
}
