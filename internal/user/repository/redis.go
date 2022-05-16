package repository

import (
	"context"

	redisClient "github.com/ZeeeUs/BMSTU-Diploma-project/pkg/redis"

	"github.com/sirupsen/logrus"
)

type SessionRepository interface {
	NewSessionCookie(context.Context, string, int) error
}

type sessionRepository struct {
	RedisConnection redisClient.Client
	logger          *logrus.Logger
}

func NewSessionRepository(conn redisClient.Client, log *logrus.Logger) SessionRepository {
	return &sessionRepository{
		RedisConnection: conn,
		logger:          log,
	}
}

func (sr *sessionRepository) NewSessionCookie(ctx context.Context, sessionCookie string, id int) error {
	err := sr.RedisConnection.Set(sessionCookie, id)
	if err != nil {
		return err
	}
	return nil
}
