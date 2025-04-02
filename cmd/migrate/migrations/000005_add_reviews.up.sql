CREATE TABLE IF NOT EXISTS reviews
(
    id         UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    user_id    UUID NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
