package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/0xdod/fileserve"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3CloudStorage struct {
	client     *s3.S3
	bucketName string
	region     string
}

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
}

func (cfg *Config) Validate() error {
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Region == "" || cfg.BucketName == "" {
		return errors.New("Invalid s3 configuration")
	}

	return nil
}

func NewFileStorage(cfg *Config) *s3CloudStorage {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	}))

	client := s3.New(sess)

	return &s3CloudStorage{
		client:     client,
		bucketName: cfg.BucketName,
		region:     cfg.Region,
	}
}

func (s *s3CloudStorage) Upload(ctx context.Context, param fileserve.UploadParam) (string, error) {
	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(param.Name),
		Body:        bytes.NewReader(param.Content),
		ContentType: aws.String(http.DetectContentType(param.Content)),
	})

	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, strings.ReplaceAll(param.Name, " ", "+"))
	fmt.Printf("File successfully uploaded to %s\n", endpoint)

	return endpoint, nil
}

func (s *s3CloudStorage) Download(ctx context.Context, name string) ([]byte, error) {
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

func (s *s3CloudStorage) createBucketWithPublicReadPolicy() {
	// Check if the bucket exists
	if _, err := s.client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	}); err != nil {
		if err != nil {
			// If the bucket doesn't exist, create it
			_, err := s.client.CreateBucket(&s3.CreateBucketInput{
				Bucket: aws.String(s.bucketName),
			})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Bucket %s created\n", s.bucketName)

			// Define the bucket policy for public read access
			policy := `{
		"Version":"2012-10-17",
		"Statement":[{
			"Sid":"PublicReadGetObject",
			"Effect":"Allow",
			"Principal": "*",
			"Action":["s3:GetObject"],
			"Resource":["arn:aws:s3:::` + s.bucketName + `/*"]
		}]
	}`

			// Set the bucket policy
			_, err = s.client.PutBucketPolicy(&s3.PutBucketPolicyInput{
				Bucket: aws.String(s.bucketName),
				Policy: aws.String(policy),
			})
			if err != nil {
				panic(err)
			}
			fmt.Printf("Bucket policy set for %s to allow public read access\n", s.bucketName)

		}
	}
}
