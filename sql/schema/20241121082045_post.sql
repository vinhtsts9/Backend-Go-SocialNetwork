-- +goose Up
-- +goose StatementBegin
CREATE TABLE post (
    id SERIAL PRIMARY KEY,                    -- ID duy nhất cho mỗi bài viết
    user_id bigint unsigned NOT NULL,                     -- ID của người đăng bài
    title VARCHAR(255) NOT NULL,              -- Tiêu đề bài viết
    image_paths JSON,
    user_nickname varchar(255) not null references user_info(user_nickname),
    created_at TIMESTAMP DEFAULT NOW(),       -- Thời gian tạo bài viết
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(), -- Thời gian cập nhật bài viết
    privacy ENUM('private', 'friends', 'public') NOT NULL DEFAULT 'public',       -- Trạng thái bài viết (đã đăng hay chưa)
    metadata JSON DEFAULT NULL,               -- Metadata bổ sung cho bài viết (tags, views, v.v.)
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES user_info(user_id) -- Ràng buộc liên kết với bảng user_info
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists post;
-- +goose StatementEnd
