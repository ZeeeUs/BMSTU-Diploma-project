package storage

import (
	"context"
	"strconv"

	redisClient "github.com/ZeeeUs/BMSTU-Diploma-project/pkg/redis"

	"github.com/sirupsen/logrus"
)

type SessionStorage interface {
	NewSessionCookie(context.Context, string, int) error
	GetSessionByToken(ctx context.Context, cookieVal string) (int, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

type sessionStorage struct {
	RedisConnection redisClient.Client
	logger          *logrus.Logger
}

func NewSessionStorage(conn redisClient.Client, log *logrus.Logger) SessionStorage {
	return &sessionStorage{
		RedisConnection: conn,
		logger:          log,
	}
}

func (sr *sessionStorage) NewSessionCookie(ctx context.Context, sessionCookie string, id int) error {
	err := sr.RedisConnection.Set(sessionCookie, id)
	if err != nil {
		return err
	}
	return nil
}

func (sr *sessionStorage) GetSessionByToken(ctx context.Context, cookieVal string) (int, error) {
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

func (sr *sessionStorage) DeleteSession(ctx context.Context, sessionId string) error {
	err := sr.RedisConnection.DeleteKeyValue(sessionId)
	if err != nil {
		return err
	}
	return nil
}
