package cms

import (
	"context"
	"encoding/base64"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type FileManagerConfig struct {
	Bucket string
}

func (fmc *FileManagerConfig) SetDefaults() {
	if fmc.Bucket == "" {
		fmc.Bucket = "uncyclo-dynamic-assets"
	}
}

type FileManager struct {
	logger *zap.SugaredLogger
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewFileManager(lc fx.Lifecycle, logger *zap.SugaredLogger, config *FileManagerConfig) (*FileManager, error) {
	ctx := context.Background()

	if creds := os.Getenv("GOOGLE_CREDENTIALS"); creds != "" {
		tmpFile, err := os.CreateTemp("", "google-credentials-*.json")

		if err != nil {
			return nil, err
		}

		data, err := base64.StdEncoding.DecodeString(creds)

		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(tmpFile.Name(), data, 0600); err != nil {
			return nil, err
		}

		if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", tmpFile.Name()); err != nil {
			return nil, err
		}
	}

	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(config.Bucket)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return &FileManager{
		logger: logger.Named("file-manager"),

		client: client,
		bucket: bucket,
	}, nil
}

func (fm *FileManager) OpenWriter(ctx context.Context, name string) io.WriteCloser {
	obj := fm.bucket.Object(name)
	wc := obj.NewWriter(ctx)

	return wc
}

func (fm *FileManager) WriteFile(ctx context.Context, name string, reader io.Reader) error {
	obj := fm.bucket.Object(name)
	wc := obj.NewWriter(ctx)

	if _, err := io.Copy(wc, reader); err != nil {
		return err
	}

	return wc.Close()
}

func (fm *FileManager) Close() error {
	return fm.client.Close()
}

func (fm *FileManager) Rename(ctx context.Context, from string, to string) error {
	fromObject := fm.bucket.Object(from)
	toObject := fm.bucket.Object(to)

	_, err := toObject.CopierFrom(fromObject).Run(ctx)

	if err != nil {
		fm.logger.Warn(err)
	}

	return fromObject.Delete(ctx)
}
