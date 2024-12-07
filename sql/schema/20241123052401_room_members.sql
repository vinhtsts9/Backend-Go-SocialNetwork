-- +goose Up
-- +goose StatementBegin
CREATE TABLE room_members (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- ID duy nhất cho mỗi bản ghi
    room_id BIGINT NOT NULL,  -- ID phòng chat
    user_id BIGINT UNSIGNED NOT NULL,  -- ID người dùng
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Thời gian tham gia
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id),  -- Liên kết với bảng chat_rooms
    FOREIGN KEY (user_id) REFERENCES user_info(user_id)  -- Liên kết với bảng users
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists room_members;
-- +goose StatementEnd
