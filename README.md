# タスク管理アプリケーションのAPI
※local mailサーバーは立てれていないので、環境変数を除外しました。よってmail機能は動きません。

### アーキテクチャ


### 主な機能
- 複数ユーザーとプロジェクトを共有できる
- プロジェクトの作成、削除
- チケット作成、削除、編集、コメント、進捗度（status）の変更
- statusの作成、編集、削除

### フロントエンド（Vue.js）
https://github.com/hizakozo/taskmanage_front

### docker start
`docker-compose up -d`
### app start
`GO_ENV=local go run main.go`
