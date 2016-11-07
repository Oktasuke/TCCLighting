# TCCLighting
FaceBookページ等に登録したWebHookを受け取りごにょごにょするレシーバー

## 設定
1. facebookアプリを作成(デベロッパーモードだとWebHookは受け取れないので注意)
1. Graph ExplorerでウォッチしたいPageのTokenを生成
1. Graph ExplorerをPOSTリクエストに変更し,
[リクエストフォーマット](https://developers.facebook.com/docs/graph-api/reference/v2.8/app/subscriptions/)をセットして`{page_id}/subscribed_apps`にリクエストする。（対象ページのAdmin権限を保持している必要がある)

| KEY | VALUE |
|:----------|:----------|
| object | page |
| callback_url | このアプリをデプロイした環境のURL |
| fields | [checkins,feed] |
| active | true |

## トリガー
* 現状
	* ページへのいいね追加
	* フィード（ページに投稿された記事）へのいいね追加
	* フィード（ページに投稿された記事）へのリアクション（超いいねなど）追加
	* フィード（ページに投稿された記事）へのコメント追加
	* ビジター記事へのいいね追加
	* ビジター記事へのコメント追加
* 追加予定
	* ページのプレイスへのチェックイン

## デバック用
* ページにいいねがついた時

```sh
curl -v -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"entry": [{"changes": [{"field": "feed","value": {"item": "like","verb": "add","user_id": 1186329501459477}}],"id": "1142469655832956","time": 1478439821}],"object": "page"}'  http://localhost:8080/facebook
```

* 記事にリアクションが追加された時

```sh
curl -v -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"entry": [{"changes": [{"field": "feed","value": {"parent_id": "1142469655832956_1164870100259578","sender_id": 1142469655832956,"item": "reaction","verb": "add","created_time": 1478440417,"post_id": "1142469655832956_1164870100259578"}}],"id": "1142469655832956","time": 1478440417}],"object": "page"}'  http://localhost:8080/facebook
```
* 記事にいいねがついた時

```sh
curl -v -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"entry": [{"changes": [{"field": "feed","value": {	"parent_id": "1142469655832956_1163393020407286",	"sender_name": "CTCC",	"sender_id": 1142469655832956,	"item": "like",	"verb": "add",	"created_time": 1478439308,	"post_id": "1142469655832956_1163393020407286"}}],"id": "1142469655832956","time": 1478439308}],"object": "page"}' http://localhost:8080/facebook
```

* 記事にコメントがついた時

```sh
curl -v -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"entry": [{"changes": [{"field": "feed","value": {"parent_id": "1142469655832956_1163393020407286","sender_name": "CTCC","comment_id": "1163393020407286_1164864096926845","sender_id": 1142469655832956,"item": "comment","verb": "add","created_time": 1478439630,"post_id": "1142469655832956_1163393020407286","message": "test"}}],"id": "1142469655832956","time": 1478439630}],"object": "page"}' http://localhost:8080/facebook
```

* 記事からいいねが削除された時

```sh
curl -v -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"entry": [{"changes": [{"field": "feed","value": {"parent_id": "1142469655832956_1163393020407286","sender_id": 1142469655832956,"item": "like","verb": "remove","created_time": 1478439307,"post_id": "1142469655832956_1163393020407286"}}],"id": "1142469655832956","time": 1478439307}],"object": "page"}' http://localhost:8080/facebook
```

## 参考
* http://ukimiawz.github.io/facebook/2015/08/11/facebook-page-subscribed-apps/
* http://ukimiawz.github.io/facebook/2015/08/12/webhook-facebook-subscriptions/
