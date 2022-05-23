package usecase

import (
	"context"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	"github.com/sirupsen/logrus"
)

type SessionUsecase interface {
	AddSession(context.Context, models.Session) error
	GetSessionByToken(context.Context, string) (models.Session, error)
	DeleteSession(ctx context.Context) error
}

type sessionUsecase struct {
	Session        repository.SessionRepository
	contextTimeout time.Duration
	logger         *logrus.Logger
}

func NewSessionUsecase(session repository.SessionRepository, timeout time.Duration, logger *logrus.Logger) SessionUsecase {
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

func (su *sessionUsecase) GetSessionByToken(ctx context.Context, token string) (models.Session, error) {
	id, err := su.Session.GetSessionByToken(ctx, token)
	if err != nil {
		return models.Session{}, err
	}

	return models.Session{Cookie: token, Id: id}, nil
}

func (su *sessionUsecase) DeleteSession(ctx context.Context) error {
	err := su.Session.DeleteSession(ctx)
	if err != nil {
		return err
	}
	return nil
}
