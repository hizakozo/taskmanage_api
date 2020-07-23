package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
)

const S3Endpoint = "http://localhost:4572"

func S3PutObject(bucketName, key string, file multipart.File) (*s3.PutObjectOutput, error) {
	s := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("tekitou", "tekitou", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.ApNortheast1RegionID),
		Endpoint:         aws.String(S3Endpoint),
	}))

	c := s3.New(s, &aws.Config{})

	p := s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		ACL:    aws.String("public-read"),
		Body:   file,
	}

	r, err := c.PutObject(&p)

	return r, err
}