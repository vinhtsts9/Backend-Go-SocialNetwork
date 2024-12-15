-- +goose Up
-- +goose StatementBegin
ALTER TABLE pre_go_acc_user_verify_9999 RENAME TO user_verify;
-- +goose StatementEnd



-- +goose StatementBegin
ALTER TABLE pre_go_acc_user_two_factor_9999 RENAME TO user_two_factor;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_verify RENAME TO pre_go_acc_user_verify_9999;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_base RENAME TO pre_go_acc_user_base_9999;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_info RENAME TO pre_go_acc_user_info_9999;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE user_two_factor RENAME TO pre_go_acc_user_two_factor_9999;
-- +goose StatementEnd
