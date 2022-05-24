--  For handler /supervisor/courses
create or replace view test_db.supers_v AS
select u.id as user_id,
       c.id as course_id,
       c.course_name,
       c.semester
from test_db.users as u,
     test_db.supervisors as s,
     test_db.supervisors_courses as sc,
     test_db.courses as c
where u.id = s.user_id
  and s.id = sc.supervisor_id
  and sc.course_id = c.id;

-- Fro handler /student
create or replace view test_db.student_group_v as
select us.id as user_id,
       g.id,
       g.group_code
from test_db.users as us,
     test_db.students as st,
     test_db.groups as g
where st.user_id = us.id
  and st.group_id = g.id;

--  For handler /table
create or replace view test_db.student_event_v as
select se.id,
       se.student_id,
       cs.id as course_id,
       cs.semester,
       cs.course_name,
       se.upload_files,
       se.grade,
       evs.event_name
from test_db.student_event as se,
     test_db.students as stud,
     test_db.events as evs,
     test_db.users as usr,
     test_db.courses as cs
where se.student_id = stud.id
  and stud.user_id = usr.id
  and se.event_id = evs.id
  and evs.course_id = cs.id;

-- Get groups by course id
create or replace view test_db.get_groups_by_course_v as
select gc.group_id,
       gc.course_id,
       g.group_code,
       c.semester,
       c.course_name
from test_db.group_course as gc,
     test_db.groups as g,
     test_db.courses as c
where gc.course_id = c.id
  and gc.group_id = g.id;

-- Get students by group
create or replace view test_db.get_stud_by_gr_v as
select st.id as stud_id,
       st.user_id,
       st.group_id,
       us.firstname,
       us.lastname,
       us.middle_name,
       us.email
from test_db.students as st,
     test_db.groups as g,
     test_db.users as us
where st.group_id = g.id
  and st.user_id = us.id;