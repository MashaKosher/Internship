-- +goose Up
-- +goose StatementBegin
-- SELECT 'up SQL query';
ALTER TABLE daily_tasks
ADD COLUMN win_reward DECIMAL(10,2) DEFAULT 0.00,
ADD COLUMN referals_reward DECIMAL(10,2) DEFAULT 0.00;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE daily_tasks
DROP COLUMN win_reward,
DROP COLUMN referals_reward;
-- +goose StatementEnd
