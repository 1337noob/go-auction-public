CREATE TABLE scheduler.tasks
(
    id           VARCHAR(255) PRIMARY KEY,
    aggregate_id VARCHAR(255) NOT NULL,
    command      VARCHAR(255) NOT NULL,
    status       VARCHAR(255) NOT NULL,
    execute_time timestamptz  NOT NULL,
    executed_at  timestamptz,
    created_at   timestamptz  NOT NULL
);

CREATE INDEX tasks_status_index ON scheduler.tasks (status);
