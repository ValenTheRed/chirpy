-- name: CreateRefreshToken :one
insert into refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
values ($1, now(), now(), $2, now() + interval '60 days', null)

returning *;

-- name: GetRefreshToken :one
select *
from refresh_tokens
where token = $1 and revoked_at is null and now() <= expires_at;

-- name: RevokeRefreshToken :exec
update refresh_tokens
set revoked_at = now(), updated_at = now()
where token = $1

returning *;
