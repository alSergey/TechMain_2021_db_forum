CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users
(
    nickname CITEXT PRIMARY KEY,
    fullname TEXT,
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
    title   TEXT,
    author  CITEXT,
    forum   CITEXT,
    message TEXT,
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
    author   CITEXT,
    message  TEXT,
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
    nickname CITEXT,
    voice    INT,
    thread   INT NOT NULL,

    FOREIGN KEY (nickname) REFERENCES "users" (nickname),
    FOREIGN KEY (thread) REFERENCES "thread" (id),

    UNIQUE (nickname, thread)
);

CREATE TABLE forum_users
(
    nickname CITEXT,
    fullname TEXT,
    about    TEXT,
    email    CITEXT,
    forum    CITEXT,

    FOREIGN KEY (nickname) REFERENCES "users" (nickname),
    FOREIGN KEY (forum) REFERENCES "forum" (slug),

    UNIQUE (nickname, forum)
);



CREATE OR REPLACE FUNCTION afterInsertVote() RETURNS TRIGGER AS
$after_insert_voice$
BEGIN
    UPDATE thread SET votes=(votes + NEW.voice) WHERE id = NEW.thread;
    RETURN NEW;
END
$after_insert_voice$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert_voice
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE afterInsertVote();


CREATE OR REPLACE FUNCTION afterUpdateVote() RETURNS TRIGGER AS
$after_update_voice$
BEGIN
    IF OLD.voice <> NEW.voice THEN
        UPDATE thread SET votes=(votes + NEW.Voice * 2) WHERE id = NEW.thread;
    END IF;

    RETURN NEW;
END
$after_update_voice$ LANGUAGE plpgsql;

CREATE TRIGGER after_update_voice
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE afterUpdateVote();



CREATE OR REPLACE FUNCTION afterInsertThread() RETURNS TRIGGER AS
$after_insert_thread$
DECLARE
    author_nickname CITEXT;
    author_fullname TEXT;
    author_about    TEXT;
    author_email    CITEXT;
BEGIN
    UPDATE forum SET threads=(threads + 1) WHERE slug = NEW.forum;

    SELECT nickname, fullname, about, email
    FROM users
    WHERE nickname = NEW.author
    INTO author_nickname, author_fullname, author_about, author_email;

    INSERT INTO forum_users (nickname, fullname, about, email, forum)
    VALUES (author_nickname, author_fullname, author_about, author_email, NEW.forum)
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END
$after_insert_thread$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert_thread
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE afterInsertThread();



CREATE OR REPLACE FUNCTION beforeInsertPost() RETURNS TRIGGER AS
$before_insert_post$
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
$before_insert_post$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_post
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE beforeInsertPost();


CREATE OR REPLACE FUNCTION afterInsertPost() RETURNS TRIGGER AS
$after_insert_post$
DECLARE
    author_nickname CITEXT;
    author_fullname TEXT;
    author_about    TEXT;
    author_email    CITEXT;
BEGIN
    SELECT nickname, fullname, about, email
    FROM users
    WHERE nickname = NEW.author
    INTO author_nickname, author_fullname, author_about, author_email;

    INSERT INTO forum_users (nickname, fullname, about, email, forum)
    VALUES (author_nickname, author_fullname, author_about, author_email, NEW.forum)
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END
$after_insert_post$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert_post
    AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE afterInsertPost();

CREATE INDEX all_users_forum ON forum_users (nickname, fullname, about, email);
CLUSTER forum_users USING all_users_forum;
CREATE INDEX nickname_users_forum ON forum_users using hash (nickname);
CREATE INDEX f_a_e_users_forum ON forum_users (fullname, about, email);

CREATE INDEX IF NOT EXISTS user_nickname ON users using hash (nickname);
CREATE INDEX IF NOT EXISTS user_email ON users using hash (email);
CREATE INDEX IF NOT EXISTS forum_slug ON forum using hash (slug);
CREATE UNIQUE INDEX IF NOT EXISTS  forum_users_unique on forum_users (forum, nickname);
-- CLUSTER users_forum USING forum_users_unique;

CREATE INDEX IF NOT EXISTS  thr_slug ON thread using hash (slug);
CREATE INDEX IF NOT EXISTS  thr_date ON thread (created);
CREATE INDEX IF NOT EXISTS  thr_forum ON thread using hash (forum);
CREATE INDEX IF NOT EXISTS  thr_forum_date ON thread (forum, created);

CREATE INDEX IF NOT EXISTS post_id_path on post (id, (path[1]));
CREATE INDEX IF NOT EXISTS post_thread_id_path1_parent on post (thread, id, (path[1]), parent);
CREATE INDEX IF NOT EXISTS post_thread_path_id on post (thread, path, id);

CREATE INDEX IF NOT EXISTS post_path1 on post ((path[1]));
CREATE INDEX IF NOT EXISTS post_thread_id on post (thread, id);
CREATE INDEX IF NOT EXISTS post_thr_id ON post (thread);

CREATE UNIQUE INDEX IF NOT EXISTS  vote_unique on votes (nickname, thread);

CREATE INDEX IF NOT EXISTS post_path1_path_id_desc ON post ((path[1]) DESC, path, id);
CREATE INDEX IF NOT EXISTS post_path1_path_id_asc ON post ((path[1]) DESC, path, id);