-- +goose Up
-- +goose StatementBegin
ALTER TABLE post
DROP COLUMN content,
ADD COLUMN image_paths JSON ,
ADD COLUMN user_nickname  varchar(255) not null references user_info(user_nickname);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
