-- +goose Up
-- +goose StatementBegin
CREATE TABLE permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,           -- ID quyền
    permission_name VARCHAR(100) NOT NULL UNIQUE, -- Tên quyền (Send Message, Delete Message, etc.)
    description TEXT                             -- Mô tả quyền
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists permissions;
-- +goose StatementEnd
