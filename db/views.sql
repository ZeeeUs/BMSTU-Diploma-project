--  For handler /supervisor/courses
create or replace view dashboard.supers_v AS
select u.id as user_id,
       c.id as course_id,
       c.course_name,
       c.semester
from dashboard.users as u,
     dashboard.supervisors as s,
     dashboard.supervisors_courses as sc,
     dashboard.courses as c
where u.id = s.user_id
  and s.id = sc.supervisor_id
  and sc.course_id = c.id;

-- Fro handler /student
create or replace view dashboard.student_group_v as
select us.id as user_id,
       g.id,
       g.group_code
from dashboard.users as us,
     dashboard.students as st,
     dashboard.groups as g
where st.user_id = us.id
  and st.group_id = g.id;

--  For handler /table
create or replace view dashboard.student_event_v as
select se.id,
       se.student_id,
       cs.id as course_id,
       cs.semester,
       cs.course_name,
       se.upload_files,
       se.grade,
       evs.event_name,
       evst.status
from dashboard.student_event as se,
     dashboard.students as stud,
     dashboard.events as evs,
     dashboard.event_status as evst,
     dashboard.users as usr,
     dashboard.courses as cs
where se.student_id = stud.id
  and stud.user_id = usr.id
  and se.event_id = evs.id
  and se.event_status = evst.id
  and evs.course_id = cs.id;