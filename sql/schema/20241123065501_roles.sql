-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id INT AUTO_INCREMENT PRIMARY KEY,           -- ID vai trò
    role_name VARCHAR(50) NOT NULL UNIQUE,        -- Tên vai trò (Admin, Member, Moderator)
    description TEXT                             -- Mô tả vai trò
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists roles;
-- +goose StatementEnd
