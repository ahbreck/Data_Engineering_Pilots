DROP DATABASE IF EXISTS go;
CREATE DATABASE go;

DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Userdata;

\c go;

CREATE TABLE Users (
    ID SERIAL,
    Username VARCHAR(100) PRIMARY KEY
);

CREATE TABLE Userdata (
    UserID Int NOT NULL,
    Name VARCHAR(100),
    Surname VARCHAR(100),
    Description VARCHAR(200)
);

CREATE TABLE msds_course_catalog (
    course_id   VARCHAR(32) PRIMARY KEY,
    course_name TEXT NOT NULL,
    prerequisite TEXT
);