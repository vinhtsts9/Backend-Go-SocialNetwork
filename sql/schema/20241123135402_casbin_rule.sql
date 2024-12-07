-- +goose Up
-- +goose StatementBegin
CREATE TABLE casbin_rule (
    ptype VARCHAR(100),
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists casbin_rule;
-- +goose StatementEnd
