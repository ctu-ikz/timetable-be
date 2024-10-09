INSERT INTO "Semester" ("codename", "start", "end")
VALUES (  'B241', '2024-09-23', '2025-01-10');

INSERT INTO "Subject" ("name", "shortcut", "code_name")
VALUES ('Právní předpisy ve zdravotnictví a normy', 'PP', 'F7PBKPPN');

INSERT INTO "Field" ("name", "shortcut")
VALUES ('Informatika a kybernetika ve zdravotnictví', 'IKZ');

insert into "Class" (start_year, end_year, field_id)
values (2022,2025,1);

insert into "SubjectClassType" ("name")
values ('Přednáška');

insert into "SubjectClass" (semester_id, subject_id, class, start_time, end_time, day, type)
values (1,1,1,'10:00','11:50',1, 1);

insert into "SubjectWeek" (week_number, subject_class_id) VALUES (1,1);
insert into "SubjectWeek" (week_number, subject_class_id) VALUES (3,1);
insert into "SubjectWeek" (week_number, subject_class_id) VALUES (5,1);
insert into "SubjectWeek" (week_number, subject_class_id) VALUES (7,1);
insert into "SubjectWeek" (week_number, subject_class_id) VALUES (9,1);
insert into "SubjectWeek" (week_number, subject_class_id) VALUES (11,1);
insert into "SubjectWeek" (week_number, subject_class_id) VALUES (13,1);