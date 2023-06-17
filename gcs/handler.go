package gcs

import (
	"context"
	goimg "image"
	"io"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/storage"

	"github.com/kalogs-c/resizenator/image"
	"github.com/kalogs-c/resizenator/semaphore"
)

type Storage struct {
	bucket *storage.BucketHandle
	ctx    context.Context
}

func NewStorage(ctx context.Context, bucketName string) (*Storage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)
	return &Storage{bucket: bucket, ctx: ctx}, nil
}

func (s *Storage) Read(filename string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	reader, err := s.bucket.Object(filename).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (s *Storage) UploadImage(
	filename string,
	img goimg.Image,
	wg *sync.WaitGroup,
	semaphore *semaphore.Semaphore,
) {
	defer func() {
		wg.Done()
		semaphore.Release()
	}()

	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	writer := s.bucket.Object(filename).NewWriter(ctx)
	err := image.ImageToWriter(img, writer, image.GetImageFormat(filename))
	if err != nil {
		log.Fatalf("encoding image error: %v", err)
	}

	if err := writer.Close(); err != nil {
		log.Fatalf("writer close error: %v", err)
	}
}
