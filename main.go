package resizenator

import (
	"context"
	"fmt"
	goimg "image"
	"log"
	"sync"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"

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

	if md, ok := gcsData.GetMetadata()["created-by"]; ok {
		if md == "Resizenator" {
			log.Println("Image is already resized")
			return nil
		}
	}

	filename := gcsData.GetName()
	if is := image.IsImage(filename); !is {
		log.Println("Image is not an image")
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

	desiredSizes := make([]image.ImageSize, 4)
	imageChan := make(chan goimg.Image, 255)
	wg := &sync.WaitGroup{}
	wg.Add(len(desiredSizes))
	uploadWg := &sync.WaitGroup{}
	uploadWg.Add(len(desiredSizes))

	go func() {
		for i, size := range desiredSizes {
			size = image.ImageSize{X: i * 300, Y: i * 300}
			if size.X == 0 || size.Y == 0 {
				wg.Done()
				uploadWg.Done()
				continue
			}
			fmt.Printf("Resizing image %s to %d x %d\n", filename, size.X, size.Y)
			go image.ResizeImageToChan(imageChan, wg, srcImg, size, "BiLinear")
		}
	}()

	go image.ResizeMonitor(imageChan, wg)

	semaphore := semaphore.NewSemaphore(10)
	go func() {
		for i := range imageChan {
			semaphore.Acquire()

			size := i.Bounds()
			imageName := fmt.Sprintf("%s/%dx%d-%s", filename, size.Dx(), size.Dy(), filename)
			go storage.UploadImage(imageName, i, uploadWg, semaphore)
		}
	}()

	uploadWg.Wait()

	storage.Delete(filename)
	storage.Close()

	return nil
}
