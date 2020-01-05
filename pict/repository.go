package pict

import (
	"context"
	"os"
)

type PictRepository interface {
	Upload(ctx context.Context, picture *os.File, bucketInfo map[string]string) (fileName string, err error)
	// Download(ctx context.Context) error
}
