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

select * from dashboard.student_group_v where user_id=2;