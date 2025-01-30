CREATE TABLE "users" (
    "id" bigint primary key,
    "username" varchar(64) not null unique,
    "password" varchar(255) not null,
    "school_id" bigint
);

CREATE TABLE "schools" (
    "id" bigint primary key,
    "name" varchar(64) unique not null,
    "shortcut" varchar(32) unique
);

CREATE TABLE "faculties" (
    "id" bigserial not null primary key,
    "school_id" bigint not null,
    "name" varchar(64) not null,
    "shortcut" varchar(64)
);

CREATE TABLE "fields" (
    "id" bigserial not null primary key,
    "faculty_id" bigint not null,
    "name" varchar(64) not null,
    "shortcut" char(4)
);

CREATE TABLE "subjects" (
    "id" bigserial not null primary key,
    "faculty_id" bigint not null,
    "name" varchar(64) not null,
    "shortcut" char(4) not null,
    "code_name" varchar(32)
);

CREATE TABLE "classes" (
    "id" bigserial not null primary key,
    "field_id" bigint not null,
    "start_year" smallint not null,
    "end_year" smallint not null
);

CREATE TABLE "semesters" (
    "id" bigserial not null primary key,
    "faculty_id" bigint not null,
    "start" date not null,
    "end" date not null
);

CREATE TABLE "subject_instance_types" (
    "id" bigserial not null primary key,
    "school_id" bigint not null,
    "name" varchar(32) not null,
    "description" varchar(128) not null
);

CREATE TABLE "subject_instances" (
    "id" bigserial not null primary key,
    "semester_id" bigint not null,
    "subject_id" bigint not null,
    "sit_id" bigint not null,
    "start_time" time not null,
    "room" varchar(32) not null,
    "end_time" time not null,
    "day" smallint not null CHECK ("day" BETWEEN 0 AND 6)
);

CREATE TABLE "si_to_users" (
    "id" bigserial not null primary key,
    "user_id" bigint not null,
    "si_id" bigint not null
);

CREATE TABLE "subject_weeks" (
    "id" bigserial not null primary key,
    "si_id" bigint not null,
    "week" smallint not null CHECK ("week" BETWEEN 1 AND 200)
);

CREATE TABLE "notes" (
    "id" bigserial not null primary key,
    "school_id" bigint not null,
    "data" text,
    "user_id" bigint,
    "subject_id" bigint
);

CREATE TABLE "one_offs" (
    "id" bigserial not null primary key,
    "school_id" bigint not null,
    "name" varchar(32) not null,
    "description" varchar(128) not null,
    "day" smallint not null CHECK ("day" BETWEEN 0 AND 6),
    "start_time" time not null,
    "end_time" time not null,
    "room" varchar(32),
    "class_id" bigint
);

CREATE TABLE "one_off_to_users" (
    "id" bigserial not null primary key,
    "user_id" bigint not null,
    "one_off_id" bigint not null
);

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY, 
    token TEXT NOT NULL, 
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    expires_at TIMESTAMP NOT NULL, 
    created_at TIMESTAMP DEFAULT NOW(), 
    is_valid BOOLEAN DEFAULT TRUE
);


ALTER TABLE "users" ADD CONSTRAINT "u_s_fk" FOREIGN KEY ("school_id") REFERENCES "schools" ("id") ON DELETE SET NULL;
ALTER TABLE "faculties" ADD CONSTRAINT "fa_s_fk" FOREIGN KEY ("school_id") REFERENCES "schools" ("id") ON DELETE CASCADE;
ALTER TABLE "fields" ADD CONSTRAINT "fi_fa_fk" FOREIGN KEY ("faculty_id") REFERENCES "faculties" ("id") ON DELETE CASCADE;
ALTER TABLE "subjects" ADD CONSTRAINT "su_fa_fk" FOREIGN KEY ("faculty_id") REFERENCES "faculties" ("id") ON DELETE CASCADE;
ALTER TABLE "classes" ADD CONSTRAINT "c_fi_fk" FOREIGN KEY ("field_id") REFERENCES "fields" ("id") ON DELETE CASCADE;
ALTER TABLE "semesters" ADD CONSTRAINT "se_fa_fk" FOREIGN KEY ("faculty_id") REFERENCES "faculties" ("id") ON DELETE CASCADE;
ALTER TABLE "subject_instance_types" ADD CONSTRAINT "sit_s_fk" FOREIGN KEY ("school_id") REFERENCES "schools" ("id") ON DELETE CASCADE;
ALTER TABLE "subject_instances" ADD CONSTRAINT "si_sit_fk" FOREIGN KEY ("sit_id") REFERENCES "subject_instance_types" ("id") ON DELETE CASCADE;
ALTER TABLE "subject_instances" ADD CONSTRAINT "si_se_fk" FOREIGN KEY ("semester_id") REFERENCES "semesters" ("id") ON DELETE CASCADE;
ALTER TABLE "subject_instances" ADD CONSTRAINT "si_su_fk" FOREIGN KEY ("subject_id") REFERENCES "subjects" ("id") ON DELETE CASCADE;
ALTER TABLE "si_to_users" ADD CONSTRAINT "stu_u_fk" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "si_to_users" ADD CONSTRAINT "stu_si_fk" FOREIGN KEY ("si_id") REFERENCES "subject_instances" ("id") ON DELETE CASCADE;
ALTER TABLE "subject_weeks" ADD CONSTRAINT "sw_si_fk" FOREIGN KEY ("si_id") REFERENCES "subject_instances" ("id") ON DELETE CASCADE;
ALTER TABLE "notes" ADD CONSTRAINT "n_s_fk" FOREIGN KEY ("school_id") REFERENCES "schools" ("id") ON DELETE CASCADE;
ALTER TABLE "notes" ADD CONSTRAINT "n_u_fk" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE SET NULL;
ALTER TABLE "notes" ADD CONSTRAINT "n_su_fk" FOREIGN KEY ("subject_id") REFERENCES "subjects" ("id") ON DELETE CASCADE;
ALTER TABLE "one_offs" ADD CONSTRAINT "oo_u_fk" FOREIGN KEY ("school_id") references "schools" ("id") ON DELETE CASCADE;
ALTER TABLE "one_offs" ADD CONSTRAINT "oo_c_fk" FOREIGN KEY ("class_id") references "classes" ("id") ON DELETE CASCADE;
ALTER TABLE "one_off_to_users" ADD CONSTRAINT "ootu_u_fk" FOREIGN KEY ("user_id") references "users" ("id") ON DELETE CASCADE;
ALTER TABLE "one_off_to_users" ADD CONSTRAINT "ootu_oo_fk" FOREIGN KEY ("one_off_id") references "one_offs" ("id") ON DELETE CASCADE;

CREATE INDEX idx_faculty_school_id ON faculties(school_id);
CREATE INDEX idx_field_faculty_id ON fields(faculty_id);
CREATE INDEX idx_subject_faculty_id ON subjects(faculty_id);
CREATE INDEX idx_class_field_id ON classes(field_id);
CREATE INDEX idx_semester_faculty_id ON semesters(faculty_id);
CREATE INDEX idx_subject_instance_semester_id ON subject_instances(semester_id);
CREATE INDEX idx_subject_instance_subject_id ON subject_instances(subject_id);
CREATE INDEX idx_si_to_user_user_id ON si_to_users(user_id);
CREATE INDEX idx_si_to_user_si_id ON si_to_users(si_id);
CREATE INDEX idx_subject_week_si_id ON subject_weeks(si_id);
CREATE INDEX idx_note_school_id ON notes(school_id);
CREATE INDEX idx_note_user_id ON notes(user_id);
CREATE INDEX idx_note_subject_id ON notes(subject_id);

INSERT INTO "schools" ("id", "name", "shortcut") VALUES (1, 'České vysoké učení technické', 'ČVUT');

INSERT INTO "faculties" ("id", "school_id", "name", "shortcut") VALUES (1, 1, 'Fakulta biomedicínského inženýrství', 'FBMI');

INSERT INTO "fields" ("faculty_id", "name", "shortcut") VALUES (1, 'Informatika a kybernetika ve zdravotnictví', 'IKZ');

INSERT INTO "semesters" ("faculty_id", "start", "end") VALUES (1, '2024-09-23', '2025-01-10');

INSERT INTO "subjects" ("id", "faculty_id", "name", "shortcut", "code_name") VALUES 
(1, 1, 'Právní předpisy ve zdravotnictví a normy', 'PP', 'F7PBKPPN'),
(2, 1, 'Zdravotnické informační systémy', 'ZIS', 'F7ZDIS'),
(3, 1, 'Biomedicínské statistiky', 'BS', 'F7BMBST'),
(4, 1, 'Zdravotnická etika', 'ZE', 'F7ZDET'),
(5, 1, 'Example Subject', 'EX', 'F7EXSUB'); 

INSERT INTO "classes" ("field_id", "start_year", "end_year") VALUES (1, 2022, 2025);

INSERT INTO "subject_instance_types" ("school_id", "name", "description") VALUES 
(1, 'Přednáška', 'Lecture'),
(1, 'Seminář', 'Seminar');

INSERT INTO "subject_instances" ("semester_id", "subject_id", "sit_id", "start_time", "end_time", "day", "room") 
VALUES 
(1, 1, 1, '10:00', '11:50', 1, 'Room 101'),
(1, 2, 1, '08:00', '09:50', 2, 'Room 102'),
(1, 4, 1, '14:00', '15:50', 3, 'Room 103'),
(1, 5, 2, '16:00', '17:50', 4, 'Room 104'),
(1, 3, 1, '10:00', '11:50', 5, 'Room 105');

DO $$
DECLARE
    week_num INT;
BEGIN
    FOR week_num IN 1..14 LOOP
        INSERT INTO "subject_weeks" ("si_id", "week") VALUES (2, week_num);
    END LOOP;
END
$$;