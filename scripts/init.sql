CREATE EXTENSION IF NOT EXISTS citext;
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

CREATE TABLE IF NOT EXISTS forums (
    title varchar(64) NOT NULL,
    author citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    slug citext UNIQUE NOT NULL PRIMARY KEY,
    posts integer NOT NULL DEFAULT 0 CONSTRAINT positive_posts CHECK (posts >= 0),
    threads integer NOT NULL DEFAULT 0 CONSTRAINT positive_threads CHECK (threads >= 0)
);
