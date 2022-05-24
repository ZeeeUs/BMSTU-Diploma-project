package repository

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	UpdateUser(ctx context.Context, pswd string, email string) (id int, err error)
	GetUserById(ctx context.Context, id int) (user models.User, err error)
	GetStudentId(ctt context.Context, id int) (getId int, err error)
	GetSuperId(ctt context.Context, id int) (getId int, err error)
	GetStudent(ctt context.Context, id int) (student models.Student, err error)
	GetSuper(ctt context.Context, id int) (supervisor models.Supervisor, err error)
}

type userRepository struct {
	conn   *pgx.ConnPool
	logger *logrus.Logger
}

func NewUserRepository(pgConn *pgx.ConnPool, logger *logrus.Logger) UserRepository {
	return &userRepository{
		conn:   pgConn,
		logger: logger,
	}
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	err = u.conn.QueryRow("select id, password, pass_status, firstname, middle_name, lastname, email, is_super from test_db.users"+
		" where email=$1", email).Scan(
		&user.Id,
		&user.Password,
		&user.PassStatus,
		&user.Firstname,
		&user.MiddleName,
		&user.Lastname,
		&user.Email,
		&user.IsSuper,
	)

	if err != nil {
		return models.User{}, err
	}

	return
}

func (u *userRepository) UpdateUser(ctx context.Context, pswd string, email string) (id int, err error) {
	err = u.conn.QueryRow("update test_db.users"+
		" set password=$1, pass_status=true"+
		" where email=$2 returning id", pswd, email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return
}

func (u *userRepository) GetUserById(ctx context.Context, id int) (user models.User, err error) {
	err = u.conn.QueryRow("select id, password, pass_status, firstname, middle_name, lastname, email, is_super from test_db.users"+
		" where id=$1", id).Scan(
		&user.Id,
		&user.Password,
		&user.PassStatus,
		&user.Firstname,
		&user.MiddleName,
		&user.Lastname,
		&user.Email,
		&user.IsSuper,
	)

	if err != nil {
		return models.User{}, err
	}

	return
}

func (u *userRepository) GetStudentId(ctt context.Context, id int) (getId int, err error) {
	err = u.conn.QueryRow("select user_id from test_db.students where user_id=$1", id).Scan(&getId)
	if err != nil {
		return 0, err
	}

	return
}

func (u *userRepository) GetStudent(ctt context.Context, id int) (student models.Student, err error) {
	err = u.conn.QueryRow("select id, user_id, group_id from test_db.students where user_id=$1", id).Scan(
		&student.Id,
		&student.UserId,
		&student.GroupId,
	)
	if err != nil {
		return models.Student{}, err
	}

	return
}

func (u *userRepository) GetSuper(ctt context.Context, id int) (supervisor models.Supervisor, err error) {
	err = u.conn.QueryRow("select id, user_id from test_db.supervisors where user_id=$1", id).Scan(
		&supervisor.Id,
		&supervisor.UserId,
	)
	if err != nil {
		return models.Supervisor{}, err
	}

	return
}

func (u *userRepository) GetSuperId(ctt context.Context, id int) (getId int, err error) {
	err = u.conn.QueryRow("select user_id from test_db.supervisors where user_id=$1", id).Scan(&getId)
	if err != nil {
		return 0, err
	}

	return
}
