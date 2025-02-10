SELECT EXISTS (
    SELECT 1
    FROM oauth2_redirect_uris
    WHERE client_id = $1 AND redirect_uri = $2
);