CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users
(
    nickname CITEXT PRIMARY KEY,
    fullname TEXT NOT NULL,
    about    TEXT,
    email    CITEXT UNIQUE
);