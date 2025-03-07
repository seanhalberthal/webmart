CREATE TABLE products
(
    id          UUID PRIMARY KEY,
    user_id     UUID           NOT NULL,
    name        TEXT           NOT NULL,
    description TEXT,
    price       DECIMAL(10, 2) NOT NULL,
    stock       INT            NOT NULL CHECK (stock >= 0),
    image_url   TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);