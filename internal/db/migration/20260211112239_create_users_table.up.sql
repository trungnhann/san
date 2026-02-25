CREATE TABLE IF NOT EXISTS users (
   id          text        NOT NULL,
   username    text        NOT NULL UNIQUE,
   email       text        NOT NULL UNIQUE,
   password    text        NOT NULL,
   bio         text,
   image       text,
   created_at  timestamptz NOT NULL DEFAULT now(),
   updated_at  timestamptz NOT NULL DEFAULT now(),
   PRIMARY KEY (id)
);