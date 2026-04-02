# 片付けチェッカー

## 開発背景
部屋の片づけを定期的に行えずに、徐々に部屋が散らかりだすことに困っており、それを解決できるアプリを作りたいと考え作りました。

## 概要
部屋の写真をアップロードすると、AIが部屋の散らかり具合を判定し、スコアを返すアプリです。

## 機能
- 部屋の写真をアップロードすると、AIが部屋の散らかり具合を判定し、スコアを返す
<img width="796" height="646" alt="image" src="https://github.com/user-attachments/assets/d911a5ad-c99f-4c28-aae1-804e1d981d9c" />

<img width="804" height="227" alt="image" src="https://github.com/user-attachments/assets/dc9787fd-6369-4b3e-81b4-6add23f46e2d" />

- これまでのスコアをグラフで確認できる
<img width="796" height="995" alt="image" src="https://github.com/user-attachments/assets/03a1f3e0-f48f-4344-a44d-dcc449bd2fb8" />

- 目標スコアを設定することで、今日片付けを行うべきか確認できる
<img width="545" height="317" alt="image" src="https://github.com/user-attachments/assets/f64944ab-05c7-4d4c-86dc-d9d7ef411a13" />


## 使用技術

### フロントエンド
- **フレームワーク**: Next.js (App Router)
- **言語**: TypeScript
- **UIコンポーネント**: Shadcn UI / Radix UI
- **スタイリング**: Tailwind CSS
- **チャート**: Recharts
- **状態管理/バリデーション**: React Hook Form / Zod

### バックエンド
- **言語**: Go
- **フレームワーク**: Gin
- **ORM**: GORM
- **データベース**: PostgreSQL
- **APIドキュメント**: Swagger (swaggo/gin-swagger)
- **ホットリロード**: Air

### AI解析 (Python)
- **フレームワーク**: FastAPI
- **推論エンジン**: PyTorch
- **モデル**: 
  - OpenAI CLIP (散らかり具合の判定)
  - Ultralytics YOLOv11 (物体検出: ボトル、カップ、缶等)
- **画像処理**: Pillow

### インフラ・その他
- **コンテナ化**: Docker / Docker Compose
- **パッケージ管理**: npm (Frontend), Go Modules (Backend), pip (AI)

### 開発補助ツール
- v0
- ChatGPT
- Claude

## 工夫した点

### AIモデル選定で工夫した点

本アプリでは、3秒以内でスコアを返すことと、散らかり具合を納得感のある形で数値化することを重視してAIモデルを選定しました。

モデル選定では、単純な精度だけでなく、以下の観点を重視しています。

- 処理速度
- コスト
- 説明可能性
- 散らかり度合いの評価という目的に対する適合性

最終的に、  
- CLIP を用いて部屋全体の散らかり感をスコア化し、  
- YOLO11l を用いてペットボトルなどの具体物を補助的に検出する構成を採用しました。

また、推論速度改善のために画像サイズやモデル構成を調整し、  
約5秒 → 1秒未満まで処理時間を短縮しています。

### 生成AIの活用
生成AIを補助的に活用し、開発スピードと検討の幅を高めました。

- ChatGPT：実装方法の整理、AIモデル選定の壁打ち
- Claude：実装補助、コード整理、リファクタリング
- v0：UI改善、画面構成の見直し

最終的な要件整理・技術選定・実装方針の判断は自身で行っています。
