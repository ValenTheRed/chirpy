-- +goose Up
create table chirps (
    id uuid primary key,
    user_id uuid references users(id) on delete cascade,
    created_at timestamp,
    updated_at timestamp,
    body text
);

-- +goose Down
drop table chirps;
