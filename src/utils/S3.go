package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"taskmanage_api/src/constants"
)

func S3PutObject(key string, file multipart.File) (*s3.PutObjectOutput, error) {
	s := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(constants.Params.S3Key, constants.Params.S3SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.ApNortheast1RegionID),
		Endpoint:         aws.String(constants.Params.S3EndPoint),
	}))

	c := s3.New(s, &aws.Config{})

	p := s3.PutObjectInput{
		Bucket: aws.String(constants.Params.S3BucketName),
		Key:    aws.String(key),
		ACL:    aws.String("public-read"),
		Body:   file,
	}

	r, err := c.PutObject(&p)

	return r, err
}
