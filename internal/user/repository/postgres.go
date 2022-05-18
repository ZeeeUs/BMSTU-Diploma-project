package repository

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/config"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	UpdateUser(ctx context.Context, pswd string, email string) (id int, err error)
	GetUserById(ctx context.Context, id int) (user models.User, err error)
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

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	err = u.Conn.QueryRow("select id, password, pass_status, firstname, middle_name, lastname, email from dashboard.users"+
		" where email=$1", email).Scan(
		&user.Id,
		&user.Password,
		&user.PassStatus,
		&user.Firstname,
		&user.MiddleName,
		&user.Lastname,
		&user.Email,
	)

	if err != nil {
		return models.User{}, err
	}

	return
}

func (u *userRepository) UpdateUser(ctx context.Context, pswd string, email string) (id int, err error) {
	err = u.Conn.QueryRow("update dashboard.users"+
		" set password=$1"+
		" where email=$2 returning id", pswd, email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return
}

func (u *userRepository) GetUserById(ctx context.Context, id int) (user models.User, err error) {
	err = u.Conn.QueryRow("select id, password, pass_status, firstname, middle_name, lastname, email from dashboard.users"+
		" where id=$1", id).Scan(
		&user.Id,
		&user.Password,
		&user.PassStatus,
		&user.Firstname,
		&user.MiddleName,
		&user.Lastname,
		&user.Email,
	)

	if err != nil {
		return models.User{}, err
	}

	return
}