-- +goose Up
-- +goose StatementBegin
alter table Comment
add column user_nickname varchar(255) not null references user_info(user_nickname),
add column reply_count int not null default 0,
DROP FOREIGN KEY Comment_ibfk_3,
modify column `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
