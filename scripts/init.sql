CREATE EXTENSION IF NOT EXISTS citext;

DROP TRIGGER IF EXISTS inc_threads on threads;
DROP TRIGGER IF EXISTS insert_author on threads;

DROP TABLE IF EXISTS
    users, forums, forum_user, threads
    CASCADE;


CREATE TABLE IF NOT EXISTS users (
    nickname citext UNIQUE NOT NULL,
    fullname varchar(32) NOT NULL,
    email varchar(32) UNIQUE NOT NULL,
    about varchar(256) NOT NULL,
    PRIMARY KEY(nickname, email)
);


CREATE TABLE IF NOT EXISTS forums (
    title varchar(64) NOT NULL,
    author citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    slug citext UNIQUE NOT NULL PRIMARY KEY,
    posts integer NOT NULL DEFAULT 0 CONSTRAINT positive_posts CHECK (posts >= 0),
    threads integer NOT NULL DEFAULT 0 CONSTRAINT positive_threads CHECK (threads >= 0)
);


CREATE TABLE IF NOT EXISTS forum_user (
    nickname citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    UNIQUE(nickname, forum)
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


-- Increment threads number in forums
CREATE OR REPLACE FUNCTION inc_threads() RETURNS trigger AS
$inc_threads$
    BEGIN
        UPDATE forums
        SET threads = threads + 1
        WHERE slug=NEW.forum;
        RETURN NEW;
    END;
$inc_threads$
LANGUAGE plpgsql;

CREATE TRIGGER inc_threads AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE inc_threads();


-- Insert user into forum_user
CREATE OR REPLACE FUNCTION insert_author() RETURNS trigger AS
$insert_autor$
    BEGIN
        INSERT INTO forum_user(nickname, forum)
        VALUES(NEW.author, NEW.forum)
        ON CONFLICT DO NOTHING;
        RETURN NEW;
    END;
$insert_autor$
LANGUAGE plpgsql;

CREATE TRIGGER insert_autor AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE insert_author();
