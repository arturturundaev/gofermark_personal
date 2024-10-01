CREATE TABLE IF NOT EXISTS users
(
    id uuid NOT NULL,
    login character varying NOT NULL,
    password character varying NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT user_login UNIQUE (login)
);