UPDATE oauth2_clients
SET client_secret = $2
WHERE client_id = $1;