INSERT INTO requests (user_id, request_type, guild_id)
VALUES ($1, $2, $3)
RETURNING id, created_at;