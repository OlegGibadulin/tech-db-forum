CREATE EXTENSION IF NOT EXISTS citext;

DROP TRIGGER IF EXISTS inc_threads on threads;
DROP TRIGGER IF EXISTS ins_author on threads;
DROP TRIGGER IF EXISTS inc_posts on posts;
DROP TRIGGER IF EXISTS upd_isEdited on posts;

DROP TABLE IF EXISTS
    users, forums, forum_user, threads, posts
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
    author citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE, -- ins_author
    slug citext UNIQUE NOT NULL PRIMARY KEY,
    posts integer NOT NULL DEFAULT 0 CONSTRAINT positive_posts CHECK (posts >= 0), -- inc_posts
    threads integer NOT NULL DEFAULT 0 CONSTRAINT positive_threads CHECK (threads >= 0) -- inc_threads
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


CREATE TABLE IF NOT EXISTS post (
    id serial PRIMARY KEY,
    parent integer NOT NULL,
    author citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    message varchar NOT NULL,
    isedited boolean DEFAULT FALSE,
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    thread integer NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    created timestamp with time zone NOT NULL DEFAULT now()
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
CREATE OR REPLACE FUNCTION ins_author() RETURNS trigger AS
$ins_author$
    BEGIN
        INSERT INTO forum_user(nickname, forum)
        VALUES(NEW.author, NEW.forum)
        ON CONFLICT DO NOTHING;
        RETURN NEW;
    END;
$ins_author$
LANGUAGE plpgsql;

CREATE TRIGGER ins_author AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE ins_author();


-- Increment threads number in forums
CREATE OR REPLACE FUNCTION inc_posts() RETURNS trigger AS
$inc_posts$
    BEGIN
        UPDATE forums
        SET posts = posts + 1
        WHERE slug=NEW.forum;
        RETURN NEW;
    END;
$inc_posts$
LANGUAGE plpgsql;

CREATE TRIGGER inc_posts AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE inc_posts();


-- Set isEdited true if message was updated
CREATE OR REPLACE FUNCTION upd_isEdited() RETURNS trigger AS
$upd_isEdited$
    BEGIN
        IF NEW.message <> OLD.message THEN
            NEW.isEdited = TRUE;
        END IF;
        RETURN NEW;
    END;
$upd_isEdited$
LANGUAGE plpgsql;

CREATE TRIGGER upd_isEdited AFTER UPDATE ON posts
    FOR EACH ROW EXECUTE PROCEDURE upd_isEdited();
