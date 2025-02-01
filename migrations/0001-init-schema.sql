CREATE TYPE request_type AS ENUM ('guild_transcripts', 'guild_data');
CREATE TYPE request_status AS ENUM ('queued', 'failed', 'completed');

CREATE TABLE requests
(
    id           uuid PRIMARY KEY        DEFAULT gen_random_uuid(),
    user_id      int8           NOT NULL,
    request_type request_type   NOT NULL,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT now(),
    guild_id     int8           NULL     DEFAULT NULL,
    status       request_status NOT NULL DEFAULT 'queued'
);

CREATE INDEX requests_user_id_idx ON requests (user_id);

CREATE TABLE task_queue
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id uuid NOT NULL,
    FOREIGN KEY (request_id) REFERENCES requests (id)
);

CREATE TABLE artifacts
(
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id uuid         NOT NULL UNIQUE,
    key        VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    size       int8         NOT NULL,
    FOREIGN KEY (request_id) REFERENCES requests (id)
);

CREATE INDEX artifacts_request_id_idx ON artifacts (request_id);
CREATE INDEX artifacts_expires_at_idx ON artifacts (expires_at);
