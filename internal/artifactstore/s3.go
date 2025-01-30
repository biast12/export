package artifactstore

import (
	"bytes"
	"context"
	"fmt"
	"github.com/TicketsBot/common/encryption"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"io"
	"math/rand"
)

type S3ArtifactStore struct {
	client        *s3.Client
	bucketName    string
	encryptionKey []byte
}

var _ ArtifactStore = (*S3ArtifactStore)(nil)

func NewS3ArtifactStore(client *s3.Client, bucketName string, encryptionKey []byte) *S3ArtifactStore {
	return &S3ArtifactStore{
		client:        client,
		bucketName:    bucketName,
		encryptionKey: encryptionKey,
	}
}

func (s *S3ArtifactStore) Fetch(ctx context.Context, requestId uuid.UUID, key string) ([]byte, error) {
	objectKey := fmt.Sprintf("%s/%s", requestId, key)

	opts := &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &objectKey,
	}

	obj, err := s.client.GetObject(ctx, opts)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	decrypted, err := encryption.Decrypt(s.encryptionKey, data)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func (s *S3ArtifactStore) Store(ctx context.Context, requestId uuid.UUID, key string, data []byte) error {
	// Encrypt data first
	encrypted, err := encryption.Encrypt(s.encryptionKey, data)
	if err != nil {
		return err
	}

	objectKey := fmt.Sprintf("%s/%s", requestId, key)

	opts := &s3.PutObjectInput{
		Bucket: &s.bucketName,
		Key:    &objectKey,
		Body:   bytes.NewReader(encrypted),
	}

	_, err = s.client.PutObject(ctx, opts)
	return err
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomKey(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(bytes)
}
