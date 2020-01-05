package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/minio/minio-go"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
)

func ToMinio(ctx context.Context, minioClient *minio.Client, file *os.File, args map[string]string) (string, error) {
	if args["folder_name"] == "" {
		return "", errors.New("desire folder_name cannot be empty")
	}

	if args["file_name"] == "" {
		return "", errors.New("desire file_name cannot be empty")
	}

	if args["aws_bucket"] == "" {
		return "", errors.New("aws_bucket cannot be empty")
	}

	bucketName := args["aws_bucket"]
	objectName := args["file_name"] + filepath.Ext(file.Name())
	filePath := fmt.Sprintf("/tmp/%s/%s", args["folder_name"], objectName)

	fmt.Println("Uploading to s3...")
	n, err := minioClient.FPutObjectWithContext(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: "multipart/form-data"})
	if err != nil {
		return "", err
	}

	// minioClient.PresignedGetObject(bucketName, objectName)
	fmt.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	return objectName, nil
}

func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader, args map[string]string) (string, error) {

	size := fileHeader.Size
	buffer := make([]byte, size)

	file.Read(buffer)
	tempFileName := args["folder_name"] + "/" + args["file_name"] + filepath.Ext(fileHeader.Filename)

	fmt.Println("Uploading to s3...")
	s3svc := s3manager.NewUploader(nil)
	result, err := s3svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(args["aws_bucket"]),
		Key:    aws.String(tempFileName),
		Body:   bytes.NewReader(buffer),
	})
	defer file.Close()

	fmt.Printf("%s/n", result.Location)
	// input := &s3.PutObjectInput{
	// 	ACL:                  aws.String("public-read"),
	// 	ContentLength:        aws.Int64(int64(size)),
	// 	ContentType:          aws.String(http.DetectContentType(buffer)),
	// 	ContentDisposition:   aws.String("attachment"),
	// 	ServerSideEncryption: aws.String("AES256"),
	// 	StorageClass:         aws.String("INTELLIGENT_TIERING"),
	// }
	// _, err := s3.New(s).PutObject(input)
	if err != nil {
		return "", err
	}

	return tempFileName, err
}
