DROP TABLE IF EXISTS
    users, forums
    CASCADE;

CREATE TABLE IF NOT EXISTS users (
    nickname citext UNIQUE NOT NULL,
    fullname varchar(32) NOT NULL,
    email varchar(32) UNIQUE NOT NULL,
    about varchar(128) NOT NULL,
    PRIMARY KEY(nickname, email)
);
