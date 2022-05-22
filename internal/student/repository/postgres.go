package repository

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type StudentRepository interface {
	GetUserGroup(ctx context.Context, id int) (models.Group, error)
}

type studentRepository struct {
	conn   *pgx.ConnPool
	logger *logrus.Logger
}

func NewStudentRepository(pgConn *pgx.ConnPool, logger *logrus.Logger) StudentRepository {
	return &studentRepository{
		conn:   pgConn,
		logger: logger,
	}
}

func (sr *studentRepository) GetUserGroup(ctx context.Context, id int) (group models.Group, err error) {
	err = sr.conn.QueryRow("select id, group_code from dashboard.student_group_v where user_id=$1", id).Scan(
		&group.Id,
		&group.GroupCode,
	)

	if err != nil {
		return models.Group{}, err
	}

	return
}
