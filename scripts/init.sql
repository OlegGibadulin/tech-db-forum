CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS
    users, forums, forum_user, threads, posts, votes
    CASCADE;


CREATE TABLE IF NOT EXISTS users (
    nickname citext UNIQUE NOT NULL,
    fullname varchar(32) NOT NULL,
    email citext UNIQUE NOT NULL,
    about varchar NOT NULL,
    PRIMARY KEY(nickname, email)
);
CREATE INDEX IF NOT EXISTS users_nickname ON users (nickname);
CREATE INDEX IF NOT EXISTS users_cover ON users (nickname, email, fullname, about);


CREATE TABLE IF NOT EXISTS forums (
    title varchar NOT NULL,
    author citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE, -- ins_author
    slug citext UNIQUE NOT NULL PRIMARY KEY,
    posts integer NOT NULL DEFAULT 0 CONSTRAINT positive_posts CHECK (posts >= 0), -- inc_posts
    threads integer NOT NULL DEFAULT 0 CONSTRAINT positive_threads CHECK (threads >= 0) -- inc_threads
);
CREATE INDEX IF NOT EXISTS forums_author ON forums (author);


CREATE TABLE IF NOT EXISTS forum_user (
    nickname citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    UNIQUE(nickname, forum)
);
CREATE INDEX IF NOT EXISTS forum_user_nickname ON forum_user (nickname);
CREATE INDEX IF NOT EXISTS forum_user_forum ON forum_user (forum);


CREATE TABLE IF NOT EXISTS threads (
    id serial PRIMARY KEY,
    title varchar NOT NULL,
    author citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    message varchar NOT NULL,
    created timestamp with time zone NOT NULL DEFAULT now(),
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    votes integer NOT NULL DEFAULT 0,
    slug citext NOT NULL
);
CREATE INDEX IF NOT EXISTS threads_forum ON threads (forum);
CREATE INDEX IF NOT EXISTS threads_slug ON threads (slug, id);
CREATE INDEX IF NOT EXISTS threads_id_forum ON threads (id, forum);
CREATE INDEX IF NOT EXISTS threads_forum_created ON threads (forum, created);


CREATE TABLE IF NOT EXISTS posts (
    id serial PRIMARY KEY,
    parent integer NOT NULL,
    author citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    message varchar NOT NULL,
    isedited boolean DEFAULT FALSE,
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    thread integer NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    created timestamp with time zone NOT NULL DEFAULT now(),
    path INTEGER[] NOT NULL
);
CREATE INDEX IF NOT EXISTS posts_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS posts_thread_id_asc ON posts (thread, id ASC);
CREATE INDEX IF NOT EXISTS posts_thread_id_desc ON posts (thread, id DESC);
CREATE INDEX IF NOT EXISTS posts_thread_path_asc ON posts (thread, path ASC);
CREATE INDEX IF NOT EXISTS posts_thread_path_desc ON posts (thread, path DESC);


CREATE TABLE IF NOT EXISTS votes (
    nickname citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    thread integer NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    voice integer NOT NULL,
    UNIQUE(nickname, thread)
);
CREATE INDEX IF NOT EXISTS votes_nickname ON votes (nickname);
CREATE INDEX IF NOT EXISTS votes_thread ON votes (thread);


DROP TRIGGER IF EXISTS inc_threads ON threads;
DROP TRIGGER IF EXISTS ins_author_on_ins_thread ON threads;
DROP TRIGGER IF EXISTS ins_author_on_ins_post ON threads;
DROP TRIGGER IF EXISTS inc_posts ON posts;
DROP TRIGGER IF EXISTS upd_isEdited ON posts;
DROP TRIGGER IF EXISTS upd_votes_on_insert ON votes;
DROP TRIGGER IF EXISTS upd_votes_on_update ON votes;
DROP TRIGGER IF EXISTS upd_path ON posts;


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

CREATE TRIGGER ins_author_on_ins_thread AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE ins_author();

CREATE TRIGGER ins_author_on_ins_post AFTER INSERT ON posts
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

CREATE TRIGGER upd_isEdited BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE PROCEDURE upd_isEdited();


-- Update sum of thread votes on insert
CREATE OR REPLACE FUNCTION upd_votes_on_insert() RETURNS trigger AS
$upd_votes_on_insert$
    BEGIN
        UPDATE threads
        SET votes = votes + NEW.voice
        WHERE id=NEW.thread;
        RETURN NEW;
    END;
$upd_votes_on_insert$
LANGUAGE plpgsql;

CREATE TRIGGER upd_votes_on_insert AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE upd_votes_on_insert();


-- Update sum of thread votes on update
CREATE OR REPLACE FUNCTION upd_votes_on_update() RETURNS trigger AS
$upd_votes_on_update$
    BEGIN
        IF NEW.voice <> OLD.voice THEN
            UPDATE threads
            SET votes = votes - OLD.voice + NEW.voice
            WHERE id=NEW.thread;
        END IF;
        RETURN NEW;
    END;
$upd_votes_on_update$
LANGUAGE plpgsql;

CREATE TRIGGER upd_votes_on_update AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE upd_votes_on_update();


-- Insert id into path on insert
CREATE OR REPLACE FUNCTION upd_path() RETURNS trigger AS
$upd_path$
    DECLARE
        parent_thread integer;
        parent_path integer[];
    BEGIN
        IF (NEW.parent = 0) THEN
            NEW.path := array_append(NEW.path, NEW.id);
            RETURN NEW;
        END IF;
        
        SELECT thread INTO parent_thread FROM posts WHERE id=NEW.parent;
        IF NOT FOUND OR NEW.thread <> parent_thread THEN
            RAISE EXCEPTION 'Can not find parent post into thread';
        END IF;

        SELECT path INTO parent_path FROM posts WHERE id=NEW.parent;
        NEW.path = array_append(parent_path, NEW.id);
        RETURN NEW;
    END;
$upd_path$
LANGUAGE plpgsql;

CREATE TRIGGER upd_path BEFORE INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE upd_path();
