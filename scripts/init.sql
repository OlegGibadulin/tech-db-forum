CREATE EXTENSION IF NOT EXISTS citext;
DROP TABLE IF EXISTS
    users, forums, threads
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

CREATE TABLE IF NOT EXISTS threads (
    id serial PRIMARY KEY,
    title varchar(64) NOT NULL,
    author citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    message varchar NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    votes integer NOT NULL DEFAULT 0 CONSTRAINT positive_threads CHECK (votes >= 0),
    slug citext NOT NULL
);
