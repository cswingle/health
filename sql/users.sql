CREATE EXTENSION pgcrypto;

CREATE TABLE users (
  id serial PRIMARY KEY,
  username text,
  password text,
  access_level text,
  UNIQUE(username)
);
