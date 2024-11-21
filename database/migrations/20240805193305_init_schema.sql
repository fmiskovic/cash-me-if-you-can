-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner VARCHAR(255) NOT NULL UNIQUE CHECK (TRIM(owner) <> ''),
    balance DECIMAL(38, 16) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    amount DECIMAL(38, 16) NOT NULL,
    type VARCHAR(15) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE INDEX idx_transactions_account_id ON transactions(account_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts, transactions;
-- +goose StatementEnd
