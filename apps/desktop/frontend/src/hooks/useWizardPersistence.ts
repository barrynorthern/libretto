import { useState, useEffect, useCallback } from 'react'

export interface ProjectData {
  title: string
  description: string
  genre: string
  coreTheme: string
  template: string
  tone: string
  notes?: string
}

const STORAGE_KEY = 'libretto-wizard-data'

export function useWizardPersistence() {
  const [projectData, setProjectData] = useState<Partial<ProjectData>>({
    title: '',
    description: '',
    genre: '',
    coreTheme: '',
    template: '',
    tone: '',
    notes: ''
  })

  const [step, setStep] = useState(1)

  // Load data from localStorage on mount
  useEffect(() => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        const parsed = JSON.parse(saved)
        setProjectData(parsed.projectData || {})
        setStep(parsed.step || 1)
      }
    } catch (error) {
      console.warn('Failed to load wizard data from localStorage:', error)
    }
  }, [])

  // Save data to localStorage whenever it changes
  useEffect(() => {
    try {
      const dataToSave = {
        projectData,
        step,
        lastSaved: new Date().toISOString()
      }
      localStorage.setItem(STORAGE_KEY, JSON.stringify(dataToSave))
    } catch (error) {
      console.warn('Failed to save wizard data to localStorage:', error)
    }
  }, [projectData, step])

  const updateData = useCallback((updates: Partial<ProjectData>) => {
    setProjectData(prev => ({ ...prev, ...updates }))
  }, [])

  const clearData = useCallback(() => {
    const emptyData = {
      title: '',
      description: '',
      genre: '',
      coreTheme: '',
      template: '',
      tone: '',
      notes: ''
    }
    setProjectData(emptyData)
    setStep(1)
    try {
      localStorage.removeItem(STORAGE_KEY)
    } catch (error) {
      console.warn('Failed to clear wizard data from localStorage:', error)
    }
  }, [])

  const hasData = useCallback(() => {
    return Object.values(projectData).some(value => value && value.trim() !== '')
  }, [projectData])

  const goToStep = useCallback((newStep: number) => {
    setStep(newStep)
  }, [])

  return {
    projectData,
    step,
    updateData,
    clearData,
    hasData,
    goToStep
  }
}
