SELECT client_id, client_secret, owner_id, label
FROM oauth2_clients
WHERE client_id = $1;