CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       username VARCHAR(50) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(50) DEFAULT 'user',
                       coins INT NOT NULL DEFAULT 1000,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
                          product_id SERIAL PRIMARY KEY,
                          name VARCHAR(100) NOT NULL UNIQUE,
                          price INT NOT NULL,
                          quantity INT NOT NULL DEFAULT 2147483647
);

INSERT INTO products (name, price) VALUES
                                                 ('t-shirt', 80),
                                                 ('cup', 20),
                                                 ('book', 50),
                                                 ('pen', 10),
                                                 ('powerbank', 200),
                                                 ('hoody', 300),
                                                 ('umbrella', 200),
                                                 ('socks', 10),
                                                 ('wallet', 50),
                                                 ('pink-hoody', 500);

CREATE TABLE purchases (
                            purchase_id SERIAL PRIMARY KEY,
                            user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
                            product_id INT REFERENCES products(product_id) ON DELETE CASCADE,
                            quantity INT NOT NULL,
                            operation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE operations (
                            operation_id SERIAL PRIMARY KEY,
                            sender_user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
                            receiver_user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
                            amount INT NOT NULL,
                            operation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
                                token_id SERIAL PRIMARY KEY,
                                user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
                                token VARCHAR(255) NOT NULL,
                                expires_at TIMESTAMP NOT NULL,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);