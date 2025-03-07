ALTER TABLE products
    ADD CONSTRAINT fk_user foreign key (user_id) REFERENCES users (id);