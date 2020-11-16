DROP TABLE IF EXISTS
    users
    CASCADE;

CREATE TABLE IF NOT EXISTS users (
    nickname varchar(32) UNIQUE NOT NULL,
    fullname varchar(32) NOT NULL,
    email varchar(32) UNIQUE NOT NULL,
    about varchar(128) NOT NULL,
    PRIMARY KEY(nickname, email)
);
