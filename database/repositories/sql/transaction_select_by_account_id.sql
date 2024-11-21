SELECT t.id, t.account_id, t.type, t.amount, t.timestamp
FROM transactions AS t
WHERE t.account_id = $1
ORDER BY t.timestamp DESC;