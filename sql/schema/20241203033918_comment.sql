-- +goose Up
-- +goose StatementBegin
create table Comment (
    id int auto_increment primary key,
    post_id bigint unsigned not null ,                    -- ID duy nhất cho mỗi bài viết
    user_id bigint unsigned NOT NULL,
    user_nickname varchar(255) not null references user_info(user_nickname),
    reply_count int not null default 0,
    
    comment_content text not null,
    comment_left int not null,
    comment_right int not null,
    comment_parent int default null ,
    isDeleted boolean default false,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
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
