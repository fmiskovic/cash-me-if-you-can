-- check if account exist
SELECT EXISTS (
    SELECT 1
    FROM accounts
    WHERE id = $1
);