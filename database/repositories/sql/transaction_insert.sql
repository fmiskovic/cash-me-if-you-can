INSERT INTO transactions (account_id, amount, type, timestamp)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
RETURNING id, timestamp;