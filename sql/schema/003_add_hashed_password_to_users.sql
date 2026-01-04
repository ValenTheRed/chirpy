-- +goose up
alter table if exists users
add column if not exists hashed_password text not null default 'unset';

-- +goose down
alter table if exists users
drop column if exists hashed_password;
