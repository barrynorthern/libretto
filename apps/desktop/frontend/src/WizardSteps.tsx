import { useState } from 'react'

interface ProjectData {
  title: string;
  description: string;
  genre: string;
  coreTheme: string;
  template: string;
  tone: string;
  notes?: string;
}

// Step 1: Project Foundation
export function ProjectFoundationStep({
  data,
  onUpdate,
  onNext,
  onCancel
}: {
  data: Partial<ProjectData>;
  onUpdate: (updates: Partial<ProjectData>) => void;
  onNext: () => void;
  onCancel?: () => void;
}) {
  const canProceed = data.title && data.description && data.genre && data.coreTheme

  return (
    <div className="w-full">
      <div className="flex items-center gap-3 mb-6">
        <span className="text-sm text-slate-400">◐</span>
        <h3 className="text-sm font-semibold text-slate-200">Project Foundation</h3>
      </div>

      <div className="space-y-6">
        <div>
          <label className="block text-xs font-medium mb-2 text-slate-300 uppercase tracking-wide">
            Project Title <span className="text-red-400">*</span>
          </label>
          <input
            type="text"
            value={data.title || ''}
            onChange={(e) => onUpdate({ title: e.target.value })}
            placeholder="The Chronicles of..."
            className="w-full px-4 py-3 bg-slate-800/50 border border-slate-600/50 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500/25 transition-all duration-200"
          />
        </div>

        <div>
          <label className="block text-xs font-medium mb-2 text-slate-300 uppercase tracking-wide">Description</label>
          <textarea
            value={data.description || ''}
            onChange={(e) => onUpdate({ description: e.target.value })}
            placeholder="Brief description of your story..."
            rows={4}
            className="w-full px-4 py-3 bg-slate-800/50 border border-slate-600/50 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500/25 transition-all duration-200 resize-none"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-xs font-medium mb-2 text-slate-300 uppercase tracking-wide">Genre</label>
            <input
              type="text"
              value={data.genre || ''}
              onChange={(e) => onUpdate({ genre: e.target.value })}
              placeholder="Fantasy, Sci-Fi, Romance..."
              className="w-full px-4 py-3 bg-slate-800/50 border border-slate-600/50 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500/25 transition-all duration-200"
            />
          </div>
          <div>
            <label className="block text-xs font-medium mb-2 text-slate-300 uppercase tracking-wide">Core Theme</label>
            <input
              type="text"
              value={data.coreTheme || ''}
              onChange={(e) => onUpdate({ coreTheme: e.target.value })}
              placeholder="What question does your story ask?"
              className="w-full px-4 py-3 bg-slate-800/50 border border-slate-600/50 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500/25 transition-all duration-200"
            />
          </div>
        </div>
      </div>

      <div className="flex justify-between mt-8">
        {onCancel && (
          <button
            onClick={onCancel}
            className="px-4 py-2 text-xs bg-slate-700/50 text-slate-400 rounded-md hover:bg-slate-600/70 hover:text-slate-300 transition-colors font-medium"
          >
            Cancel
          </button>
        )}
        <button
          onClick={onNext}
          disabled={!canProceed}
          className="px-4 py-2 text-xs bg-purple-600 text-white rounded-md hover:bg-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium ml-auto"
        >
          Next
        </button>
      </div>
    </div>
  )
}

// Step 2: Template Selection
export function TemplateSelectionStep({ 
  data, 
  onUpdate, 
  onNext, 
  onPrevious 
}: {
  data: Partial<ProjectData>;
  onUpdate: (updates: Partial<ProjectData>) => void;
  onNext: () => void;
  onPrevious: () => void;
}) {
  const templates = [
    {
      id: 'three-act',
      name: 'Three-Act Structure',
      description: 'Classic beginning, middle, and structure',
      icon: '◐',
      scenes: '~12 scenes'
    },
    {
      id: 'heros-journey',
      name: "Hero's Journey",
      description: "Campbell's monomyth template",
      icon: '◊',
      scenes: '~17 scenes'
    },
    {
      id: 'romance-arc',
      name: 'Romance Arc',
      description: 'Meet-cute to happily ever after',
      icon: '◈',
      scenes: '~10 scenes'
    },
    {
      id: 'mystery',
      name: 'Mystery Structure',
      description: 'Crime, clues, and revelation',
      icon: '◉',
      scenes: '~15 scenes'
    },
    {
      id: 'thriller',
      name: 'Thriller Pacing',
      description: 'High-tension escalating stakes',
      icon: '◆',
      scenes: '~14 scenes'
    }
  ]

  const canProceed = !!data.template

  return (
    <div className="w-full">
      <div className="flex items-center gap-3 mb-6">
        <span className="text-sm text-slate-400">◊</span>
        <h3 className="text-sm font-semibold text-slate-200">Choose Your Template</h3>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
        {templates.map((template) => (
          <TemplateCard
            key={template.id}
            template={template}
            selected={data.template === template.id}
            onSelect={() => onUpdate({ template: template.id })}
          />
        ))}
      </div>

      <div className="flex justify-between">
        <button
          onClick={onPrevious}
          className="px-4 py-2 text-xs bg-slate-700/50 text-slate-400 rounded-md hover:bg-slate-600/70 hover:text-slate-300 transition-colors font-medium"
        >
          Previous
        </button>
        <button
          onClick={onNext}
          disabled={!canProceed}
          className="px-4 py-2 text-xs bg-purple-600 text-white rounded-md hover:bg-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
        >
          Continue
        </button>
      </div>
    </div>
  )
}

function TemplateCard({ template, selected, onSelect }: {
  template: any;
  selected: boolean;
  onSelect: () => void;
}) {
  return (
    <div
      onClick={onSelect}
      className={`p-4 rounded-lg border cursor-pointer transition-colors ${
        selected
          ? 'border-purple-500/60 bg-purple-500/10'
          : 'border-slate-600/50 bg-slate-800/30 hover:border-slate-500/70 hover:bg-slate-700/40'
      }`}
    >
      <div className="text-sm mb-3 text-slate-400">{template.icon}</div>
      <h4 className="font-semibold mb-1 text-xs text-slate-200">{template.name}</h4>
      <p className="text-xs text-slate-500 mb-2 leading-relaxed">{template.description}</p>
      <div className="text-xs text-slate-600 font-medium">{template.scenes}</div>
    </div>
  )
}

// Step 3: Tone Selection
export function ToneSelectionStep({ 
  data, 
  onUpdate, 
  onComplete, 
  onPrevious 
}: {
  data: Partial<ProjectData>;
  onUpdate: (updates: Partial<ProjectData>) => void;
  onComplete: () => void;
  onPrevious: () => void;
}) {
  const tones = [
    {
      id: 'light',
      name: 'Light & Optimistic',
      icon: '☀️',
      description: 'Uplifting, hopeful, and positive. Characters overcome challenges with resilience and find joy.',
      color: 'bg-yellow-500'
    },
    {
      id: 'dark',
      name: 'Dark & Serious',
      icon: '🌙',
      description: 'Intense, gritty, and realistic. Explores difficult themes with emotional depth.',
      color: 'bg-slate-500'
    },
    {
      id: 'mysterious',
      name: 'Mysterious & Suspenseful',
      icon: '🔮',
      description: 'Intriguing, enigmatic, and full of secrets. Keeps readers guessing until the end.',
      color: 'bg-purple-500'
    },
    {
      id: 'comedic',
      name: 'Comedic & Playful',
      icon: '🎭',
      description: 'Humorous, witty, and entertaining. Uses comedy to explore characters and situations.',
      color: 'bg-green-500'
    },
    {
      id: 'dramatic',
      name: 'Dramatic & Emotional',
      icon: '🎪',
      description: 'Intense emotions, high stakes, and powerful character moments that resonate deeply.',
      color: 'bg-red-500'
    }
  ]

  const canProceed = !!data.tone

  return (
    <div className="w-full">
      <div className="flex items-center gap-3 mb-6">
        <span className="text-sm text-slate-400">◈</span>
        <div>
          <h3 className="text-sm font-semibold text-slate-200">Set the Tone</h3>
          <p className="text-xs text-slate-500 mt-1">Choose the overall mood and atmosphere for your narrative. This can be adjusted later.</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
        {tones.map((tone) => (
          <ToneCard
            key={tone.id}
            tone={tone}
            selected={data.tone === tone.id}
            onSelect={() => onUpdate({ tone: tone.id })}
          />
        ))}
      </div>

      {/* Additional Notes */}
      <div className="mb-8">
        <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wide mb-3">
          Additional Notes (Optional)
        </label>
        <textarea
          value={data.notes || ''}
          onChange={(e) => onUpdate({ notes: e.target.value })}
          placeholder="Any specific tone considerations, mood preferences, or stylistic notes for your narrative..."
          className="w-full px-4 py-3 bg-slate-800/50 border border-slate-600/50 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 focus:ring-1 focus:ring-purple-500/25 transition-all duration-200 resize-none"
          rows={3}
        />
        <p className="text-xs text-slate-600 mt-2">These notes will help guide the AI Orchestra in maintaining consistent tone throughout your narrative.</p>
      </div>

      <div className="flex justify-between">
        <button
          onClick={onPrevious}
          className="px-4 py-2 text-xs bg-slate-700/50 text-slate-400 rounded-md hover:bg-slate-600/70 hover:text-slate-300 transition-colors font-medium"
        >
          Previous
        </button>
        <button
          onClick={onComplete}
          disabled={!canProceed}
          className="px-4 py-2 text-xs bg-purple-600 text-white rounded-md hover:bg-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium flex items-center gap-2"
        >
          <span className="text-xs">◈</span>
          Begin Orchestra
        </button>
      </div>
    </div>
  )
}

function ToneCard({ tone, selected, onSelect }: {
  tone: any;
  selected: boolean;
  onSelect: () => void;
}) {
  return (
    <div
      onClick={onSelect}
      className={`p-4 rounded-lg border cursor-pointer transition-colors ${
        selected
          ? 'border-purple-500/60 bg-purple-500/10'
          : 'border-slate-600/50 bg-slate-800/30 hover:border-slate-500/70 hover:bg-slate-700/40'
      }`}
    >
      <div className="flex items-center gap-3 mb-3">
        <span className="text-2xl">{tone.icon}</span>
        <div className={`w-3 h-3 ${tone.color} rounded-full`}></div>
      </div>
      <h4 className="font-semibold text-sm text-slate-200 mb-2">{tone.name}</h4>
      <p className="text-xs text-slate-400 leading-relaxed">{tone.description}</p>
    </div>
  )
}
