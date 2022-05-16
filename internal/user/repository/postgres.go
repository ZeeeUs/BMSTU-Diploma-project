package repository

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/config"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

type userRepository struct {
	Conn   *pgx.ConnPool
	logger *logrus.Logger
}

func NewUserRepository(config *config.Config, logger *logrus.Logger) (UserRepository, error) {
	pgConn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     config.DbConfig.DbHostName,
			Port:     uint16(config.DbConfig.DbPort),
			Database: config.DbConfig.DbName,
			User:     config.DbConfig.DbUser,
			Password: config.DbConfig.DbPassword,
		},
	})
	if err != nil {
		logger.Fatalf("Error %s occurred during connection to database", err)
	}

	return &userRepository{pgConn, logger}, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var newUser models.User

	err := u.Conn.QueryRow("insert into bmstudb.user (password, firstname, middle_name, lastname, email)"+
		" values ($1, $2, $3, $4, $5)"+
		" returning password, firstname, middle_name, lastname, email",
		user.Password, user.Firstname, user.MiddleName, user.Lastname, user.Email).Scan(
		&newUser.Password,
		&newUser.Firstname,
		&newUser.MiddleName,
		&newUser.Lastname,
		&newUser.Email,
	)
	if err != nil {
		return models.User{}, err
	}

	return newUser, err
}
