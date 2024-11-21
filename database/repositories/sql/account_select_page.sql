SELECT a.id, a.owner, a.balance FROM accounts AS a ORDER BY a.created_at LIMIT $1 OFFSET $2;


