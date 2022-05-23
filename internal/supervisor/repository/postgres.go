package repository

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type SupersRepository interface {
	GetSupersCourses(ctx context.Context, id int) ([]models.Course, error)
}

type supersRepo struct {
	conn   *pgx.ConnPool
	logger *logrus.Logger
}

func NewSupersRepository(pgConn *pgx.ConnPool, log *logrus.Logger) SupersRepository {
	return &supersRepo{
		conn:   pgConn,
		logger: log,
	}
}

func (su *supersRepo) GetSupersCourses(ctx context.Context, id int) ([]models.Course, error) {
	rows, err := su.conn.Query("select course_id, course_name, semester from test_db.supers_v where user_id=$1", id)
	defer rows.Close()

	var (
		courseId   int
		courseName string
		semester   int
	)

	courses := make([]models.Course, 0)
	for rows.Next() {
		err = rows.Scan(&courseId, &courseName, &semester)
		if err != nil {
			su.logger.Errorf("GetSupersCourses: can't scan from rows to vars: %s", err)
			return nil, err
		}

		courses = append(courses, models.Course{
			Id:         courseId,
			Semester:   semester,
			CourseName: courseName,
		})
	}
	return courses, nil
}
