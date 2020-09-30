-- name: CreateEntry :one
INSERT INTO entrys (
    account_id,
    amount
) VALUES (
  $1, $2, $3
) RETURNING *;


-- name: GetEntry :one
SELECT * FROM entrys
WHERE id = $1 LIMIT 1;

-- name: ListEntrys :many
SELECT * FROM entrys
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :one
UPDATE entrys SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entrys WHERE id = $1;