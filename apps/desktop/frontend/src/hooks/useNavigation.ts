import { useState, useCallback } from 'react'

export type NavigationScreen = 'canvas' | 'new-project' | 'ai-orchestra'

interface NavigationState {
  currentScreen: NavigationScreen
  canNavigate: boolean
  isInWizardFlow: boolean
}

interface NavigationActions {
  navigateTo: (screen: NavigationScreen) => void
  startWizardFlow: () => void
  exitWizardFlow: () => void
  setCanNavigate: (canNavigate: boolean) => void
}

export function useNavigation(): NavigationState & NavigationActions {
  const [currentScreen, setCurrentScreen] = useState<NavigationScreen>('canvas')
  const [canNavigate, setCanNavigate] = useState(true)
  const [isInWizardFlow, setIsInWizardFlow] = useState(false)

  const navigateTo = useCallback((screen: NavigationScreen) => {
    setCurrentScreen(screen)

    // Handle special navigation logic
    if (screen === 'new-project') {
      setIsInWizardFlow(true)
    } else {
      setIsInWizardFlow(false)
    }
  }, [])

  const startWizardFlow = useCallback(() => {
    setCurrentScreen('new-project')
    setIsInWizardFlow(true)
  }, [])

  const exitWizardFlow = useCallback(() => {
    setIsInWizardFlow(false)
    setCurrentScreen('canvas') // Return to main screen
  }, [])

  return {
    currentScreen,
    canNavigate: true, // Always allow navigation now
    isInWizardFlow,
    navigateTo,
    startWizardFlow,
    exitWizardFlow,
    setCanNavigate
  }
}
