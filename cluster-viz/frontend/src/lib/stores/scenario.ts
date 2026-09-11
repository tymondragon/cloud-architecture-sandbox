import { writable, derived } from 'svelte/store';
import type { Scenario, ScenarioSummary, LogicalGraph } from '../types/scenario';

// Current scenario ID
export const currentScenarioId = writable<string>('scenario-01');

// Available scenarios
export const scenarios = writable<ScenarioSummary[]>([]);

// Current scenario metadata
export const currentScenario = writable<Scenario | null>(null);

// Logical graph for current scenario
export const logicalGraph = writable<LogicalGraph | null>(null);

// Pattern overlay visibility
export const showPatternOverlay = writable<boolean>(true);

// Initialize: Check localStorage for overlay preference
if (typeof localStorage !== 'undefined') {
  const hideOverlay = localStorage.getItem('hidePatternOverlay');
  if (hideOverlay === 'true') {
    showPatternOverlay.set(false);
  }
}

// Fetch available scenarios
export async function fetchScenarios() {
  try {
    const response = await fetch('/api/scenarios');
    if (!response.ok) {
      throw new Error('Failed to fetch scenarios');
    }
    const data: ScenarioSummary[] = await response.json();
    scenarios.set(data);
  } catch (error) {
    console.error('Error fetching scenarios:', error);
  }
}

// Fetch specific scenario metadata
export async function fetchScenario(scenarioId: string) {
  try {
    const response = await fetch(`/api/scenario/${scenarioId}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch scenario ${scenarioId}`);
    }
    const data: Scenario = await response.json();
    currentScenario.set(data);
    return data;
  } catch (error) {
    console.error(`Error fetching scenario ${scenarioId}:`, error);
    return null;
  }
}

// Fetch logical graph for scenario
export async function fetchLogicalGraph(scenarioId: string) {
  try {
    const response = await fetch(`/api/graph?scenario=${scenarioId}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch graph for ${scenarioId}`);
    }
    const data: LogicalGraph = await response.json();
    logicalGraph.set(data);
    return data;
  } catch (error) {
    console.error(`Error fetching graph for ${scenarioId}:`, error);
    return null;
  }
}

// Switch to a different scenario
export async function switchScenario(scenarioId: string) {
  currentScenarioId.set(scenarioId);
  showPatternOverlay.set(true); // Show overlay when switching
  await fetchScenario(scenarioId);
  await fetchLogicalGraph(scenarioId);
}
