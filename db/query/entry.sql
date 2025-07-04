
-- name: CreateEntry :one
insert into entries ( account_id, amount)
values ($1, $2) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries WHERE id = $1;

-- name: ListEntries :many
SELECT * FROM entries WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3;