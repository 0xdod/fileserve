package filestorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Backend struct {
	client     *s3.S3
	bucketName string
	region     string
}

type S3BackendConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
}

func (cfg *S3BackendConfig) Validate() error {
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Region == "" || cfg.BucketName == "" {
		return errors.New("Invalid s3 configuration")
	}

	return nil
}

func NewS3StorageBackend(cfg *S3BackendConfig) *s3Backend {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	}))

	return &s3Backend{
		client:     s3.New(sess),
		bucketName: cfg.BucketName,
		region:     cfg.Region,
	}
}

func (s *s3Backend) Upload(ctx context.Context, param UploadParam) (string, error) {
	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(param.Name),
		Body:        bytes.NewReader(param.Content),
		ContentType: aws.String(http.DetectContentType(param.Content)),
	})

	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", s.region, s.bucketName, param.Name)
	fmt.Printf("File successfully uploaded to %s\n", endpoint)

	return endpoint, nil
}

func (s *s3Backend) Download(ctx context.Context, name string) ([]byte, error) {
	res, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &name,
	})

	if err != nil {
		return nil, err
	}

	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(res.Body)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
