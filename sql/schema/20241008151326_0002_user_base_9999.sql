-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_base (
	user_id INT AUTO_INCREMENT PRIMARY KEY,
    user_account VARCHAR(255) NOT NULL,
    user_password VARCHAR(255) NOT NULL,
    user_salt VARCHAR(255) NOT NULL,
    is_two_factor_enabled int(1) default 0 comment "authentication is enable for the userbase",
    user_login_time TIMESTAMP NULL DEFAULT NULL,
    user_logout_time TIMESTAMP NULL DEFAULT NULL,
    user_login_ip VARCHAR(255),
    user_state tinyint unsigned not null default 3 comment'1-online, 2-active,3-unactive',
    
    user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE index unique_user_account (user_account)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='user_base';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_base;
-- +goose StatementEnd
