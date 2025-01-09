-- +goose Up
-- +goose StatementBegin
create table Comment (
    id int auto_increment primary key,
    post_id bigint unsigned not null ,                    -- ID duy nhất cho mỗi bài viết
    user_id bigint unsigned NOT NULL,
    comment_content text,
    comment_left int,
    comment_right int,
    comment_parent int,
    isDeleted boolean,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,

    foreign key (post_id) references post(id),
    foreign key (user_id) references user_info(user_id) ,
    foreign key (comment_parent) references Comment(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists Comment;
-- +goose StatementEnd
