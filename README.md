# タスク管理アプリケーションのAPI
※local mailサーバーは立てれていないので、環境変数を除外しました。なのでmail機能は動きません。

### アーキテクチャ
![001](https://user-images.githubusercontent.com/47819815/89103672-72fde180-d44e-11ea-88d8-541a3e9fe64c.jpeg)

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
