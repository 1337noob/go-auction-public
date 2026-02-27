CREATE TABLE auction.auction_read_model
(
    id           VARCHAR(255) PRIMARY KEY,
    lot_id       VARCHAR(255) NOT NULL,
    lot_name     VARCHAR(255) NOT NULL,
    start_price  INTEGER      NOT NULL,
    min_bid_step INTEGER      NOT NULL,
    seller_id    VARCHAR(255) NOT NULL,
    current_bid  JSONB,
    bids         JSONB,
    winner_id    VARCHAR(255),
    final_price  INTEGER,
    status       VARCHAR(255) NOT NULL,
    start_time   timestamptz  NOT NULL,
    end_time     timestamptz  NOT NULL,
    timeout      VARCHAR(255) NOT NULL,
    created_at   timestamptz  NOT NULL,
    started_at   timestamptz,
    completed_at timestamptz,
    updated_at   timestamptz  NOT NULL,
    version      INTEGER      NOT NULL
);
