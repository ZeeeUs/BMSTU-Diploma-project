package repository

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	redisClient "github.com/ZeeeUs/BMSTU-Diploma-project/pkg/redis"

	"github.com/sirupsen/logrus"
)

type SessionRepository interface {
	NewSessionCookie(context.Context, string, int) error
	GetSessionByToken(ctx context.Context, cookieVal string) (int, error)
	DeleteSession(ctx context.Context) error
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

func (sr *sessionRepository) GetSessionByToken(ctx context.Context, cookieVal string) (int, error) {
	var id interface{}
	id, err := sr.RedisConnection.GetValue(cookieVal)
	if err != nil {
		sr.logger.Errorf("Check auth cookie: cookie not found: %s", err)
		return 0, err
	}

	newId, err := strconv.Atoi(id.(string))
	if err != nil {
		sr.logger.Println(err)
		return 0, err
	}

	return newId, err
}

func (sr *sessionRepository) DeleteSession(ctx context.Context) error {
	ctxSession := ctx.Value("userId")
	if ctxSession == nil {
		return errors.New("context nil error")
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return errors.New("convert to model session error")
	}

	log.Println(currentSession.Cookie)

	err := sr.RedisConnection.DeleteKeyValue(currentSession.Cookie)
	if err != nil {
		return err
	}
	return nil
}
