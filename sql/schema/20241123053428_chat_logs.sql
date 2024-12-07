-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- ID duy nhất
    room_id INT,  -- ID phòng chat
    user_id BIGINT UNSIGNED,  -- ID người dùng
    event_type ENUM('joined', 'left', 'role_changed', 'name_changed', 'avatar_changed') NOT NULL,  -- Loại sự kiện
    event_details TEXT,  -- Chi tiết sự kiện
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Thời gian sự kiện
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id),  -- Liên kết với bảng chat_rooms
    FOREIGN KEY (user_id) REFERENCES user_info(user_id)  -- Liên kết với bảng user_info
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists chat_logs;
-- +goose StatementEnd
