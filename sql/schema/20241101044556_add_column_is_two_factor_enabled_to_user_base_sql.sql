-- +goose Up
-- +goose StatementBegin
alter table pre_go_acc_user_base_9999
add column is_two_factor_enabled int(1) default 0 comment "authentication is enable for the userbase";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table pre_go_acc_user_base_9999
drop column is_two_factor_enabled;
-- +goose StatementEnd
