import { ReactNode } from 'react'

interface NavItemProps {
  icon: string
  label: string
  active?: boolean
  onClick?: () => void
  disabled?: boolean
}

function NavItem({ icon, label, active = false, onClick, disabled = false }: NavItemProps) {
  const isClickable = !disabled && onClick

  return (
    <div
      onClick={disabled ? undefined : onClick}
      className={`flex items-center gap-3 px-3 py-2 rounded-md transition-all duration-200 ${
        disabled
          ? 'text-slate-600 cursor-not-allowed'
          : isClickable
            ? 'cursor-pointer'
            : 'cursor-default'
      } ${
        active
          ? 'bg-slate-700/60 text-slate-200 border-l-2 border-purple-500/60'
          : disabled
            ? 'text-slate-600'
            : isClickable
              ? 'text-slate-400 hover:bg-slate-700/30 hover:text-slate-300'
              : 'text-slate-400'
      }`}
    >
      <span className="text-xs">{icon}</span>
      <span className="text-xs font-medium">{label}</span>
    </div>
  )
}

interface StatusItemProps {
  label: string
  value: string
}

function StatusItem({ label, value }: StatusItemProps) {
  return (
    <div className="flex justify-between items-center">
      <span className="text-xs text-slate-500">{label}</span>
      <span className="text-xs text-slate-400 font-medium">{value}</span>
    </div>
  )
}

export type NavigationScreen = 'canvas' | 'new-project' | 'ai-orchestra'

interface SidebarProps {
  currentScreen: NavigationScreen
  onNavigate: (screen: NavigationScreen) => void
  canNavigate?: boolean
  children?: ReactNode
}

export function Sidebar({ currentScreen, onNavigate, canNavigate = true, children }: SidebarProps) {
  return (
    <div className="w-64 min-w-[256px] bg-slate-800 border-r border-slate-700 flex flex-col">
      {/* Logo/Header */}
      <div className="p-4 border-b border-slate-700">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-purple-600 rounded-lg flex items-center justify-center">
            <span className="text-white font-bold text-sm">L</span>
          </div>
          <div>
            <div className="font-semibold">Libretto</div>
            <div className="text-xs text-slate-400">Narrative Orchestration Engine</div>
          </div>
        </div>
      </div>

      {/* Navigation */}
      <div className="p-4">
        <div className="text-xs font-semibold text-slate-400 uppercase tracking-wide mb-3">Navigation</div>
        <nav className="space-y-1">
          <NavItem 
            icon="◐" 
            label="Orchestration Canvas" 
            active={currentScreen === 'canvas'}
            onClick={canNavigate ? () => onNavigate('canvas') : undefined}
            disabled={!canNavigate}
          />
          <NavItem 
            icon="+" 
            label="New Project" 
            active={currentScreen === 'new-project'}
            onClick={canNavigate ? () => onNavigate('new-project') : undefined}
            disabled={!canNavigate}
          />
          <NavItem 
            icon="◊" 
            label="AI Orchestra" 
            active={currentScreen === 'ai-orchestra'}
            onClick={canNavigate ? () => onNavigate('ai-orchestra') : undefined}
            disabled={!canNavigate}
          />
        </nav>
      </div>

      {/* AI Orchestra Status */}
      <div className="p-4 border-t border-slate-700">
        <div className="text-xs font-semibold text-slate-400 uppercase tracking-wide mb-3">AI Orchestra Status</div>
        <div className="space-y-2">
          <StatusItem label="Active Agents" value="7" />
          <StatusItem label="Processing" value="3 tasks" />
        </div>
      </div>

      {/* Additional content can be passed as children */}
      {children && (
        <div className="flex-1 flex flex-col">
          {children}
        </div>
      )}

      {/* Spacer to push content to bottom if needed */}
      <div className="flex-1"></div>
    </div>
  )
}
