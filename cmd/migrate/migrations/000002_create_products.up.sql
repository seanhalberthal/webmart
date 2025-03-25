CREATE TABLE products
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID           NOT NULL,
    title       TEXT           NOT NULL,
    description TEXT,
    rating      INT            NOT NULL,
    price       DECIMAL(10, 2) NOT NULL,
    stock       INT            NOT NULL CHECK (stock >= 0),
    image_url   TEXT,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);