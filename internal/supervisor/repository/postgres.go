package repository

import (
	"context"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/models"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type SupersRepository interface {
	GetSupersCourses(ctx context.Context, id int) ([]models.Course, error)
	GetSuperId(ctx context.Context, id int) (int, error)
	GetGroupsByCourseId(ctx context.Context, id int) ([]models.GroupByCourse, error)
	GetStudentsByGroup(ctx context.Context, groupId int) ([]models.StudentByGroup, error)
	GetEventsByCourseId(ctx context.Context, id int) ([]models.Event, error)
	GetStudentEvents(ctx context.Context, studentId int, courseId int) ([]models.StudentEvent, error)
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

func (su *supersRepo) GetSuperId(ctx context.Context, id int) (int, error) {
	var getId int
	err := su.conn.QueryRow("select id from test_db.supervisors where user_id=$1", id).Scan(&getId)
	su.logger.Info(id)
	if err != nil {
		return 0, err
	}

	return getId, nil
}

func (su *supersRepo) GetGroupsByCourseId(ctx context.Context, id int) ([]models.GroupByCourse, error) {
	rows, err := su.conn.Query("select * from test_db.get_groups_by_course_v where course_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.GroupByCourse
	for rows.Next() {
		var group models.GroupByCourse

		err := rows.Scan(
			&group.GroupId,
			&group.CourseId,
			&group.GroupCode,
			&group.Semester,
			&group.CourseName,
		)
		if err != nil {
			su.logger.Errorf("GetGroupsByCourseId: can't scan object: %s", err)
			continue
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (su *supersRepo) GetStudentsByGroup(ctx context.Context, groupId int) ([]models.StudentByGroup, error) {
	rows, err := su.conn.Query("select stud_id, user_id, group_id, firstname, lastname, middle_name, email"+
		" from test_db.get_stud_by_gr_v where group_id=$1", groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.StudentByGroup
	for rows.Next() {
		var student models.StudentByGroup

		err := rows.Scan(
			&student.StudentId,
			&student.UserId,
			&student.GroupId,
			&student.FirstName,
			&student.LastName,
			&student.MiddleName,
			&student.Email,
		)
		if err != nil {
			su.logger.Errorf("GetStudentsByGroup: can't scan object: %s", err)
			continue
		}

		students = append(students, student)
	}

	return students, nil
}

func (su *supersRepo) GetEventsByCourseId(ctx context.Context, id int) ([]models.Event, error) {
	rows, err := su.conn.Query("select id, event_name, event_date, deadline, description from test_db.events where course_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.EventId,
			&event.EventName,
			&event.EventDate,
			&event.Deadline,
			&event.Description,
		)
		if err != nil {
			su.logger.Errorf("GetGroupsByCourseId: can't scan object: %s", err)
			continue
		}

		events = append(events, event)
	}

	return events, nil
}

func (su *supersRepo) GetStudentEvents(ctx context.Context, studentId int, courseId int) ([]models.StudentEvent, error) {
	rows, err := su.conn.Query("select student_event.id,"+
		" student_event.student_id,"+
		" student_event.event_id,"+
		" student_event.upload_files,"+
		" student_event.event_status"+
		" from test_db.students,"+
		" test_db.student_event,"+
		" test_db.events,"+
		" test_db.courses"+
		" where student_event.student_id = students.id"+
		" and student_event.event_id = events.id"+
		" and events.course_id = courses.id"+
		" and student_id = $1"+
		" and course_id = $2", studentId, courseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.StudentEvent
	for rows.Next() {
		var event models.StudentEvent

		err := rows.Scan(
			&event.Id,
			&event.StudentId,
			&event.EventId,
			&event.UploadFiles,
			&event.Status,
		)
		if err != nil {
			su.logger.Errorf("GetStudentEvents: can't scan object: %s", err)
			continue
		}

		events = append(events, event)
	}

	return events, nil
}
