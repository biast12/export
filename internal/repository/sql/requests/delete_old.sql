DELETE FROM requests
WHERE created_at < NOW() - $1::INTERVAL;