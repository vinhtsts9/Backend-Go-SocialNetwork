-- +goose Up
-- +goose StatementBegin
create table user_follows (
    follow_id bigint auto_increment primary key,
    follower_id bigint unsigned ,
    following_id bigint unsigned,
    created_at timestamp default current_timestamp,
    foreign key (follower_id) references user_info(user_id),
    foreign key (following_id) references user_info(user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_follows;
-- +goose StatementEnd
