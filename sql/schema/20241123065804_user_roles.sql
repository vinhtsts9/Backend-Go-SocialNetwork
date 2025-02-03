-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_roles (
    id INT AUTO_INCREMENT PRIMARY KEY,           -- ID duy nhất
    user_id BIGINT UNSIGNED,                     -- ID người dùng
    role_id INT,                                 -- ID vai trò
    room_id bigint,                                 -- ID phòng chat
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Thời gian gán vai trò
    FOREIGN KEY (user_id) REFERENCES user_info(user_id), -- Liên kết với bảng user_info
    FOREIGN KEY (role_id) REFERENCES roles(id), -- Liên kết với bảng roles
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id) -- Liên kết với bảng chat_rooms
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_roles;
-- +goose StatementEnd
