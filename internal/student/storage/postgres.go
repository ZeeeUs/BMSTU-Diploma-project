package storage

import (
	"context"
	"sort"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type StudentStorage interface {
	GetUserGroup(ctx context.Context, id int) (models.Group, int, error)
	GetTable(ctx context.Context, id int) (table models.Table, err error)
	GetGroup(ctx context.Context, id int) (group models.Group, err error)
	AddFileName(ctx context.Context, fileName string, studentEventId int) error
	ChangeEventStatus(ctx context.Context, status int, studEventId int) error
}

type studentStorage struct {
	conn   *pgx.ConnPool
	logger *logrus.Logger
}

func NewStudentRepository(pgConn *pgx.ConnPool, logger *logrus.Logger) StudentStorage {
	return &studentStorage{
		conn:   pgConn,
		logger: logger,
	}
}

func (sr *studentStorage) GetUserGroup(ctx context.Context, id int) (group models.Group, studentId int, err error) {
	err = sr.conn.QueryRow("select student_id, id, group_code from test_db.student_group_v where user_id=$1", id).Scan(
		&studentId,
		&group.Id,
		&group.GroupCode,
	)

	if err != nil {
		return models.Group{}, 0, err
	}

	return
}

func (sr *studentStorage) GetTable(ctx context.Context, id int) (models.Table, error) {
	var (
		cId     int
		cName   string
		tbl     models.Table
		names   []string
		cCourse []models.CCourse
	)

	doubleString := make(map[int]string, 0)

	rows, err := sr.conn.Query("select courses.id, courses.course_name from test_db.courses, test_db.students"+
		" where students.id=$1", id)
	if err != nil {
		return models.Table{}, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&cId, &cName)

		doubleString[cId] = cName
		names = append(names, cName)
	}

	for key, val := range doubleString {
		var rrows *pgx.Rows
		rrows, err = sr.conn.Query("select student_event.id,"+
			" events.event_date,"+
			" events.deadline,"+
			" student_event.event_status,"+
			" events.event_name,"+
			" events.description,"+
			" student_event.upload_files"+
			" from test_db.events, test_db.student_event"+
			" where course_id=$1 and events.id=student_event.event_id and student_event.student_id=$2"+
			" order by event_date", key, id)
		if err != nil {
			return models.Table{}, err
		}

		var events []models.Event
		for rrows.Next() {
			var event models.Event
			err = rrows.Scan(
				&event.EventId,
				&event.EventDate,
				&event.Deadline,
				&event.Status,
				&event.EventName,
				&event.Description,
				&event.Files,
			)

			if err != nil {
				sr.logger.Error(err)
			}

			events = append(events, event)
		}

		if events != nil {
			cCourse = append(cCourse, models.CCourse{
				CourseId:   key,
				CourseName: val,
				Events:     events,
			})
		}

		rrows.Close()
	}

	sort.Strings(names)
	for _, v := range names {
		for _, val := range cCourse {
			if v == val.CourseName {
				tbl.Courses = append(tbl.Courses, val)
			}
		}
	}

	return tbl, nil
}

func (sr *studentStorage) GetGroup(ctx context.Context, id int) (group models.Group, err error) {
	err = sr.conn.QueryRow("select id, group_code from test_db.student_group_v where user_id=$1", id).Scan(
		&group.Id,
		&group.GroupCode,
	)

	if err != nil {
		return models.Group{}, err
	}

	return
}

func (sr *studentStorage) AddFileName(ctx context.Context, fileName string, studentEventId int) error {
	_, err := sr.conn.Exec("update test_db.student_event set upload_files=array_append(upload_files, $1)"+
		" where id=$2", fileName, studentEventId)
	if err != nil {
		return err
	}

	return nil
}

func (sr *studentStorage) ChangeEventStatus(ctx context.Context, status int, studEventId int) error {
	_, err := sr.conn.Exec("update test_db.student_event set event_status=$1 where id=$2", status, studEventId)
	if err != nil {
		return err
	}

	sr.logger.Infof("event_status = %d\nid = %d\n", status, studEventId)

	return nil
}
