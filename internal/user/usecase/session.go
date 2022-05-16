package usecase

import (
	"context"
	"time"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
)

type SessionUsecas interface {
	AddSession(context.Context, models.Session) error
}

type sessionUsecase struct {
	Session        repository.SessionRepository
	contextTimeout time.Duration
}

func NewSessionUsecase(session repository.SessionRepository, timeout time.Duration) SessionUsecas {
	return &sessionUsecase{
		Session:        session,
		contextTimeout: timeout,
	}
}

func (su *sessionUsecase) AddSession(ctx context.Context, session models.Session) error {
	err := su.Session.NewSessionCookie(ctx, session.Cookie, session.Id)
	if err != nil {
		return err
	}
	return nil
}
