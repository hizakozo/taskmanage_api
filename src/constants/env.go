package constants

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type EnvParam struct {
	DbUrl     string
	DbUser    string
	DbPass    string
	RedisAddr string
	FrontUrl  string
	S3EndPoint string
	S3Key string
	S3SecretKey string
	S3BucketName string
	MailUserName string
	MailPassword string
}

var Params EnvParam

func init() {
	if err := godotenv.Load(fmt.Sprintf("env/%s.env", os.Getenv("GO_ENV"))); err != nil {
	}
	Params.DbUrl = os.Getenv("DB_URL")
	Params.DbUser = os.Getenv("DB_USER")
	Params.DbPass = os.Getenv("DB_PASS")
	Params.RedisAddr = os.Getenv("REDIS_ADDR")
	Params.FrontUrl = os.Getenv("FRONT_URL")
	Params.S3EndPoint = os.Getenv("S3_ENDPOINT")
	Params.S3Key = os.Getenv("S3_KEY")
	Params.S3SecretKey = os.Getenv("S3_SECRET_KEY")
	Params.S3BucketName = os.Getenv("S3_BUCKET_NAME")
	Params.MailUserName = os.Getenv("MAIL_USER_NAME")
	Params.MailPassword = os.Getenv("MAIL_PASSWORD")
}
