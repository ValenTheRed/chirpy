-- +goose Up
create table refresh_tokens (
    token text primary key,
    created_at timestamp,
    updated_at timestamp,
    user_id uuid references users(id) on delete cascade,
    expires_at timestamp,
    revoked_at timestamp
);

-- +goose Down
drop table refresh_tokens;
