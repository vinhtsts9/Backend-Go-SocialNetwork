-- +goose Up
-- +goose StatementBegin
Create table if not exists `pre_go_acc_user_two_factor_9999` (
    `two_factor_id` int unsigned auto_increment primary key,
    `user_id` int unsigned not null,
    `two_factor_auth_type` enum('SMS','EMAIL','APP') not null,
    `two_factor_auth_secret` varchar(255) not null,
    `two_factor_phone` varchar(20) null,
    `two_factor_email` varchar(255) null,
    `two_factor_is_active` boolean not null default true,
    `two_factor_created_at` timestamp default current_timestamp,
    `two_factor_updated_at` timestamp default current_timestamp on update current_timestamp,

    index `idx_user_id` (`user_id`),
    index `idx_auth_type` (`two_factor_auth_type`)
) ENGINE=InnoDb default charset=utf8mb4 comment='pre_go_acc_user_two_factor_9999';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists `pre_go_acc_user_two_factor_9999`;
-- +goose StatementEnd
