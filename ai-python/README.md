# 部屋の散らかり度判定 API

CLIPとYOLOv11を使って寝室の散らかり具合を判定するFastAPI サーバーです。

---

## 必要要件

| ツール | バージョン |
|--------|-----------|
| Docker | 20.10以上 |
| Docker Compose (任意) | v2以上 |

> **注意:** CLIPモデル（ViT-L/14）とYOLOv11lは合計で約3GB程度のディスクを消費します。初回起動時にダウンロードされます。

---

## ディレクトリ構成

```
project/
├── main.py       # アプリケーション本体
├── Dockerfile    # Dockerイメージ定義
└── README.md     # このファイル
```

---

## セットアップ手順

### 1. ファイルの配置

`main.py` と `Dockerfile` を同じディレクトリに置きます。

```
mkdir room-analyzer
cd room-analyzer
# main.py と Dockerfile をここに配置
```

### 2. Dockerイメージのビルド

```bash
docker build -t room-analyzer .
```

> 初回ビルドはCLIPモデル・YOLOモデルのダウンロードがあるため、**10〜20分程度**かかります。

### 3. コンテナの起動

```bash
docker run -p 8000:8000 room-analyzer
```

起動後、以下のログが出れば成功です：

```
INFO:     Uvicorn running on http://0.0.0.0:8000 (Press CTRL+C to quit)
```

---

## APIの使い方

### ヘルスチェック

```bash
curl http://localhost:8000/
```

**レスポンス例:**
```json
{"message": "Hello World"}
```

### 画像の分析

```bash
curl -X POST http://localhost:8000/analyze/ \
  -F "file=@/path/to/your/bedroom.jpg"
```

**レスポンス例:**
```json
{
  "image": "bedroom.jpg",
  "clipv4_messiness_score": 42.3,
  "yolo": {
    "bottle": 2,
    "cup": 1
  }
}
```

#### レスポンスフィールドの説明

| フィールド | 説明 |
|-----------|------|
| `image` | アップロードされたファイル名 |
| `clipv4_messiness_score` | 散らかり度スコア（0〜100、高いほど散らかっている） |
| `yolo` | 検出されたオブジェクトと個数（bottle, cup, wine glass, bowl, can） |

#### スコアの目安

| スコア | 状態 |
|--------|------|
| 0〜25 | きれいに整頓されている |
| 25〜50 | 普通の生活感がある状態 |
| 50〜75 | やや散らかっている |
| 75〜100 | かなり散らかっている |

---

## ブラウザでAPIを確認する

Swagger UIが自動生成されます：

```
http://localhost:8000/docs
```

---

## コンテナの停止

```bash
# 起動中のコンテナIDを確認
docker ps

# 停止
docker stop <CONTAINER_ID>
```

---

## トラブルシューティング

### ビルド時にエラーが出る

```bash
# Dockerのディスク容量を確認
docker system df

# 不要なキャッシュを削除してから再ビルド
docker system prune
docker build --no-cache -t room-analyzer .
```

### メモリ不足エラー

CLIPのViT-L/14は大きなモデルです。**最低8GB以上のメモリ**を推奨します。
Docker Desktopを使用している場合は、設定からメモリ割り当てを増やしてください。

### ポート競合

8000番ポートが使用中の場合は別のポートにマッピングします：

```bash
docker run -p 9000:8000 room-analyzer
# → http://localhost:9000/ でアクセス
```

---

## モデルのキャッシュについて

コンテナを削除するたびにCLIPとYOLOのモデルが再ダウンロードされます。
ダウンロードを省略したい場合は、Dockerボリュームを使ってキャッシュを永続化できます：

```bash
docker run -p 8000:8000 \
  -v room-analyzer-cache:/root/.cache \
  room-analyzer
```
