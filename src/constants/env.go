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

	fmt.Println(Params)
}
