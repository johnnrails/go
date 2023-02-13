-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :exec
INSERT INTO users (id, name, email, bio, password) VALUES (?, ?, ?, ?, ?);

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: UpdateUser :exec
UPDATE users set name = ?, email = ?, bio = ?, password = ? WHERE id = ?;

-- name: GetProduct :one
SELECT * FROM products WHERE id = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY name;

-- name: CreateProduct :exec
INSERT INTO products (id, name, code, price, user_id) VALUES (?, ?, ?, ?, ?);

-- name: UpdateProduct :exec
UPDATE products set name = ?, code = ?, price = ? WHERE id = ?;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;