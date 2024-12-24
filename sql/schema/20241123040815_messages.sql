-- +goose Up
-- +goose StatementBegin
create table messages (
    id int auto_increment primary key,
    room_id int,
    sender_id bigint unsigned,
    message_context text,
    message_type enum('text','image','video','file') not null default 'text',
    is_pinned boolean default false,
    is_announcement boolean default false,
    user_nickname varchar(255) not null references user_info(user_nickname),
    created_at timestamp default current_timestamp,
    foreign key (room_id) references chat_rooms(id),
    foreign key (sender_id) references user_info(user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists messages;
-- +goose StatementEnd
