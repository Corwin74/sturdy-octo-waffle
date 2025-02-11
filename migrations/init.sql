CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance INT CHECK(balance BETWEEN 0 AND 3000000) NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    price INT CHECK(price BETWEEN 0 AND 1000000) NOT NULL
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
    ('pink-hoody', 500)
ON CONFLICT (name) DO NOTHING;



CREATE TABLE IF NOT EXISTS users_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_user UUID NOT NULL REFERENCES users(id) ,
    id_product UUID NOT NULL REFERENCES products(id)  
);

CREATE TABLE IF NOT EXISTS transfers_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_user UUID NOT NULL REFERENCES users(id) ,
    amount INT CHECK(amount BETWEEN 0 AND 3000000) NOT NULL
);
