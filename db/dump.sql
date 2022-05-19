CREATE TABLE users
(
    id          serial PRIMARY KEY NOT NULL,
    password    varchar(50)        NOT NULL,
    pass_status bool,
    firstname   varchar(255)       NOT NULL,
    middle_name varchar(255)       NOT NULL,
    lastname    varchar(255)       NOT NULL,
    email       varchar unique     not null
);
CREATE TABLE groups
(
    id         serial PRIMARY KEY NOT NULL,
    group_code varchar(25)        not null
);

CREATE TABLE students
(
    id       serial PRIMARY KEY NOT NULL,
    user_id  int2               not null,
    group_id int2               not null
);

CREATE TABLE events
(
    id         serial PRIMARY KEY NOT NULL,
    course_id  int                not null,
    event_name varchar(255)       not null,
    event_date date,
    deadline   date
);

CREATE TABLE student_event
(
    id           serial PRIMARY KEY NOT NULL,
    student_id   int                not null,
    event_id     int                not null,
    upload_files text,
    grade        int,
    event_status int
);

CREATE TABLE supervisors
(
    id      serial PRIMARY KEY NOT NULL,
    user_id int2               not null
);

CREATE TABLE supervisors_courses
(
    course_id     int not null,
    supervisor_id int not null
);

CREATE TABLE courses
(
    id          serial PRIMARY KEY NOT NULL,
    semester    int                not null,
    course_name varchar(255)       not null
);

CREATE TABLE events_eventsName
(
    event_id        int not null,
    "event-name_id" int not null
);

CREATE TABLE events_names
(
    id   serial PRIMARY KEY NOT NULL,
    name varchar(255)       not null
);

CREATE TABLE event_status
(
    id          serial PRIMARY KEY NOT NULL,
    status_name varchar(255)       not null
);

CREATE TABLE group_course
(
    group_id  int not null,
    course_id int not null
);

CREATE TABLE comments
(
    id            serial primary key not null,
    comment_field text
);

alter table students add foreign key (user_id) references users (id) on delete cascade;
alter table students add   foreign key (group_id) references groups (id) on delete set null;

alter table events add foreign key (course_id) references courses (id) on delete cascade;

alter table student_event add foreign key (student_id) references students (id) on delete cascade;
alter table student_event add foreign key (event_id) references events (id) on delete cascade;

alter table supervisors add foreign key (user_id) references users (id) on delete cascade;

alter table supervisors_courses add foreign key (course_id) references courses (id) on delete cascade;
alter table supervisors_courses add foreign key (supervisor_id) references supervisors (id) on delete cascade;

alter table supervisors_courses add foreign key (course_id) references courses (id) on delete cascade;
alter table supervisors_courses add foreign key (supervisor_id) references supervisors (id) on delete cascade;

alter table events_eventsName add foreign key (event_id) references events (id) on delete cascade;
alter table events_eventsName add foreign key ("event-name_id") references events_names (id) on delete cascade;

alter table group_course add foreign key (group_id) references groups (id) on delete cascade;
alter table group_course add foreign key (course_id) references courses (id) on delete cascade;