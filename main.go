package resizenator

import (
	"context"
	"fmt"
	goimg "image"
	"log"
	"strings"
	"sync"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kalogs-c/resizenator/config"
	"github.com/kalogs-c/resizenator/gcs"
	"github.com/kalogs-c/resizenator/image"
	"github.com/kalogs-c/resizenator/semaphore"
)

func init() {
	functions.CloudEvent("Resizenator", Resize)
}

func Resize(ctx context.Context, e event.Event) error {
	log.Printf("Received event -> ID: %s\n", e.ID())

	var gcsData storagedata.StorageObjectData
	if err := protojson.Unmarshal(e.Data(), &gcsData); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w\n", err)
	}

	filename := gcsData.GetName()
	if is := image.IsImage(filename); !is {
		log.Println("Image is not an image")
		return nil
	}

	config, err := config.Load("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to load config: %w\n", err)
	}

	if strings.Contains(filename, config.Prefix) {
		log.Println("Image is already resized")
		return nil
	}

	storage, err := gcs.NewStorage(ctx, gcsData.GetBucket())
	if err != nil {
		return fmt.Errorf("Error creating storage %w\n", err)
	}

	reader, err := storage.Read(filename)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	srcImg, err := image.ReaderToImage(reader)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	desiredSizes := config.Sizes
	imageChan := make(chan goimg.Image, len(desiredSizes))
	wg := &sync.WaitGroup{}
	wg.Add(len(desiredSizes))
	uploadWg := &sync.WaitGroup{}
	uploadWg.Add(len(desiredSizes))

	resizeSemaphore := semaphore.NewSemaphore(config.MaxConcurrency)
	go func() {
		for _, size := range desiredSizes {
			resizeSemaphore.Acquire()
			if size.X == 0 || size.Y == 0 {
				wg.Done()
				uploadWg.Done()
				continue
			}
			log.Printf("Resizing image %s to %d x %d\n", filename, size.X, size.Y)
			go image.ResizeImageToChan(
				imageChan,
				resizeSemaphore,
				wg,
				srcImg,
				size,
				config.Algorithm,
			)
		}
	}()

	go image.ResizeMonitor(imageChan, wg)

	uploadSemaphore := semaphore.NewSemaphore(config.MaxConcurrency)
	go func() {
		for i := range imageChan {
			uploadSemaphore.Acquire()
			size := i.Bounds()
			imageName := fmt.Sprintf(
				"%s/%s%dx%d-%s",
				filename,
				config.Prefix,
				size.Dx(),
				size.Dy(),
				filename,
			)
			go storage.UploadImage(imageName, i, uploadWg, uploadSemaphore, config.TargetFormat)
		}
	}()

	uploadWg.Wait()

	if config.DeleteAfterUpload {
		storage.Delete(filename)
	}
	storage.Close()

	return nil
}

func main() {
	config, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if config.TargetFormat == "" {
		log.Printf("Config: %+v\n", config)
	}
}
