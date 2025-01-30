SELECT id, request_id, key, expires_at
FROM artifacts
WHERE request_id = $1;