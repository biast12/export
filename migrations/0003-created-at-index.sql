CREATE INDEX requests_created_at_idx ON requests(created_at);

ALTER TABLE task_queue DROP CONSTRAINT task_queue_request_id_fkey;
ALTER TABLE task_queue ADD CONSTRAINT task_queue_request_id_fkey
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE CASCADE;

ALTER TABLE artifacts DROP CONSTRAINT artifacts_request_id_fkey;
ALTER TABLE artifacts ADD CONSTRAINT artifacts_request_id_fkey
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE CASCADE;

ALTER TABLE downloads DROP CONSTRAINT downloads_artifact_id_fkey;
ALTER TABLE downloads ADD CONSTRAINT downloads_artifact_id_fkey
    FOREIGN KEY (artifact_id) REFERENCES artifacts(id) ON DELETE CASCADE;