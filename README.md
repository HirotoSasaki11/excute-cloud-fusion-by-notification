## 概要
Cloud Storageにファイルがアップロードされたのをトリガーに、Data Fusionのパイプラインを実行するためのサンプルコード

## 事前準備
パイプラインを実行するためのAPIには、アクセストークンを渡す必要がある。しかしアクセストークンの有効期限は短い。
そのため、あらかじめGCP環境にてクライアントIDを生成して、アクセストークンを取得するためのリフレッシュトークンを取得しておく必要がある。

## deploy
```console
gcloud --project ${project_id} functions deploy Excute \
  --runtime go113 \
  --entry-point Excute \
  --trigger-resource {bucket} \
  --region asia-northeast1 \
  --set-env-vars="ENDPOINT"="${endpoint}" \
  --set-env-vars="CLIENTID"="${client_id}" \
  --set-env-vars="CLIENTSECRET"="${client_secret}" \
  --set-env-vars="RERFESH_TOKEN"="${refresh_token}" \
  --set-env-vars="PIPLINENAME"="${pipline_name}" \
  --trigger-event google.storage.object.finalize
```