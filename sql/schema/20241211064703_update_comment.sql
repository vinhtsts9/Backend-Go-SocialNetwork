-- +goose Up
-- +goose StatementBegin

alter table Comment
modify column comment_content text not null,
modify column comment_left int not null,
modify column comment_right int not NULL,
modify column comment_parent int DEFAULT NULL,
modify column isDeleted bool DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
