from random import randint

import psycopg2
from faker import Faker
from psycopg2 import Error

fake = Faker('ru_RU')

STUDENTS_COUNT = 100
SUPERUSERS_COUNT = 10
groupsName = ["ИУ5-71Б", "ИУ5-72Б", "ИУ5-83Б", "ИУ5-84Б", "ИУ5-85Б"]
coursesName = [{"semester": 8, "name": "Беспроводные сети"},
               {"semester": 8, "name": "Защита информации"},
               {"semester": 8, "name": "Русский язык делового общения"},
               {"semester": 8, "name": "Экономика часть 2"},
               {"semester": 8, "name": "Эксплуатация АСОИУ"},
               {"semester": 7, "name": "Методы проектирования АСОИУ"},
               {"semester": 7, "name": "Элементы управления в АСОиУ"},
               {"semester": 7, "name": "Имитационное моделирование дискретных процессов"},
               ]


# Groups
def fillGroups(connection):
    cursor = connection.cursor()
    groups = []
    for i in range(len(groupsName)):
        group = {"groupCode": groupsName[i]}
        insert_query = """ INSERT INTO test_db.groups (group_code)
            VALUES (%s) returning id;"""
        cursor.execute(insert_query, (group["groupCode"],))
        connection.commit()
        id = cursor.fetchone()[0]
        groups.append({
            "id": id,
            "code": groupsName[i],
        })

    cursor.close()
    return groups


# Students
def fillStudents(connection, groups):
    # usersStudents = []
    cursor = connection.cursor()
    students = []
    for i in range(STUDENTS_COUNT):
        user = {
            "pass": "qwe",
            "passStatus": False,
            "fName": fake.first_name(),
            "mName": fake.middle_name(),
            "lName": fake.last_name(),
            "email": f"stud{i}",
            "isSuper": False,
        }
        # вставка юзера в бд получение айдишника, пока хардкод айдишника
        insert_query = """ INSERT INTO test_db.users (password, pass_status, firstname, middle_name, lastname, email, is_super)
            VALUES (%s, %s, %s, %s, %s, %s, %s) returning id;"""
        cursor.execute(insert_query, (
            user["pass"], user["passStatus"], user["fName"], user["mName"], user["lName"], user["email"],
            user["isSuper"],))
        connection.commit()
        id = cursor.fetchone()[0]
        user["id"] = id
        # usersStudents.append(user)

        group = len(students) % len(groups)
        student = {
            "userId": user["id"],
            "groupId": groups[group]["id"],
        }
        # вставка студента в бд
        insert_query = """ INSERT INTO test_db.students (user_id, group_id)
            VALUES (%s, %s) returning id;"""
        cursor.execute(insert_query, (student["userId"], student["groupId"],))
        connection.commit()
        id = cursor.fetchone()[0]

        student["id"] = id  # хардкод айдишника
        students.append(student)

    cursor.close()
    return students


# Superuser
def fillSuperuser(connection):
    # usersSuper = []
    cursor = connection.cursor()
    superusers = []
    for i in range(SUPERUSERS_COUNT):
        user = {
            "pass": "qwe",
            "passStatus": False,
            "fName": fake.first_name(),
            "mName": fake.middle_name(),
            "lName": fake.last_name(),
            "email": f"super{i}",
            "isSuper": True,
        }
        # вставка юзера в бд получение айдишника, пока хардкод айдишника
        insert_query = """ INSERT INTO test_db.users (password, pass_status, firstname, middle_name, lastname, email, is_super)
            VALUES (%s, %s, %s, %s, %s, %s, %s) returning id;"""
        cursor.execute(insert_query, (
            user["pass"], user["passStatus"], user["fName"], user["mName"], user["lName"], user["email"],
            user["isSuper"],))
        connection.commit()
        id = cursor.fetchone()[0]
        user["id"] = id

        # вставка препода
        superuser = {
            "userId": user["id"],
        }

        insert_query = """ INSERT INTO test_db.supervisors (user_id)
            VALUES (%s) returning id;"""
        cursor.execute(insert_query, (superuser["userId"],))
        connection.commit()
        id = cursor.fetchone()[0]

        superuser["id"] = id  # хардкод айдишника
        superusers.append(superuser)

    cursor.close()
    return superusers


# Courses
def fillCourses(connection):
    cursor = connection.cursor()
    courses = []
    for i in range(len(coursesName)):
        course = coursesName[i]
        # добавляем курс в бд, получаем айдишник
        insert_query = """ INSERT INTO test_db.courses (semester, course_name)
            VALUES (%s, %s) returning id;"""
        cursor.execute(insert_query, (course["semester"], course["name"],))
        connection.commit()
        id = cursor.fetchone()[0]
        course["id"] = id
        courses.append(course)

    cursor.close()
    return courses


# Supervisor courses
def fillSupervisorCourses(connection, courses, supervisors):
    cursor = connection.cursor()
    i = 0
    for course in courses:
        supervisorCourse = {
            "courseId": course["id"],
            "supervisorId": supervisors[i % SUPERUSERS_COUNT]["id"],
        }
        # вставка в бд
        insert_query = """ INSERT INTO test_db.supervisors_courses (course_id, supervisor_id)
            VALUES (%s, %s);"""
        cursor.execute(insert_query, (supervisorCourse["courseId"], supervisorCourse["supervisorId"],))
        connection.commit()
        i += 1

    cursor.close()


# Group courses
def fillGroupCourses(connection, courses, groups):
    cursor = connection.cursor()
    groupCourses = []
    for course in courses:
        if course["semester"] == 8:
            for i in range(2, 5):
                groupCourse = {
                    "courseId": course["id"],
                    "groupId": groups[i]["id"]
                }
                insert_query = """ INSERT INTO test_db.group_course (group_id, course_id)
                VALUES (%s,%s);"""
                cursor.execute(insert_query, (groupCourse["groupId"], groupCourse["courseId"],))
                connection.commit()
                groupCourses.append(groupCourse)

        if course["semester"] == 7:
            for i in range(0, 2):
                groupCourse = {
                    "courseId": course["id"],
                    "groupId": groups[i]["id"]
                }
                insert_query = """ INSERT INTO test_db.group_course (group_id, course_id)
                VALUES (%s,%s);"""
                cursor.execute(insert_query, (groupCourse["groupId"], groupCourse["courseId"],))
                connection.commit()
                groupCourses.append(groupCourse)

    cursor.close()
    return groupCourses


# Events
def fillEvents(connection, courses):
    cursor = connection.cursor()
    eventsCourses = []
    for course in courses:
        eventsCourse = {
            "courseId": course["id"],
        }
        eventCount = randint(3, 7)

        events = []
        for i in range(eventCount):
            event = {
                "courseId": course["id"],
                "eventName": f"Лаб {i + 1}",
                "eventDate": f"2022-05-{1+i}",
                "deadline": f"2022-06-{1+i}",
                "description": fake.text(),
            }
            # вставляем в бд добавляем айдишник
            insert_query = """ INSERT INTO test_db.events (course_id, event_name, event_date, deadline, description)
            VALUES (%s, %s, %s, %s, %s) returning id;"""
            cursor.execute(insert_query, (event["courseId"], event["eventName"], event["eventDate"], event["deadline"], event["description"]))
            connection.commit()
            id = cursor.fetchone()[0]
            event["id"] = id
            events.append(event)

        eventsCourse["events"] = events
        eventsCourses.append(eventsCourse)

    cursor.close()
    return eventsCourses


# Student events
def fillStudentEvents(connection, eventCourses, groupCourses, students):
    cursor = connection.cursor()
    for course in eventCourses:
        for group in groupCourses:
            if group["courseId"] == course["courseId"]:
                for student in students:
                    if student["groupId"] == group["groupId"]:
                        # студент из группы которая учавствует в курсе, добавляем для него все эвенты
                        for ev in course["events"]:
                            eventStudent = {
                                "studentId": student["id"],
                                "eventStatus": randint(0, 6),
                                "eventId": ev["id"],
                            }
                            # добавляем в бд получаем айдишник
                            # eventStudent["id"] = 0
                            insert_query = """ INSERT INTO test_db.student_event (student_id, event_id, event_status)
                            VALUES (%s, %s, %s) returning id;"""
                            cursor.execute(insert_query, (
                                eventStudent["studentId"], eventStudent["eventId"], eventStudent["eventStatus"]))
                            connection.commit()
    cursor.close()


try:
    # Подключиться к существующей базе данных
    print(1)
    connection = psycopg2.connect(user="buser",
                                  # пароль, который указали при установке PostgreSQL
                                  password="bpassword",
                                  host="127.0.0.1",
                                  port="5432",
                                  database="bdb")
    # connection = psycopg2.connect(user="bmstuUser",
    #                               # пароль, который указали при установке PostgreSQL
    #                               password="pgpwd4bmstu",
    #                               host="127.0.0.1",
    #                               port="5432",
    #                               database="bmstuDb")

    groups = fillGroups(connection)
    students = fillStudents(connection, groups)
    supers = fillSuperuser(connection)
    courses = fillCourses(connection)
    fillSupervisorCourses(connection, courses, supers)
    eventsCourses = fillEvents(connection, courses)
    groupCourses = fillGroupCourses(connection, courses, groups)
    fillStudentEvents(connection, eventsCourses, groupCourses, students)
    print("1 элемент успешно добавлен")

except (Exception, Error) as error:
    print("Ошибка при работе с PostgreSQL", error)
