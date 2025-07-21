-- name: CreateUser :one
insert into users (username, hashed_password, full_name, email)
values ($1, $2, $3, $4) RETURNING *;

-- name: GetUser :one
SELECT * From users
WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    hashed_password = coalesce(sqlc.narg(hashed_password),hashed_password),
    full_name = coalesce(sqlc.narg(full_name),full_name),
    email = coalesce(sqlc.narg(email), email)
WHERE
    username = sqlc.arg(username)
RETURNING *;