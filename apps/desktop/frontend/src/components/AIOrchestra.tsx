import { useState } from 'react'

interface Agent {
  id: string
  name: string
  icon: string
  description: string
  status: 'active' | 'inactive'
  eventsProcessed: number
  successRate: number
  avgResponse: number
  subscriptions: number
  eventSubscriptions: string[]
}

const mockAgents: Agent[] = [
  {
    id: 'world-architect',
    name: 'World Architect',
    icon: '🏛️',
    description: 'Creates and maintains consistent world-building and setting details',
    status: 'inactive',
    eventsProcessed: 52,
    successRate: 85,
    avgResponse: 5.1,
    subscriptions: 3,
    eventSubscriptions: ['WorldBuilding', 'LocationCreation', 'LoreGeneration']
  },
  {
    id: 'the-empath',
    name: 'The Empath',
    icon: '🎭',
    description: 'Steward of emotional core and character arc development',
    status: 'active',
    eventsProcessed: 22,
    successRate: 87,
    avgResponse: 5.0,
    subscriptions: 3,
    eventSubscriptions: ['SceneGenerated', 'EmotionalAnalysis', 'ArcDevelopment']
  },
  {
    id: 'continuity-steward',
    name: 'Continuity Steward',
    icon: '🔗',
    description: 'Guardian of narrative consistency, monitoring for contradictions',
    status: 'active',
    eventsProcessed: 41,
    successRate: 90,
    avgResponse: 4.7,
    subscriptions: 3,
    eventSubscriptions: ['NarrativeRestructured', 'SceneAdded', 'CharacterUpdated']
  },
  {
    id: 'thematic-steward',
    name: 'Thematic Steward',
    icon: '🎯',
    description: 'Guardian of the story\'s central thematic question',
    status: 'active',
    eventsProcessed: 52,
    successRate: 94,
    avgResponse: 5.3,
    subscriptions: 3,
    eventSubscriptions: ['ThematicAnalysis', 'ContentReview', 'FocusCheck']
  },
  {
    id: 'plot-weaver',
    name: 'Plot Weaver',
    icon: '🕸️',
    description: 'Generates and refines narrative structure, causality, and plot points',
    status: 'active',
    eventsProcessed: 48,
    successRate: 86,
    avgResponse: 2.7,
    subscriptions: 3,
    eventSubscriptions: ['DirectiveIssued', 'PlotDevelopment', 'StructuralChange']
  },
  {
    id: 'the-dramaturg',
    name: 'The Dramaturg',
    icon: '✂️',
    description: 'Ruthless editor focused on narrative focus and impact',
    status: 'active',
    eventsProcessed: 22,
    successRate: 90,
    avgResponse: 4.7,
    subscriptions: 3,
    eventSubscriptions: ['ActCompleted', 'DraftFinished', 'EditingPass']
  },
  {
    id: 'character-troupe',
    name: 'Character Troupe',
    icon: '👥',
    description: 'Method actors for each character, ensuring authentic voice and motivation',
    status: 'active',
    eventsProcessed: 45,
    successRate: 88,
    avgResponse: 4.7,
    subscriptions: 3,
    eventSubscriptions: ['CharacterDevelopment', 'DialogueGeneration', 'CharacterInteraction']
  }
]

function StatusBadge({ status }: { status: 'active' | 'inactive' }) {
  return (
    <div className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${
      status === 'active' 
        ? 'bg-green-500/20 text-green-400' 
        : 'bg-slate-500/20 text-slate-400'
    }`}>
      <div className={`w-1.5 h-1.5 rounded-full ${
        status === 'active' ? 'bg-green-400' : 'bg-slate-400'
      }`} />
      {status === 'active' ? 'Active' : 'Inactive'}
    </div>
  )
}

function MetricCard({ label, value, className = '' }: { label: string; value: string | number; className?: string }) {
  return (
    <div className={`bg-slate-800/50 border border-slate-700/50 rounded-lg p-4 ${className}`}>
      <div className="text-xs text-slate-500 mb-1">{label}</div>
      <div className="text-2xl font-semibold text-slate-200">{value}</div>
    </div>
  )
}

function AgentCard({ agent }: { agent: Agent }) {
  return (
    <div className="bg-slate-800/50 border border-slate-700/50 rounded-lg p-4">
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-3">
          <div className="text-2xl">{agent.icon}</div>
          <div>
            <h3 className="font-semibold text-sm text-slate-200">{agent.name}</h3>
            <StatusBadge status={agent.status} />
          </div>
        </div>
      </div>
      
      <p className="text-xs text-slate-400 mb-4 leading-relaxed">{agent.description}</p>
      
      <div className="grid grid-cols-2 gap-3 mb-4">
        <div>
          <div className="text-xs text-slate-500">Events Processed</div>
          <div className="text-lg font-semibold text-slate-200">{agent.eventsProcessed}</div>
        </div>
        <div>
          <div className="text-xs text-slate-500">Success Rate</div>
          <div className="text-lg font-semibold text-slate-200">{agent.successRate}%</div>
        </div>
        <div>
          <div className="text-xs text-slate-500">Avg Response</div>
          <div className="text-lg font-semibold text-slate-200">{agent.avgResponse}s</div>
        </div>
        <div>
          <div className="text-xs text-slate-500">Subscriptions</div>
          <div className="text-lg font-semibold text-slate-200">{agent.subscriptions}</div>
        </div>
      </div>
      
      <div>
        <div className="text-xs text-slate-500 mb-2">Event Subscriptions:</div>
        <div className="flex flex-wrap gap-1">
          {agent.eventSubscriptions.map((event) => (
            <span
              key={event}
              className="px-2 py-1 bg-slate-700/50 text-slate-400 rounded text-xs"
            >
              {event}
            </span>
          ))}
        </div>
      </div>
    </div>
  )
}

export function AIOrchestra() {
  const activeAgents = mockAgents.filter(agent => agent.status === 'active').length
  const totalAgents = mockAgents.length
  const processingTasks = 3
  const avgSuccessRate = Math.round(mockAgents.reduce((acc, agent) => acc + agent.successRate, 0) / mockAgents.length)

  return (
    <div className="flex-1 flex flex-col min-w-0">
      {/* Header */}
      <div className="h-16 bg-slate-800 border-b border-slate-700 flex items-center justify-between px-6">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-purple-600 rounded-lg flex items-center justify-center">
            <span className="text-white font-bold text-sm">◊</span>
          </div>
          <div>
            <h1 className="text-lg font-semibold text-slate-200">AI Orchestra</h1>
            <div className="text-xs text-slate-500">Manage your specialized narrative agents</div>
          </div>
        </div>
      </div>

      <div className="flex-1 p-6 overflow-auto">
        {/* Metrics Overview */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <MetricCard label="Total Agents" value={totalAgents} />
          <MetricCard label="Active" value={activeAgents} />
          <MetricCard label="Processing" value={`${processingTasks} tasks`} />
          <MetricCard label="Avg Success Rate" value={`${avgSuccessRate}%`} />
        </div>

        {/* Agents Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {mockAgents.map((agent) => (
            <AgentCard key={agent.id} agent={agent} />
          ))}
        </div>
      </div>
    </div>
  )
}
