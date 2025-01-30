SELECT task_queue.id, task_queue.request_id, requests.id, requests.user_id, requests.request_type, requests.created_at, requests.guild_id, requests.status
FROM task_queue
INNER JOIN requests ON task_queue.request_id = requests.id
WHERE requests.status = 'queued'
ORDER BY requests.created_at ASC
LIMIT 1