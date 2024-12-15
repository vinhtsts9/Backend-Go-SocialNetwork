-- +goose Up
-- +goose StatementBegin
alter table user_base
add column user_state tinyint unsigned not null default 3 comment'1-online, 2-active,3-unactive';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
