package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"object-storage/internal/config"
	"object-storage/internal/entity"
	"object-storage/pkg/filereader"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type ObjectStorage struct {
	store jetstream.ObjectStore
}

func New(nc *nats.Conn, config *config.Config) (*ObjectStorage, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	store, err := js.ObjectStore(ctx, config.Nats.Bucket)
	if err != nil {
		if errors.Is(err, jetstream.ErrBucketNotFound) {
			if store, err = js.CreateObjectStore(ctx, jetstream.ObjectStoreConfig{
				Bucket: config.Nats.Bucket,
			}); err != nil {
				return nil, err
			}
		}
	}

	return &ObjectStorage{
		store: store,
	}, nil
}

func (o *ObjectStorage) Put(ctx context.Context, reader filereader.Reader) (*entity.FileInfo, error) {
	fileMap := strings.Split(reader.Filename(), ".")

	ext := fileMap[len(fileMap)-1]

	filename := fmt.Sprintf("%s.%s", uuid.New().String(), ext)

	info, err := o.store.Put(ctx, jetstream.ObjectMeta{
		Name: filename,
	}, reader)
	if err != nil {
		return nil, err
	}

	slog.Info("putted object", slog.Any("info", info))

	return &entity.FileInfo{
		Id:       info.NUID,
		Filename: info.Name,
	}, nil
}

func (o *ObjectStorage) Get(ctx context.Context, filename string) (*entity.FileResult, uint64, error) {
	res, err := o.store.Get(ctx, filename)
	if err != nil {
		return nil, 0, err
	}

	info, err := res.Info()
	if err != nil {
		return nil, 0, err
	}

	return &entity.FileResult{
		ReadCloser: res,
	}, info.Size, nil
}
