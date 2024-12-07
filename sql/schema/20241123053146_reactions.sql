-- +goose Up
-- +goose StatementBegin
CREATE TABLE reactions (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- ID phản ứng
    user_id BIGINT UNSIGNED,  -- ID người dùng
    message_id INT,  -- ID tin nhắn
    reaction_type ENUM('like', 'love', 'laugh', 'angry', 'sad') NOT NULL,  -- Loại phản ứng
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Thời gian phản ứng
    FOREIGN KEY (user_id) REFERENCES user_info(user_id),  -- Liên kết với bảng user_info
    FOREIGN KEY (message_id) REFERENCES messages(id)  -- Liên kết với bảng messages
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists reactions;
-- +goose StatementEnd
