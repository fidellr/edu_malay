package repository

import (
	"context"
	"os"

	"github.com/minio/minio-go"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/fidellr/edu_malay/utils"
)

type S3Pict struct {
	MinioClient *minio.Client
	S3Cfg       *aws.Config
}

func NewS3Pict(minioClient *minio.Client) *S3Pict {
	return &S3Pict{
		MinioClient: minioClient,
	}
}

func (r *S3Pict) Upload(ctx context.Context, picture *os.File, bucketInfo map[string]string) (fileName string, err error) {
	// file, err := picture.Open()
	// if err != nil {
	// 	return "", err
	// }

	fileName, err = utils.ToMinio(ctx, r.MinioClient, picture, bucketInfo)
	if err != nil {
		return "", err
	}
	// fileName, err = utils.UploadFileToS3(sess, file, picture, bucketInfo, r.S3Cfg)
	// if err != nil {
	// 	return "", err
	// }

	return fileName, err
}
