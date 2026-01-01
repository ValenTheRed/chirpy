-- +goose Up
create table users (
    id uuid primary key,
    created_at timestamp,
    updated_at timestamp,
    email text
);

-- +goose Down
drop table users;
