import { useState, useRef } from 'react'
import './App.css'

function App() {
  const [score, setScore] = useState<number | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [fileName, setFileName] = useState<string | null>(null)
  const [preview, setPreview] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setFileName(file.name)
      setScore(null)
      setError(null)
      const reader = new FileReader()
      reader.onload = (ev) => setPreview(ev.target?.result as string)
      reader.readAsDataURL(file)
    }
  }

  const handleSubmit = async () => {
    const file = fileInputRef.current?.files?.[0]
    if (!file) {
      setError('画像ファイルを選択してください')
      return
    }

    setLoading(true)
    setError(null)
    setScore(null)

    try {
      const formData = new FormData()
      formData.append('image', file)

      const res = await fetch('/analyze', {
        method: 'POST',
        body: formData,
      })

      if (!res.ok) {
        throw new Error(`HTTP Error: ${res.status}`)
      }

      const data = await res.json()
      setScore(data.score)
    } catch (err) {
      setError(err instanceof Error ? err.message : '通信エラーが発生しました')
    } finally {
      setLoading(false)
    }
  }

  const getScoreColor = (s: number) => {
    if (s >= 70) return '#22c55e'
    if (s >= 40) return '#f59e0b'
    return '#ef4444'
  }

  return (
    <div className="app">
      <div className="card">
        <h1 className="title">
          <span className="icon">🔍</span>
          Image Analyzer
        </h1>
        <p className="subtitle">画像をアップロードしてスコアを確認</p>

        <div className="upload-area" onClick={() => fileInputRef.current?.click()}>
          {preview ? (
            <img src={preview} alt="preview" className="preview-image" />
          ) : (
            <div className="upload-placeholder">
              <span className="upload-icon">📁</span>
              <p>クリックして画像を選択</p>
            </div>
          )}
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            className="file-input"
          />
        </div>

        {fileName && <p className="file-name">📎 {fileName}</p>}

        <button
          onClick={handleSubmit}
          disabled={loading || !fileName}
          className="submit-btn"
        >
          {loading ? (
            <span className="spinner" />
          ) : (
            '分析する'
          )}
        </button>

        {error && <p className="error">⚠️ {error}</p>}

        {score !== null && (
          <div className="result">
            <p className="result-label">分析スコア</p>
            <div
              className="score"
              style={{ color: getScoreColor(score) }}
            >
              {score}
            </div>
            <div className="score-bar-bg">
              <div
                className="score-bar"
                style={{
                  width: `${score}%`,
                  backgroundColor: getScoreColor(score),
                }}
              />
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
