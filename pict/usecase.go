package pict

import (
	"context"
	"mime/multipart"
)

type PictUsecase interface {
	Upload(ctx context.Context, picture *multipart.FileHeader, bucketInfo map[string]string) (fileName string, err error)
	// Download(ctx context.Context) error
}
