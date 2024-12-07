-- +goose Up
-- +goose StatementBegin
CREATE TABLE role_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,           -- ID duy nhất
    role_id INT,                                 -- ID vai trò
    permission_id INT,                           -- ID quyền
    FOREIGN KEY (role_id) REFERENCES roles(id),  -- Liên kết với bảng roles
    FOREIGN KEY (permission_id) REFERENCES permissions(id) -- Liên kết với bảng permissions
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists role_permissions;
-- +goose StatementEnd
