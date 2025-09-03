import { useEffect, useMemo, useState } from 'react'
import './style.css'
import { ListScenes } from '../wailsjs/go/main/App'

interface SceneDTO { id: string; title: string; summary: string; content: string }

type AnyWindow = Window & { go?: any }

function App() {
  const [scenes, setScenes] = useState<SceneDTO[]>([])
  const hasWails = useMemo(() => {
    if (typeof window === 'undefined') return false
    const w = window as AnyWindow
    return !!w.go?.main?.App?.ListScenes
  }, [])

  useEffect(() => {
    if (!hasWails) return
    ListScenes().then(setScenes).catch(() => setScenes([]))
  }, [hasWails])

  return (
    <div id="App" className="min-h-screen p-6">
      <header className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-zinc-100">Libretto — Scenes</h1>
        <div className="flex gap-2">
          <NewSceneButton onCreated={(s) => setScenes((prev) => [s, ...prev])} />
          <DarkToggle />
        </div>
      </header>
      {!hasWails && (
        <div className="mb-4 rounded border border-amber-700 bg-amber-900/30 p-3 text-amber-300">
          Open the Wails DevServer URL (not the raw Vite port) to access Go bindings.
        </div>
      )}
      {scenes.length === 0 ? (
        <div className="text-zinc-400">No scenes yet.</div>
      ) : (
        <ul className="space-y-3">
          {scenes.map((s) => (
            <li key={s.id} className="rounded border border-zinc-700 bg-zinc-800 p-4">
              <div className="text-lg font-medium">{s.title}</div>
              <div className="text-sm text-zinc-400">{s.summary}</div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

function DarkToggle() {
  const [dark, setDark] = useState(true)
  useEffect(() => {
    const saved = localStorage.getItem('theme')
    const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
    const initial = saved ? saved === 'dark' : prefersDark
    document.documentElement.classList.toggle('dark', initial)
    setDark(initial)
  }, [])
  const toggle = () => {
    const next = !dark
    setDark(next)
    document.documentElement.classList.toggle('dark', next)
    localStorage.setItem('theme', next ? 'dark' : 'light')
  }
  return (
    <button onClick={toggle} className="rounded border border-zinc-700 bg-zinc-800 px-3 py-1 text-zinc-200 hover:bg-zinc-700">
      {dark ? 'Dark' : 'Light'}
    </button>
  )
}

function NewSceneButton({ onCreated }: { onCreated: (s: SceneDTO) => void }) {
  const [title, setTitle] = useState('Untitled Scene')
  const [summary, setSummary] = useState('')
  const [content, setContent] = useState('')
  const [busy, setBusy] = useState(false)
  const click = async () => {
    try {
      setBusy(true)
      // Dynamic import to avoid calling when bindings aren’t present
      const mod = await import('../wailsjs/go/main/App')
      const created = await mod.CreateScene(title, summary, content)
      onCreated(created as unknown as SceneDTO)
      setSummary('')
      setContent('')
    } finally {
      setBusy(false)
    }
  }
  return (
    <div className="flex items-center gap-2">
      <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Title" className="rounded border border-zinc-700 bg-zinc-900 px-2 py-1 text-zinc-100" />
      <button onClick={click} disabled={busy} className="rounded border border-zinc-700 bg-zinc-800 px-3 py-1 text-zinc-100 hover:bg-zinc-700 disabled:opacity-50">
        New Scene
      </button>
    </div>
  )
}

export default App
