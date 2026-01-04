-- name: Login :one
select *
from users
where email = $1 and hashed_password = $2;
