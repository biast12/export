CREATE TABLE downloads (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    artifact_id UUID NOT NULL REFERENCES artifacts(id),
    download_time TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX downloads_user_id_idx ON downloads(user_id);
CREATE INDEX downloads_download_time_idx ON downloads(download_time);