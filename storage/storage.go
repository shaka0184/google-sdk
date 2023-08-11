package storage

import (
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"time"
)

func GetReader(ctx context.Context, bucketName, objectPath string) (*storage.Reader, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	obj := client.Bucket(bucketName).Object(objectPath)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return reader, nil
}

func GetByteSlice(ctx context.Context, bucketName, objectPath string) ([]byte, error) {
	r, err := GetReader(ctx, bucketName, objectPath)
	if err != nil {
		return nil, err
	}

	res, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer r.Close()

	return res, nil
}

func Upload(ctx context.Context, bucketName, objectPath string, f *os.File) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	obj := client.Bucket(bucketName).Object(objectPath)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer reader.Close()

	// 書き込み
	tee := io.TeeReader(reader, f)
	s := bufio.NewScanner(tee)
	for s.Scan() {
	}
	if err := s.Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func UploadFile(bucket, object string, f io.Reader) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
