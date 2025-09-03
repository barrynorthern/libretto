import { useEffect, useMemo, useState } from 'react'
import './style.css'
import { ListScenes } from '../wailsjs/go/main/App'
import { ProjectFoundationStep, TemplateSelectionStep, ToneSelectionStep } from './WizardSteps'
import { Sidebar } from './components/Sidebar'
import { useNavigation } from './hooks/useNavigation'
import { useWizardPersistence, type ProjectData } from './hooks/useWizardPersistence'
import { ConfirmDialog } from './components/ConfirmDialog'
import { AIOrchestra } from './components/AIOrchestra'

interface SceneDTO { id: string; title: string; summary: string; content: string; created: string }

type AnyWindow = Window & { go?: any }

function App() {
  const [scenes, setScenes] = useState<SceneDTO[]>([])
  const [selectedScene, setSelectedScene] = useState<SceneDTO | null>(null)
  const [currentProject, setCurrentProject] = useState({
    title: 'Serious Games',
    description: 'Coming-of-Age Tech Thriller'
  })

  const navigation = useNavigation()
  const wizard = useWizardPersistence()

  const hasWails = useMemo(() => {
    if (typeof window === 'undefined') return false
    const w = window as AnyWindow
    return !!w.go?.main?.App?.ListScenes
  }, [])

  useEffect(() => {
    if (!hasWails) return
    ListScenes().then(setScenes).catch(() => setScenes([]))
  }, [hasWails])

  // Handle project wizard completion
  const handleProjectComplete = (projectData: ProjectData) => {
    setCurrentProject(projectData)
    setScenes([]) // Reset scenes for new project
    navigation.exitWizardFlow()
  }

  return (
    <div id="App" className="flex h-screen bg-slate-900 text-white overflow-hidden">
      <Sidebar
        currentScreen={navigation.currentScreen}
        onNavigate={navigation.navigateTo}
        canNavigate={navigation.canNavigate}
      >
        {/* Conductor */}
        <div className="mt-auto p-4 border-t border-slate-700">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-slate-600 rounded-full flex items-center justify-center">
              <span className="text-xs">C</span>
            </div>
            <div>
              <div className="text-sm font-medium">Conductor</div>
              <div className="text-xs text-slate-400">Master of the Orchestra</div>
            </div>
          </div>
        </div>
      </Sidebar>

      {/* Main Content */}
      {navigation.currentScreen === 'canvas' && (
        <div className="flex-1 flex flex-col min-w-0">
          {/* Top Bar */}
          <div className="h-16 bg-slate-800 border-b border-slate-700 flex items-center justify-between px-6">
            <div>
              <h1 className="text-lg font-semibold text-slate-200">{currentProject.title}</h1>
              <div className="text-xs text-slate-500">
                {currentProject.description} • {scenes.length} scenes • 0 characters
              </div>
            </div>
            <div className="flex gap-2">
              <button className="px-3 py-1.5 text-xs bg-slate-700/50 text-slate-400 rounded-md hover:bg-slate-600/70 hover:text-slate-300 transition-colors">
                Overview
              </button>
              <button className="px-3 py-1.5 text-xs bg-slate-700/50 text-slate-400 rounded-md hover:bg-slate-600/70 hover:text-slate-300 transition-colors">
                Notes
              </button>
              <button className="px-3 py-1.5 text-xs bg-purple-600 text-white rounded-md hover:bg-purple-500 transition-colors">
                Generate
              </button>
            </div>
          </div>

          <div className="flex-1 flex">
            {/* Canvas Area */}
            <div className="flex-1 flex flex-col">
              <div className="flex-1 flex items-center justify-center bg-slate-900">
                {scenes.length === 0 ? (
                  <EmptyCanvas />
                ) : (
                  <ScenesCanvas scenes={scenes} onSelectScene={setSelectedScene} />
                )}
              </div>

              {/* Baton Input */}
              <BatonInput onSceneCreated={(s) => setScenes(prev => [s, ...prev])} />
            </div>

            {/* Right Inspector */}
            <Inspector selectedScene={selectedScene} />
          </div>
        </div>
      )}

      {/* New Project Wizard */}
      {navigation.currentScreen === 'new-project' && (
        <NewProjectWizard
          projectData={wizard.projectData}
          step={wizard.step}
          onUpdateData={wizard.updateData}
          onGoToStep={wizard.goToStep}
          onCancel={navigation.exitWizardFlow}
          onComplete={handleProjectComplete}
          onClearData={wizard.clearData}
          hasData={wizard.hasData()}
        />
      )}

      {/* AI Orchestra Screen */}
      {navigation.currentScreen === 'ai-orchestra' && (
        <AIOrchestra />
      )}
    </div>
  )
}



function EmptyCanvas() {
  return (
    <div className="text-center">
      <div className="w-12 h-12 bg-slate-700/30 rounded-md flex items-center justify-center mx-auto mb-4 border border-slate-600/30">
        <span className="text-sm text-slate-500">◐</span>
      </div>
      <h3 className="text-sm font-semibold mb-2 text-slate-300">Empty Canvas</h3>
      <p className="text-xs text-slate-500 mb-6">Use the Baton below to begin crafting your narrative</p>
      <div className="flex flex-wrap gap-2 justify-center max-w-md mx-auto">
        <SuggestionChip text="Create opening scene" />
        <SuggestionChip text="Introduce protagonist" />
        <SuggestionChip text="Set the world" />
      </div>
    </div>
  )
}

function SuggestionChip({ text }: { text: string }) {
  return (
    <button className="px-2 py-1 bg-slate-700/30 text-slate-500 rounded text-xs hover:bg-slate-600/50 hover:text-slate-400 whitespace-nowrap overflow-hidden text-ellipsis max-w-[180px] transition-colors border border-slate-600/30">
      {text}
    </button>
  )
}

function ScenesCanvas({ scenes, onSelectScene }: { scenes: SceneDTO[]; onSelectScene: (scene: SceneDTO) => void }) {
  return (
    <div className="w-full max-w-4xl p-6">
      <div className="grid gap-4">
        {scenes.map((scene) => (
          <div
            key={scene.id}
            onClick={() => onSelectScene(scene)}
            className="p-4 bg-slate-800 border border-slate-700 rounded-lg cursor-pointer hover:border-purple-500 transition-colors"
          >
            <h3 className="font-semibold mb-2">{scene.title}</h3>
            <p className="text-slate-400 text-sm">{scene.summary}</p>
            <div className="text-xs text-slate-500 mt-2">
              Created: {new Date(scene.created).toLocaleString()}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

function BatonInput({ onSceneCreated }: { onSceneCreated: (scene: SceneDTO) => void }) {
  const [directive, setDirective] = useState('')
  const [busy, setBusy] = useState(false)

  const handleSubmit = async () => {
    if (!directive.trim()) return
    try {
      setBusy(true)
      const mod = await import('../wailsjs/go/main/App')
      const created = await mod.CreateScene(
        `Scene from: ${directive.slice(0, 30)}...`,
        directive,
        `Generated from directive: "${directive}"`
      )
      onSceneCreated(created as unknown as SceneDTO)
      setDirective('')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="bg-slate-800 border-t border-slate-700 p-4">
      <div className="flex items-center gap-3 mb-3">
        <span className="text-xs text-slate-500">◐ Direct your AI Orchestra...</span>
      </div>
      <div className="flex gap-2">
        <input
          value={directive}
          onChange={(e) => setDirective(e.target.value)}
          placeholder="e.g., 'Create a dramatic confrontation scene'"
          className="flex-1 px-3 py-2 text-xs bg-slate-800/50 border border-slate-600/50 rounded-md text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500/25 transition-all duration-200"
          onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
        />
        <button
          onClick={handleSubmit}
          disabled={busy || !directive.trim()}
          className="px-4 py-2 text-xs bg-purple-600 text-white rounded-md hover:bg-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
        >
          {busy ? 'Processing...' : 'Execute'}
        </button>
      </div>
      <div className="flex flex-wrap gap-2 mt-3 items-center">
        <span className="text-xs text-slate-500 whitespace-nowrap">Suggested Directives:</span>
        <div className="flex flex-wrap gap-2">
          <SuggestionChip text="Create an opening scene where the protag..." />
          <SuggestionChip text="Introduce a mentor character with myster..." />
          <SuggestionChip text="Add tension between two main characters" />
        </div>
      </div>
    </div>
  )
}

function Inspector({ selectedScene }: { selectedScene: SceneDTO | null }) {
  return (
    <div className="w-80 min-w-[320px] bg-slate-800 border-l border-slate-700 flex flex-col">
      <div className="p-4 border-b border-slate-700">
        <h3 className="text-sm font-semibold text-slate-200">Inspector</h3>
        <p className="text-xs text-slate-500">
          {selectedScene ? 'Scene details' : 'Select an element on the canvas to view details'}
        </p>
      </div>
      <div className="flex-1 p-4">
        {selectedScene ? (
          <div className="space-y-4">
            <div>
              <label className="text-xs font-semibold text-slate-500 uppercase tracking-wide">Title</label>
              <div className="mt-1 text-xs text-slate-300">{selectedScene.title}</div>
            </div>
            <div>
              <label className="text-xs font-semibold text-slate-500 uppercase tracking-wide">Summary</label>
              <div className="mt-1 text-xs text-slate-400">{selectedScene.summary}</div>
            </div>
            <div>
              <label className="text-xs font-semibold text-slate-500 uppercase tracking-wide">Content</label>
              <div className="mt-1 text-xs text-slate-400">{selectedScene.content}</div>
            </div>
            <div>
              <label className="text-xs font-semibold text-slate-500 uppercase tracking-wide">Created</label>
              <div className="mt-1 text-xs text-slate-400">
                {new Date(selectedScene.created).toLocaleString()}
              </div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <div className="w-10 h-10 bg-slate-700/30 rounded-md flex items-center justify-center mb-3 border border-slate-600/30">
              <span className="text-sm text-slate-500">◊</span>
            </div>
            <h4 className="text-xs font-medium mb-2 text-slate-400">Waiting for Selection</h4>
            <p className="text-xs text-slate-500">
              Click on any scene or character to inspect its details
            </p>
          </div>
        )}
      </div>
    </div>
  )
}

// New Project Wizard Component
function NewProjectWizard({
  projectData,
  step,
  onUpdateData,
  onGoToStep,
  onCancel,
  onComplete,
  onClearData,
  hasData
}: {
  projectData: Partial<ProjectData>;
  step: number;
  onUpdateData: (updates: Partial<ProjectData>) => void;
  onGoToStep: (step: number) => void;
  onCancel: () => void;
  onComplete: (data: ProjectData) => void;
  onClearData: () => void;
  hasData: boolean;
}) {
  const [showClearConfirm, setShowClearConfirm] = useState(false)

  const handleComplete = () => {
    onComplete(projectData as ProjectData)
  }

  const handleNewProject = () => {
    if (hasData) {
      setShowClearConfirm(true)
    } else {
      onClearData()
    }
  }

  const handleConfirmClear = () => {
    onClearData()
    setShowClearConfirm(false)
  }

  return (
    <div className="fixed inset-0 bg-slate-900 z-50">
      <div className="h-screen flex">
        <Sidebar
          currentScreen="new-project"
          onNavigate={() => {}} // Disabled during wizard
          canNavigate={false}
        >
          {/* Conductor */}
          <div className="mt-auto p-4 border-t border-slate-700">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-slate-600 rounded-full flex items-center justify-center">
                <span className="text-xs">C</span>
              </div>
              <div>
                <div className="text-sm font-medium">Conductor</div>
                <div className="text-xs text-slate-400">Master of the Orchestra</div>
              </div>
            </div>
          </div>
        </Sidebar>

        {/* Main Content */}
        <div className="flex-1 flex flex-col min-w-0">
          {/* Header */}
          <div className="p-6 border-b border-slate-700">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-slate-700 rounded-lg flex items-center justify-center">
                  <span className="text-lg">📄</span>
                </div>
                <div>
                  <h2 className="text-xl font-semibold">Bootstrap New Narrative</h2>
                  <p className="text-slate-400">Configure your AI Orchestra to begin the creative process</p>
                </div>
              </div>
              <div className="flex gap-2">
                {step === 1 && (
                  <button
                    onClick={handleNewProject}
                    className="px-3 py-1.5 text-xs bg-red-600/80 text-white rounded-md hover:bg-red-500 transition-colors"
                  >
                    New
                  </button>
                )}
                <button
                  onClick={onCancel}
                  className="px-3 py-1.5 text-xs bg-slate-700/50 text-slate-400 rounded-md hover:bg-slate-600/70 hover:text-slate-300 transition-colors"
                >
                  Cancel
                </button>
              </div>
            </div>

            {/* Progress Steps */}
            <div className="flex items-center gap-4">
              <StepIndicator number={1} active={step === 1} completed={step > 1} />
              <div className="flex-1 h-px bg-slate-600"></div>
              <StepIndicator number={2} active={step === 2} completed={step > 2} />
              <div className="flex-1 h-px bg-slate-600"></div>
              <StepIndicator number={3} active={step === 3} completed={step > 3} />
            </div>
          </div>

          {/* Step Content */}
          <div className="flex-1 p-6 flex justify-center">
            <div className="w-full max-w-4xl">
                {step === 1 && (
                  <ProjectFoundationStep
                    data={projectData}
                    onUpdate={onUpdateData}
                    onNext={() => onGoToStep(2)}
                    onCancel={onCancel}
                  />
                )}
                {step === 2 && (
                  <TemplateSelectionStep
                    data={projectData}
                    onUpdate={onUpdateData}
                    onNext={() => onGoToStep(3)}
                    onPrevious={() => onGoToStep(1)}
                  />
                )}
                {step === 3 && (
                  <ToneSelectionStep
                    data={projectData}
                    onUpdate={onUpdateData}
                    onComplete={handleComplete}
                    onPrevious={() => onGoToStep(2)}
                  />
                )}
            </div>
          </div>
        </div>
      </div>

      <ConfirmDialog
        isOpen={showClearConfirm}
        title="Clear Project Setup?"
        message={
          <div>
            <p className="mb-2">This will permanently clear all project setup data including:</p>
            <ul className="list-disc list-inside text-xs space-y-1 text-slate-500">
              <li>Project title and description</li>
              <li>Genre and core theme</li>
              <li>Selected template</li>
              <li>Tone and notes</li>
            </ul>
            <p className="mt-3 text-xs">This action cannot be undone.</p>
          </div>
        }
        confirmText="Clear All"
        cancelText="Keep Data"
        onConfirm={handleConfirmClear}
        onCancel={() => setShowClearConfirm(false)}
        destructive
      />
    </div>
  )
}

function StepIndicator({ number, active, completed }: {
  number: number;
  active: boolean;
  completed: boolean;
}) {
  return (
    <div className={`w-7 h-7 rounded-full flex items-center justify-center text-xs font-semibold transition-colors ${
      completed ? 'bg-purple-600 text-white' :
      active ? 'bg-purple-600 text-white' :
      'bg-slate-700/50 text-slate-500 border border-slate-600/50'
    }`}>
      {completed ? '✓' : number}
    </div>
  )
}

export default App
