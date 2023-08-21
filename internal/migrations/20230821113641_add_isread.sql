-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS social.messages
    ADD COLUMN is_read bool;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS social.messages
    DROP COLUMN is_read CASCADE ;
-- +goose StatementEnd
