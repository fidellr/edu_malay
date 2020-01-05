package usecase

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/fidellr/edu_malay/pict"
)

type pictUsecase struct {
	pictRepo       pict.PictRepository
	contextTimeout time.Duration
}

func NewPictUsecase(pr pict.PictRepository, timeout time.Duration) pict.PictUsecase {
	return &pictUsecase{
		pictRepo:       pr,
		contextTimeout: timeout,
	}
}

func (u *pictUsecase) Upload(c context.Context, pictFile *multipart.FileHeader, s3BucketInfo map[string]string) (fileName string, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	src, err := pictFile.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	folderName := s3BucketInfo["folder_name"]
	fileName = s3BucketInfo["file_name"]

	dir := os.TempDir() + "/" + folderName
	filePath := dir + "/" + fileName + filepath.Ext(pictFile.Filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	dst, err := createFile(filePath, src)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	fileName, err = u.pictRepo.Upload(ctx, dst, s3BucketInfo)
	if err != nil {
		return "", err
	}
	defer os.Remove(filePath)

	return fileName, err
}

func createFile(desireFileName string, sourceFile multipart.File) (dstFile *os.File, err error) {
	dstFile, err = os.Create(desireFileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(dstFile, sourceFile)
	if err != nil {
		return nil, err
	}

	return dstFile, nil
}
