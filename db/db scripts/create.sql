CREATE TABLE "Semester" (
    "id" serial   NOT NULL,
    "codename" char(4)   NOT NULL,
    "start" date NOT NULL,
    "end" date NOT NULL,
    CONSTRAINT "pk_Semester" PRIMARY KEY (
        "id"
     )
);

CREATE TABLE "Field" (
    "id" serial   NOT NULL,
    "name" varchar(64)   NOT NULL,
    "shortcut" char(4) not null,
    CONSTRAINT "pk_Field" PRIMARY KEY (
        "id"
     )
);

CREATE TABLE "Class" (
    "id" serial   NOT NULL,
    "start_year" int   NOT NULL,
    "end_year" int   NOT NULL,
    "field_id" int   NOT NULL,
    CONSTRAINT "pk_Class" PRIMARY KEY (
        "id"
     )
);

CREATE TABLE "Subject" (
    "id" serial   NOT NULL,
    "name" varchar(64)   NOT NULL,
    "shortcut" char(3)   NOT NULL,
    "code_name" varchar(32),
    CONSTRAINT "pk_Subject" PRIMARY KEY (
        "id"
     )
);


create table "SubjectClassType" (
    "id" smallserial not null,
    "name" varchar(32) not null,
    constraint "pk_SubjectClassType" primary key (
                                            "id"
        )
);

CREATE TABLE "SubjectClass" (
    "id" serial   NOT NULL,
    "semester_id" int   NOT NULL,
    "subject_id" int   NOT NULL,
    "class" int   NOT NULL,
    "start_time" char(5)   NOT NULL,
    "end_time" char(5)   NOT NULL,
    "day" smallint NOT NULL,
    "type" smallint   NOT NULL,
    CONSTRAINT "pk_SubjectClass" PRIMARY KEY (
        "id"
     )
);

CREATE TABLE "SubjectWeek" (
    "id" serial   NOT NULL,
    "week_number" smallint   NOT NULL,
    "subject_class_id" int   NOT NULL,
    CONSTRAINT "pk_SubjectWeek" PRIMARY KEY (
        "id"
     )
);

ALTER TABLE "Class" ADD CONSTRAINT "fk_Class_field_id" FOREIGN KEY("field_id")
REFERENCES "Field" ("id");

ALTER TABLE "SubjectClass" ADD CONSTRAINT "fk_SubjectClass_semester_id" FOREIGN KEY("semester_id")
REFERENCES "Semester" ("id");

ALTER TABLE "SubjectClass" ADD CONSTRAINT "fk_SubjectClass_subject_id" FOREIGN KEY("subject_id")
REFERENCES "Subject" ("id");

ALTER TABLE "SubjectClass" ADD CONSTRAINT "fk_SubjectClass_class" FOREIGN KEY("class")
REFERENCES "Class" ("id");

ALTER TABLE "SubjectWeek" ADD CONSTRAINT "fk_SubjectWeek_subject_class_id" FOREIGN KEY("subject_class_id")
REFERENCES "SubjectClass" ("id");

alter table "SubjectClass" add constraint  "fk_SubjectClassType_SubjectClass" foreign key ("type")
references "SubjectClassType" ("id");

