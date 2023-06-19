package gcs

import (
	"bytes"
	"context"
	goimg "image"
	"io"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/kalogs-c/resizenator/image"
	"github.com/kalogs-c/resizenator/semaphore"
)

type Storage struct {
	client *storage.Client
	bucket *storage.BucketHandle
	ctx    context.Context
}

func NewStorage(ctx context.Context, bucketName string) (*Storage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)
	return &Storage{client: client, bucket: bucket, ctx: ctx}, nil
}

func (s *Storage) ListObjects() ([]string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	query := &storage.Query{Prefix: ""}

	var names []string
	it := s.bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return names, err
		}
		names = append(names, attrs.Name)
	}

	return names, nil
}

func (s *Storage) Read(filename string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	reader, err := s.bucket.Object(filename).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func (s *Storage) UploadImage(
	filename string,
	img goimg.Image,
	wg *sync.WaitGroup,
	semaphore semaphore.Semaphore,
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
	log.Printf("Uploading image %s to bucket\n", filename)

	if err := writer.Close(); err != nil {
		log.Fatalf("writer close error: %v", err)
	}
	log.Printf("Uploaded image %s to bucket\n", filename)
}

func (s *Storage) Delete(filename string) error {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	err := s.bucket.Object(filename).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() error {
	return s.client.Close()
}
