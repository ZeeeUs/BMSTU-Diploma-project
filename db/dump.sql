CREATE schema dashboard;

CREATE TABLE dashboard.users
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
CREATE TABLE dashboard.groups
(
    id         serial PRIMARY KEY NOT NULL,
    group_code varchar(25)        not null
);

CREATE TABLE dashboard.students
(
    id       serial PRIMARY KEY NOT NULL,
    user_id  int2               not null,
    group_id int2               not null
);

CREATE TABLE dashboard.events
(
    id         serial PRIMARY KEY NOT NULL,
    course_id  int                not null,
    event_name varchar(255)       not null,
    event_date date,
    deadline   date
);

CREATE TABLE dashboard.student_event
(
    id           serial PRIMARY KEY NOT NULL,
    student_id   int                not null,
    event_id     int                not null,
    upload_files text,
    grade        int,
    event_status int
);

CREATE TABLE dashboard.supervisors
(
    id      serial PRIMARY KEY NOT NULL,
    user_id int2               not null
);

CREATE TABLE dashboard.supervisors_courses
(
    course_id     int not null,
    supervisor_id int not null
);

CREATE TABLE dashboard.courses
(
    id          serial PRIMARY KEY NOT NULL,
    semester    int                not null,
    course_name varchar(255)       not null
);

CREATE TABLE dashboard.events_eventsName
(
    event_id        int not null,
    "event-name_id" int not null
);

CREATE TABLE dashboard.events_names
(
    id   serial PRIMARY KEY NOT NULL,
    name varchar(255)       not null
);

CREATE TABLE dashboard.event_status
(
    id          serial PRIMARY KEY NOT NULL,
    status_name varchar(255)       not null
);

CREATE TABLE dashboard.group_course
(
    group_id  int not null,
    course_id int not null
);

CREATE TABLE dashboard.comments
(
    id            serial primary key not null,
    comment_field text
);

alter table dashboard.students
    add foreign key (user_id) references dashboard.users (id) on delete cascade;
alter table dashboard.students
    add foreign key (group_id) references dashboard.groups (id) on delete set null;

alter table dashboard.events
    add foreign key (course_id) references dashboard.courses (id) on delete cascade;

alter table dashboard.student_event
    add foreign key (student_id) references dashboard.students (id) on delete cascade;
alter table dashboard.student_event
    add foreign key (event_id) references dashboard.events (id) on delete cascade;

alter table dashboard.supervisors
    add foreign key (user_id) references dashboard.users (id) on delete cascade;

alter table dashboard.supervisors_courses
    add foreign key (course_id) references dashboard.courses (id) on delete cascade;
alter table dashboard.supervisors_courses
    add foreign key (supervisor_id) references dashboard.supervisors (id) on delete cascade;

alter table dashboard.supervisors_courses
    add foreign key (course_id) references dashboard.courses (id) on delete cascade;
alter table dashboard.supervisors_courses
    add foreign key (supervisor_id) references dashboard.supervisors (id) on delete cascade;

alter table dashboard.events_eventsName
    add foreign key (event_id) references dashboard.events (id) on delete cascade;
alter table dashboard.events_eventsName
    add foreign key ("event-name_id") references dashboard.events_names (id) on delete cascade;

alter table dashboard.group_course
    add foreign key (group_id) references dashboard.groups (id) on delete cascade;
alter table dashboard.group_course
    add foreign key (course_id) references dashboard.courses (id) on delete cascade;

insert into dashboard.users (password, pass_status, firstname, middle_name, lastname, email, is_super) VALUES ('hello', false, 'Test', 'Test', 'Test', 'test@test', true);
insert into dashboard.users (password, pass_status, firstname, middle_name, lastname, email, is_super) VALUES ('world', false, 'Борис', 'Иванович', 'Джонсон', 'ivanovbi@bmstu.ru', false);
insert into dashboard.users (password, pass_status, firstname, middle_name, lastname, email, is_super) VALUES ('world', false, 'Валерий', 'Михайлович', 'Чёрненький', 'test2', true);


insert into dashboard.groups (group_code) values ('ИУ5-85Б');
insert into dashboard.students (user_id, group_id) VALUES (2, 1);

