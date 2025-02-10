CREATE TABLE oauth2_clients
(
    client_id     VARCHAR(32) NOT NULL PRIMARY KEY,
    client_secret VARCHAR(64) NOT NULL,
    owner_id      BIGINT      NOT NULL,
    label         VARCHAR(255)
);

CREATE TABLE oauth2_redirect_uris
(
    client_id    VARCHAR(32)  NOT NULL,
    redirect_uri VARCHAR(255) NOT NULL,
    PRIMARY KEY (client_id, redirect_uri),
    FOREIGN KEY (client_id) REFERENCES oauth2_clients (client_id) ON DELETE CASCADE
);

CREATE INDEX oauth2_redirect_uris_client_id ON oauth2_redirect_uris (client_id);

CREATE TABLE oauth2_authorized_clients
(
    user_id   BIGINT       NOT NULL,
    client_id VARCHAR(32)  NOT NULL,
    scopes    VARCHAR(255) NOT NULL,
    PRIMARY KEY (user_id, client_id),
    FOREIGN KEY (client_id) REFERENCES oauth2_clients (client_id) ON DELETE CASCADE
);

CREATE INDEX oauth2_authorized_clients_user_id ON oauth2_authorized_clients (user_id);

CREATE TABLE oauth2_codes
(
    code       VARCHAR(32) NOT NULL,
    client_id  VARCHAR(32) NOT NULL,
    user_id    BIGINT      NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (code),
    FOREIGN KEY (client_id) REFERENCES oauth2_clients (client_id) ON DELETE CASCADE
);

CREATE TABLE oauth2_code_authorities
(
    code      VARCHAR(32)  NOT NULL,
    authority VARCHAR(255) NOT NULL,
    PRIMARY KEY (code, authority),
    FOREIGN KEY (code) REFERENCES oauth2_codes (code) ON DELETE CASCADE
);