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

CREATE TABLE students
(
    id       serial PRIMARY KEY NOT NULL,
    user_id  int2               not null,
    group_id int2               not null,
    foreign key (user_id) references users (id) on delete cascade,
    foreign key (group_id) references groups (id) on delete set null
);

CREATE TABLE groups
(
    id         serial
        primary key,
    group_code varchar(25) not null
);

CREATE TABLE student_event
(
    id           serial PRIMARY KEY NOT NULL,
    student_id   int                not null,
    event_id     int                not null,
    upload_files text,
    grade        int,
    event_status int,
    foreign key (student_id) references students (id) on delete cascade,
    foreign key (event_id) references events (id) on delete cascade
);

CREATE TABLE supervisors
(
    id      serial PRIMARY KEY NOT NULL,
    user_id int2               not null,
    foreign key (user_id) references users (id) on delete cascade
);

CREATE TABLE supervisors_courses
(
    course_id     int not null,
    supervisor_id int not null,
    foreign key (course_id) references courses (id) on delete cascade,
    foreign key (supervisor_id) references supervisors (id) on delete cascade
);

CREATE TABLE groups
(
    id         serial PRIMARY KEY NOT NULL,
    group_code varchar(25)        not null
);

CREATE TABLE courses
(
    id          serial PRIMARY KEY NOT NULL,
    semester    int                not null,
    course_name varchar(255)       not null
);

CREATE TABLE events
(
    id         serial PRIMARY KEY NOT NULL,
    course_id  int                not null,
    event_name varchar(255)       not null,
    event_date date,
    deadline   date,
    foreign key (course_id) references courses (id) on delete cascade
);

CREATE TABLE events_eventsName
(
    event_id        int not null,
    "event-name_id" int not null,
    foreign key (event_id) references events (id) on delete cascade,
    foreign key ("event-name_id") references events_names (id) on delete cascade
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

CREATE TABLE groups
(
    id         serial PRIMARY KEY NOT NULL,
    group_code varchar(25)
);

CREATE TABLE group_course
(
    group_id  int not null,
    course_id int not null,
    foreign key (group_id) references groups (id) on delete cascade,
    foreign key (course_id) references courses (id) on delete cascade
);

CREATE TABLE comments
(
    id            serial primary key not null,
    comment_field text
);
