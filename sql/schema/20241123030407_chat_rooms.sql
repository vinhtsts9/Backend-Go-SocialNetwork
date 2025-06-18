-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_rooms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_group BOOLEAN NOT NULL,
    admin_id BIGINT unsigned NOT NULL,
    avatar_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    constraint fk_admin_id foreign key (admin_id)
    references user_info(user_id)
    on delete cascade
    on update cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists chat_rooms;
-- +goose StatementEnd
