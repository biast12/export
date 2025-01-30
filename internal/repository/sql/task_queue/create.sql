INSERT INTO task_queue (request_id)
VALUES ($1)
RETURNING "id";