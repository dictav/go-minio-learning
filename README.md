# go-minio-learning

https://minio.io/ と https://github.com/minio/minio-go を使うサンプル

## Usage

### minio

事前に https://minio.io/ から実行ファイルをダウンロードする

```
$ minio server ./storage

AccessKey: YOUR_ACCESS_KEY  SecretKey: YOUR_SECRET_KEY
```

サーバーを起動すると `ACCESS_KEY` と `SECRET_KEY` が分かる


キーとホストを環境変数で設定するとバケットを作る。

```
$ env HOST="localhost:9000" ACCESS_KEY=YOUR_ACCESS_KEY SECRET_KEY=YOUR_SECRET_KEY go run minio.go
Success: I made a bucket.
PresignedPostPolicy:
curl ...
```

`PresignedPostPolicy:` 以下の curl コマンドを実行すると `post-object` というファイルをアップロードする。便利。

```
$ tree storage/
storage/
├── $buckets.json
├── $multiparts-session.json
└── go-minio-learning
    └── post-object
```

### AWS S3

minio は AWS S3 と API 互換性があるので、ソースコードそのままに AWS S3 の `ACCESS_KEY` と `SECRET_KEY` を設定すると同じように動作する。便利。

