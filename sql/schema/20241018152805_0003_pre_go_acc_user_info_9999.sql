-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_acc_user_info_9999 (
    user_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT'User ID',
    user_account VARCHAR(255) NOT NULL,
    user_nickname VARCHAR(255),
    user_avatar VARCHAR(255),
    user_state TINYINT UNSIGNED NOT NULL,
    user_mobile VARCHAR(255),
    user_gender TINYINT UNSIGNED,
    user_birthday DATE,
    user_email VARCHAR(255),
    user_is_authencation TINYINT UNSIGNED NOT NULL COMMENT 'Authentication status: 0-not, 1-pending,2-authen',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'record creation time',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'record update time',

    UNIQUE KEY unique_user_account (user_account),
    INDEX idx_user_mobile (user_mobile),
    INDEX idx_user_email (user_email),
    INDEX idx_user_state (user_state),
    INDEX idx_user_is_authencation (user_is_authencation)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='pre_go_acc_user_9999';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_9999`;
-- +goose StatementEnd
