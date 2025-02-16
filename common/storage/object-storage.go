package storage

import (
	"common/config"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/gofiber/fiber/v2"
	"log"
	"mime/multipart"
	"path"
)

type ObjectStorageInfo struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
	EndPoint  string
	BasePath  string
}
type ObjectStorage struct {
	Configuration *config.Config
	StorageInfo   *ObjectStorageInfo
	Client        *s3.Client
}

func NewObjectStorage(configuration *config.Config) *ObjectStorage {
	storageInfo := &ObjectStorageInfo{
		AccessKey: configuration.ObjectStorageAccessKey,
		SecretKey: configuration.ObjectStorageSecretKey,
		Region:    configuration.ObjectStorageRegion,
		Bucket:    configuration.ObjectStorageBucket,
		EndPoint:  configuration.ObjectStorageEndpoint,
		BasePath:  configuration.DocumentBasePath,
	}

	s3Config := aws.Config{
		Credentials:  credentials.NewStaticCredentialsProvider(storageInfo.AccessKey, storageInfo.SecretKey, ""),
		Region:       storageInfo.Region,
		BaseEndpoint: &storageInfo.EndPoint,
	}
	s3Client := s3.NewFromConfig(s3Config)

	return &ObjectStorage{
		Configuration: configuration,
		StorageInfo:   storageInfo,
		Client:        s3Client,
	}
}

func (obs *ObjectStorage) Upload(ctx *fiber.Ctx, username string, documents ...*multipart.FileHeader) error {
	for _, document := range documents {
		documentPath := fmt.Sprintf("%s/%s/%s", obs.StorageInfo.BasePath, username, document.Filename)
		if err := obs.uploadFile(ctx, documentPath, document); err != nil {
			log.Printf("Error uploading file %v: %v\n", document.Filename, err)
			return err
		}
	}
	return nil
}

func (obs *ObjectStorage) getPath(paths ...string) string {
	allPaths := append([]string{obs.StorageInfo.BasePath}, paths...)
	return path.Join(allPaths...)
}

func (obs *ObjectStorage) uploadFile(ctx *fiber.Ctx, documentPath string, document *multipart.FileHeader) error {
	file, err := document.Open()
	if err != nil {
		log.Printf("Couldn't open uploaded file %v. Here's why: %v\n", document.Filename, err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Error closing file %v: %v\n", document.Filename, cerr)
		}
	}()

	_, err = obs.Client.PutObject(ctx.Context(), &s3.PutObjectInput{
		Bucket: aws.String(obs.StorageInfo.Bucket),
		Key:    aws.String(documentPath),
		Body:   file,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			log.Printf("Error while uploading object to %s. The object is too large.\n"+
				"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
				"or the multipart upload API (5TB max).", obs.StorageInfo.Bucket)
		} else {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				document.Filename, obs.StorageInfo.Bucket, documentPath, err)
		}
	}
	return err
}
