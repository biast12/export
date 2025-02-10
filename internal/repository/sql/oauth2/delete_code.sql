DELETE FROM oauth2_codes
WHERE code = $1 AND client_id = $2;