-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (gen_random_uuid(), now(), now(), $1, $2)

returning *;

-- name: DeleteAllUsers :exec
delete from users

returning *;

-- name: GetUsersHashedPassword :one
select hashed_password
from users
where email = $1;
