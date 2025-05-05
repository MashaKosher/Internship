-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS task_status_tbl(
    task_date DATE NOT NULL,
    user_id INTEGER NOT NULL,
    win INTEGER DEFAULT 0,
    win_status BOOLEAN NOT NULL DEFAULT FALSE,
    referals INTEGER DEFAULT 0,
    referals_status BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (task_date, user_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (task_date) REFERENCES daily_tasks (task_date) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task_status_tbl;
-- +goose StatementEnd
