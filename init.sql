CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users
(
    nickname CITEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    about    TEXT,
    email    CITEXT UNIQUE
);

CREATE TABLE forum
(
    title   TEXT,
    "user"  CITEXT,
    slug    CITEXT PRIMARY KEY,
    posts   BIGINT DEFAULT 0,
    threads BIGINT DEFAULT 0,

    FOREIGN KEY ("user") REFERENCES "users" (nickname)
);

CREATE TABLE thread
(
    id      SERIAL PRIMARY KEY,
    title   TEXT NOT NULL,
    author  CITEXT,
    forum   CITEXT,
    message TEXT NOT NULL,
    votes   INT                      DEFAULT 0,
    slug    CITEXT UNIQUE,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    FOREIGN KEY (author) REFERENCES "users" (nickname),
    FOREIGN KEY (forum) REFERENCES "forum" (slug)
);

CREATE TABLE post
(
    id       BIGSERIAL PRIMARY KEY,
    parent   BIGINT                   DEFAULT 0,
    author   CITEXT NOT NULL,
    message  TEXT   NOT NULL,
    isEdited BOOLEAN                  DEFAULT FALSE,
    forum    CITEXT,
    thread   INT,
    created  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    path     BIGINT[]                 DEFAULT ARRAY []::INTEGER[],

    FOREIGN KEY (author) REFERENCES "users" (nickname),
    FOREIGN KEY (forum) REFERENCES "forum" (slug),
    FOREIGN KEY (thread) REFERENCES "thread" (id)
--     FOREIGN KEY (thread) REFERENCES "post" (id)
);

CREATE TABLE votes
(
    id       BIGSERIAL PRIMARY KEY,
    nickname CITEXT NOT NULL,
    voice    INT    NOT NULL,
    thread   INT    NOT NULL,

    FOREIGN KEY (nickname) REFERENCES "users" (nickname),
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    UNIQUE (nickname, thread)
);


CREATE OR REPLACE FUNCTION insertVote() RETURNS TRIGGER AS
$update_thread$
BEGIN
    UPDATE thread SET votes=(votes + NEW.voice) WHERE id = NEW.thread;
    return NEW;
END
$update_thread$ LANGUAGE plpgsql;

CREATE TRIGGER insert_voice
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insertVote();

CREATE OR REPLACE FUNCTION updateVote() RETURNS TRIGGER AS
$update_thread$
BEGIN
    IF OLD.voice <> NEW.voice THEN
        UPDATE thread SET votes=(votes + NEW.Voice * 2) WHERE id = NEW.thread;
    END IF;
    return NEW;
END
$update_thread$ LANGUAGE plpgsql;

CREATE TRIGGER update_voice
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE updateVote();


CREATE OR REPLACE FUNCTION updateThread() RETURNS TRIGGER AS
$update_forum$
BEGIN
    UPDATE forum SET Threads=(Threads + 1) WHERE slug = NEW.forum;
    return NEW;
END
$update_forum$ LANGUAGE plpgsql;

CREATE TRIGGER insert_thread
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE updateThread();


CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_path$
DECLARE
    parent_path         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS NULL) THEN
        NEW.path := array_append(NEW.path, NEW.id);
    ELSE
        SELECT path FROM post WHERE id = NEW.parent INTO parent_path;
        SELECT thread FROM post WHERE id = parent_path[1] INTO first_parent_thread;

        IF NOT FOUND OR first_parent_thread <> NEW.thread THEN
            RAISE EXCEPTION 'parent post was created in another thread' USING ERRCODE = '12345';
        END IF;

        NEW.path := NEW.path || parent_path || NEW.id;
    END IF;

    UPDATE forum SET posts=posts + 1 WHERE forum.slug = NEW.forum;
    RETURN NEW;
END
$update_path$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();