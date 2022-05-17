package usecase

import (
	"context"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	"github.com/sirupsen/logrus"
)

type SessionUsecas interface {
	AddSession(context.Context, models.Session) error
}

type sessionUsecase struct {
	Session        repository.SessionRepository
	contextTimeout time.Duration
	logger         *logrus.Logger
}

func NewSessionUsecase(session repository.SessionRepository, timeout time.Duration, logger *logrus.Logger) SessionUsecas {
	return &sessionUsecase{
		Session:        session,
		contextTimeout: timeout,
		logger:         logger,
	}
}

func (su *sessionUsecase) AddSession(ctx context.Context, session models.Session) error {
	err := su.Session.NewSessionCookie(ctx, session.Cookie, session.Id)
	if err != nil {
		return err
	}
	return nil
}
