-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages_status (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- ID duy nhất
    message_id INT,  -- ID tin nhắn
    user_id BIGINT UNSIGNED,  -- ID người dùng
    status ENUM('sent', 'received', 'read','sending','deleted','edited') NOT NULL,  -- Trạng thái tin nhắn
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- Thời gian cập nhật
    FOREIGN KEY (message_id) REFERENCES messages(id),  -- Liên kết với bảng messages
    FOREIGN KEY (user_id) REFERENCES user_info(user_id)  -- Liên kết với bảng user_info
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists messages_status;
-- +goose StatementEnd
