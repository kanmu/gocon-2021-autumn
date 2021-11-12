# "Go" beyond your proxy

## 必要なもの

- Docker

## 問題

以下のイメージを実行すると 8000 番ポートで HTTP サーバが起動します。起動するサーバは本リポジトリに含まれる main.go の実装です。HTTP リクエストを投げて flag.txt に含まれる秘密の答えを見つけてください！Go で HTTP リクエストを送信するサンプルは `http_client_example` 配下に置いてありますので参考にしてください。

```
docker run --name go-beyond-your-proxy -p 8000:8000 --rm ghcr.io/kanmu/go-beyond-your-proxy:latest
```

**※注意事項**: docker exec や docker cp などを用いると簡単に flag.txt が取れてしまいます。本日は Go Conference です。Go のソースコードを読み解き、上記コマンドで公開された 8000 番ポートにアクセスする形で flag.txt の奪取に挑戦してみてください。

# ヒント

https://github.com/kanmu/gocon-2021-autumn/issues/1

まずはヒントとしてリバースプロキシや X-Forwarded-For、出題されたアプリケーションの説明をこちらの issue の概要欄に書いてありますので参考にしてみてください。ある程度時間が経過したらこちらの issue のコメント欄に徐々にヒントを投下していきます。自信がある方はヒントを見ずにやってみてください。

