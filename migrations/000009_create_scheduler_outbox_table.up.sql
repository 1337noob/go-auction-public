CREATE TABLE scheduler.outbox
(
    id             VARCHAR(255) PRIMARY KEY,
    aggregate_id   VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    event_type     VARCHAR(255) NOT NULL,
    event_data     JSONB        NOT NULL,
    status         VARCHAR(255) NOT NULL,
    created_at     timestamptz  NOT NULL
);

CREATE INDEX outbox_status_index ON scheduler.outbox (status);