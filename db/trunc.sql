truncate table test_db.comments cascade;
truncate table test_db.courses cascade;
truncate table test_db.events cascade;
truncate table test_db.group_course cascade;
truncate table test_db.groups cascade;
truncate table test_db.student_event cascade;
truncate table test_db.students cascade;
truncate table test_db.supervisors cascade;
truncate table test_db.supervisors_courses cascade;
truncate table test_db.users cascade;

alter sequence test_db.comments_id_seq restart with 1;
alter sequence test_db.courses_id_seq restart with 1;
alter sequence test_db.events_id_seq restart with 1;
alter sequence test_db.groups_id_seq restart with 1;
alter sequence test_db.students_id_seq restart with 1;
alter sequence test_db.student_event_id_seq restart with 1;
alter sequence test_db.supervisors_id_seq restart with 1;
alter sequence test_db.users_id_seq restart with 1;