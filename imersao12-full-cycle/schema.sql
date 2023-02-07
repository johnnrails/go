CREATE TABLE authors (
  id INTEGER PRIMARY KEY,
  name text NOT NULL,
  bio text
);

CREATE TABLE products (
  id INTEGER PRIMARY KEY,
  price FLOAT NOT NULL,
  name text NOT NULL
);
