SELECT COALESCE(SUM (size), 0)
FROM artifacts
INNER JOIN downloads ON artifacts.id = downloads.artifact_id
WHERE downloads.download_time > NOW() - INTERVAL '1 DAY';