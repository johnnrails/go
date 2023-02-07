
-- name: GetAuthor :one
SELECT * FROM authors WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors ORDER BY name;

-- name: CreateAuthor :exec
INSERT INTO authors (name, bio) VALUES (?, ?);

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = ?;

-- name: UpdateAuthor :exec
UPDATE authors set name = ?, bio = ? WHERE id = ?;

-- name: GetProduct :one
SELECT * FROM products WHERE id = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY name;

-- name: CreateProduct :exec
INSERT INTO products (price, name) VALUES (?, ?);

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;

-- name: UpdateProduct :exec
UPDATE products set name = ?, price = ? WHERE id = ?;
