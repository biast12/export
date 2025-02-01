package transcriptstore

import (
	"context"
	"fmt"
	"github.com/TicketsBot/common/encryption"
	"github.com/TicketsBot/data-self-service/internal/config"
	"github.com/TicketsBot/data-self-service/internal/utils"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/sync/errgroup"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
)

type S3Client struct {
	logger *slog.Logger
	config config.WorkerConfig
	client *s3.Client
}

var _ Client = (*S3Client)(nil)

func NewS3Client(logger *slog.Logger, config config.WorkerConfig, client *s3.Client) *S3Client {
	return &S3Client{
		logger: logger,
		config: config,
		client: client,
	}
}

func (c *S3Client) GetTranscriptsForGuild(ctx context.Context, guildId uint64) (map[int][]byte, error) {
	keys, err := c.listForGuild(ctx, guildId)
	if err != nil {
		return nil, err
	}

	transcriptCount := 0
	for _, objKeys := range keys {
		transcriptCount += len(objKeys)
	}

	c.logger.Info("Found transcripts for guild", slog.Uint64("guild_id", guildId), slog.Int("transcript_count", transcriptCount))

	if transcriptCount > 500_000 {
		return nil, fmt.Errorf("too many transcripts for guild %d: %d", guildId, transcriptCount)
	}

	prefix := fmt.Sprintf("%d/", guildId)

	type object struct {
		bucket string
		key    string
	}

	keysCh := make(chan object)

	mu := sync.Mutex{}
	files := make(map[int][]byte)

	wrappedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, _ := errgroup.WithContext(wrappedCtx)
	for i := 0; i < c.config.Daemon.DownloadWorkers; i++ {
		group.Go(func() error {
			for objMetadata := range keysCh {
				if !strings.HasPrefix(objMetadata.key, prefix) {
					cancel()
					return fmt.Errorf("unexpected key: %s", objMetadata.key)
				}

				ticketId, err := strconv.Atoi(strings.TrimPrefix(objMetadata.key, prefix))
				if err != nil {
					cancel()
					return err
				}

				obj, err := c.client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: utils.Ptr(objMetadata.bucket),
					Key:    utils.Ptr(objMetadata.key),
				})
				if err != nil {
					cancel()
					return err
				}

				bytes, err := io.ReadAll(obj.Body)
				if err != nil {
					cancel()
					return err
				}

				mu.Lock()
				files[ticketId] = bytes
				c.logger.Debug("Downloaded transcript", slog.Uint64("guild_id", guildId), slog.Int("ticket_id", ticketId))
				mu.Unlock()
			}
			return nil
		})
	}

	for bucket, objKeys := range keys {
		for _, key := range objKeys {
			keysCh <- object{
				bucket: bucket,
				key:    key,
			}
		}
	}
	close(keysCh)

	if err := group.Wait(); err != nil {
		return nil, err
	}

	// Decompress + decrypt
	decryptedFiles := make(map[int][]byte)
	for ticketId, data := range files {
		decompressed, err := encryption.Decompress(data)
		if err != nil {
			return nil, err
		}

		decrypted, err := encryption.Decrypt([]byte(c.config.TranscriptS3.EncryptionKey), decompressed)
		if err != nil {
			return nil, err
		}

		decryptedFiles[ticketId] = decrypted
	}

	return decryptedFiles, nil
}

func (c *S3Client) listForGuild(ctx context.Context, guildId uint64) (map[string][]string, error) {
	prefix := fmt.Sprintf("%d/", guildId)

	keys := make(map[string][]string)

	for _, bucket := range c.config.TranscriptS3.Buckets {
		paginator := s3.NewListObjectsV2Paginator(c.client, &s3.ListObjectsV2Input{
			Bucket: utils.Ptr(bucket),
			Prefix: utils.Ptr(prefix),
		})

		bucketKeys := make([]string, 0)
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, obj := range output.Contents {
				if obj.Key == nil || !strings.HasPrefix(*obj.Key, prefix) {
					return nil, fmt.Errorf("unexpected key: %s", *obj.Key)
				}

				bucketKeys = append(bucketKeys, *obj.Key)
			}

			keys[bucket] = bucketKeys
		}
	}

	return keys, nil
}
