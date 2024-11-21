-- +goose Up
-- +goose StatementBegin
INSERT INTO accounts (id, owner, balance)
VALUES
    ('a1b2c3d4-1111-2222-3333-444455556666', 'Alice', 7467.8976),
    ('b1c2d3e4-2222-3333-4444-555566667777', 'Bob', 100.0000),
    ('c1d2e3f4-3333-4444-5555-666677778888', 'Charlie', 0.0000),
    ('2f6f112a-a8e2-42c3-a6b0-c15e86d01704', 'David', 0.0000);

INSERT INTO transactions (id, account_id, amount, type, timestamp)
VALUES
    ('a1b2c3d4-1111-2222-3333-444455556666', 'a1b2c3d4-1111-2222-3333-444455556666', 7467.8976, 'deposit', '2024-08-16 21:51:58'),
    ('b1c2d3e4-2222-3333-4444-555566667777', 'b1c2d3e4-2222-3333-4444-555566667777', 50.0000, 'deposit', '2024-08-16 21:51:58'),
    ('b1c2d3e4-2222-3333-4444-555566669999', 'b1c2d3e4-2222-3333-4444-555566667777', 50.0000, 'deposit', '2024-08-16 21:51:58'),
    ('c1d2e3f4-3333-4444-5555-666677778888', 'c1d2e3f4-3333-4444-5555-666677778888', 0.0000, 'deposit', '2024-08-16 21:51:58');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE accounts, transactions RESTART IDENTITY CASCADE;
-- +goose StatementEnd
