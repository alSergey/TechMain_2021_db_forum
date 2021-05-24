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
