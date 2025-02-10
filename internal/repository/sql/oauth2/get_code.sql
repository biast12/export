SELECT code, client_id, user_id, created_at
FROM oauth2_codes
WHERE code = $1 AND client_id = $2 AND created_at > NOW() - $3::interval
FOR UPDATE;