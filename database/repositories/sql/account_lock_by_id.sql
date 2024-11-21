SELECT id, owner, balance FROM accounts WHERE id = $1 FOR UPDATE;
-- Locks the selected row for update.