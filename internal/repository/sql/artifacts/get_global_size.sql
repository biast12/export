SELECT COALESCE(SUM(size), 0)
FROM artifacts
WHERE expires_at > NOW();