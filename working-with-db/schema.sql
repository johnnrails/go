CREATE TABLE users (
  id text PRIMARY KEY,
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  bio VARCHAR,
  password VARCHAR NOT NULL
);

CREATE TABLE products (
  id text PRIMARY KEY,
  name VARCHAR NOT NULL,
  code VARCHAR NOT NULL,
  price integer NOT NULL,
  user_id text NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id)
);