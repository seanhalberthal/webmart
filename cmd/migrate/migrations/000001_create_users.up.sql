CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY                     DEFAULT gen_random_uuid(),
    name       VARCHAR(255)                NOT NULL,
    username   VARCHAR(255) UNIQUE         NOT NULL,
    email      citext UNIQUE               NOT NULL,
    password   bytea                       NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
