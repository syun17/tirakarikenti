Docker開発環境構築 (React + Go)
画像を送信してダミーのスコアを返すアプリをDockerで開発できるようにする。

Proposed Changes
Backend (Go/Gin)
[NEW] 
go.mod
Go moduleの定義。gin-gonic/ginとgin-contrib/corsを依存に含む。

[NEW] 
main.go
POST /analyze: multipart/form-dataからimageを受け取り、0〜100のランダムスコアをJSONで返す
CORS設定: すべてのオリジンを許可
[NEW] 
Dockerfile
Go 1.22ベースのマルチステージビルド。開発用にホットリロード (air) を使用。

Frontend (React/Vite/TypeScript)
[NEW] Viteプロジェクト一式
npx create-vite でReact + TypeScriptテンプレートから初期化。

[NEW] 
App.tsx
ファイル選択UI
POST /analyzeに画像をmultipart/form-dataで送信
レスポンスのscoreを表示
[NEW] 
Dockerfile
Node 20ベース。Vite dev serverを起動。

Docker Compose
[NEW] 
docker-compose.yml
backend: ポート8080
frontend: ポート5173、backendに依存
frontendからbackendへのプロキシ設定（Vite config）
Verification Plan
Automated Tests
bash
docker-compose up --build
ビルドが成功しエラーなく起動することを確認。

Manual Verification
ブラウザで http://localhost:5173 を開く
画像ファイルを選択して送信ボタンを押す
0〜100のスコアが画面に表示されることを確認