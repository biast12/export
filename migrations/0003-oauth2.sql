CREATE TABLE oauth2_clients (
    client_id VARCHAR(32) NOT NULL PRIMARY KEY,
    client_secret VARCHAR(64) NOT NULL,
    label VARCHAR(255)
);

CREATE TABLE oauth2_redirect_uris (
    client_id VARCHAR(32) NOT NULL,
    redirect_uri VARCHAR(255) NOT NULL,
    PRIMARY KEY (client_id, redirect_uri),
    FOREIGN KEY (client_id) REFERENCES oauth2_clients(client_id)
);

CREATE INDEX oauth2_redirect_uris_client_id ON oauth2_redirect_uris(client_id);

CREATE TABLE oauth2_authorized_clients (
    user_id BIGINT NOT NULL,
    client_id VARCHAR(32) NOT NULL,
    scopes VARCHAR(255) NOT NULL,
    PRIMARY KEY (user_id, client_id)
);

CREATE INDEX oauth2_authorized_clients_user_id ON oauth2_authorized_clients(user_id);
