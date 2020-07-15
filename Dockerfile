# ベースとなるDockerイメージ指定
FROM golang:1.14

RUN go get github.com/labstack/echo/...
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/jinzhu/gorm
RUN go get github.com/go-redis/redis
RUN go get github.com/aws/aws-sdk-go/aws
# コンテナ内に作業ディレクトリを作成
RUN mkdir /app
# コンテナログイン時のディレクトリ指定
WORKDIR /app/src

ADD . /app