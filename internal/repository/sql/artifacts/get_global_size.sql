SELECT sum(size)
FROM artifacts
WHERE expires_at > NOW();