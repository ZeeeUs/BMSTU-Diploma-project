CREATE SCHEMA test_db;

CREATE TABLE test_db.users
(
    id          serial PRIMARY KEY NOT NULL,
    password    text               NOT NULL,
    pass_status bool,
    firstname   varchar(255)       NOT NULL,
    middle_name varchar(255)       NOT NULL,
    lastname    varchar(255)       NOT NULL,
    email       varchar unique     not null,
    is_super    boolean            not null
);
CREATE TABLE test_db.groups
(
    id         serial PRIMARY KEY NOT NULL,
    group_code varchar(25)        not null
);

CREATE TABLE test_db.students
(
    id       serial PRIMARY KEY NOT NULL,
    user_id  int2               not null,
    group_id int2               not null
);

CREATE TABLE test_db.events
(
    id          serial PRIMARY KEY NOT NULL,
    course_id   int                not null,
    event_name  varchar(255)       not null,
    event_date  date,
    deadline    date,
    description text
);

CREATE TABLE test_db.student_event
(
    id           serial PRIMARY KEY NOT NULL,
    student_id   int                not null,
    event_id     int                not null,
    upload_files text,
    grade        int,
    event_status int
);

CREATE TABLE test_db.supervisors
(
    id      serial PRIMARY KEY NOT NULL,
    user_id int2               not null
);

CREATE TABLE test_db.supervisors_courses
(
    course_id     int not null,
    supervisor_id int not null
);

CREATE TABLE test_db.courses
(
    id          serial PRIMARY KEY NOT NULL,
    semester    int                not null,
    course_name varchar(255)       not null
);

CREATE TABLE test_db.group_course
(
    group_id  int not null,
    course_id int not null
);

CREATE TABLE test_db.comments
(
    id               serial primary key not null,
    student_event_id int                not null,
    comment_field    text
);

alter table test_db.students
    add foreign key (user_id) references test_db.users (id) on delete cascade;
alter table test_db.students
    add foreign key (group_id) references test_db.groups (id) on delete set null;

alter table test_db.events
    add foreign key (course_id) references test_db.courses (id) on delete cascade;

alter table test_db.student_event
    add foreign key (student_id) references test_db.students (id) on delete cascade;
alter table test_db.student_event
    add foreign key (event_id) references test_db.events (id) on delete cascade;

alter table test_db.supervisors
    add foreign key (user_id) references test_db.users (id) on delete cascade;

alter table test_db.supervisors_courses
    add foreign key (course_id) references test_db.courses (id) on delete cascade;
alter table test_db.supervisors_courses
    add foreign key (supervisor_id) references test_db.supervisors (id) on delete cascade;

alter table test_db.supervisors_courses
    add foreign key (course_id) references test_db.courses (id) on delete cascade;
alter table test_db.supervisors_courses
    add foreign key (supervisor_id) references test_db.supervisors (id) on delete cascade;

alter table test_db.group_course
    add foreign key (group_id) references test_db.groups (id) on delete cascade;
alter table test_db.group_course
    add foreign key (course_id) references test_db.courses (id) on delete cascade;

alter table test_db.comments
    add foreign key (student_event_id) references test_db.student_event (id) on delete cascade;