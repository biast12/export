SELECT
    requests.id, requests.user_id, requests.request_type, requests.created_at, requests.guild_id, requests.status,
    artifacts.id, artifacts.request_id, artifacts.key, artifacts.expires_at
FROM requests
LEFT OUTER JOIN artifacts ON requests.id = artifacts.request_id
WHERE requests.user_id = $1;