SELECT t.id, t.account_id, t.type, t.amount, t.timestamp
FROM transactions AS t
WHERE t.id = $1;