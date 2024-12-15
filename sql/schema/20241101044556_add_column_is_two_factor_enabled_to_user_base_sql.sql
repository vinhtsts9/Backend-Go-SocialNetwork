-- +goose Up
-- +goose StatementBegin
alter table user_base
add column is_two_factor_enabled int(1) default 0 comment "authentication is enable for the userbase";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table user_base
drop column is_two_factor_enabled;
-- +goose StatementEnd
