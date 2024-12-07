-- +goose Up
-- +goose StatementBegin
alter table permissions
add column resource varchar(255) not null,
add column action varchar(50) not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
