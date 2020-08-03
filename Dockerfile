# ベースとなるDockerイメージ指定
FROM golang:1.14
# コンテナ内に作業ディレクトリを作成
RUN mkdir /app
# コンテナログイン時のディレクトリ指定
WORKDIR /app

ADD . /app