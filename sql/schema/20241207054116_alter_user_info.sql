-- +goose Up
-- +goose StatementBegin
alter table user_info
modify column user_id BIGINT UNSIGNED AUTO_INCREMENT not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
