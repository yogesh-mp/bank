CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(id),
    amount DECIMAL(15,2) NOT NULL,
    type TEXT CHECK (type IN ('deposit', 'withdraw', 'account_creation')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
