CREATE TABLE auction.event_store
(
    id           VARCHAR(255) PRIMARY KEY,
    aggregate_id VARCHAR(255) NOT NULL,
    version      INTEGER      NOT NULL,
    event_type   VARCHAR(255) NOT NULL,
    event_data   JSONB        NOT NULL,
    timestamp    timestamptz  NOT NULL
);

CREATE INDEX auction_event_store_aggregate_id_index ON auction.event_store (aggregate_id);