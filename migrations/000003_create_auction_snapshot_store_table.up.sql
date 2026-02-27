CREATE TABLE auction.snapshot_store
(
    aggregate_id VARCHAR(255) PRIMARY KEY,
    version      INTEGER     NOT NULL,
    data         JSONB,
    created_at   timestamptz NOT NULL
);

CREATE INDEX auction_snapshot_store_aggregate_id_index ON auction.snapshot_store (aggregate_id);