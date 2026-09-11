<script lang="ts">
  import { onMount } from 'svelte';
  import {
    currentScenario,
    logicalGraph,
    showPatternOverlay,
    fetchScenarios,
    fetchScenario,
    fetchLogicalGraph,
    currentScenarioId
  } from './lib/stores/scenario';
  import { connectionStatus } from './lib/stores/websocket';
  import PatternOverlay from './lib/components/PatternOverlay.svelte';
  import LayerSidebar from './lib/components/LayerSidebar.svelte';
  import LayoutRouter from './lib/components/LayoutRouter.svelte';
  import ScenarioSwitcher from './lib/components/ScenarioSwitcher.svelte';
  import NodeDetails from './lib/components/NodeDetails.svelte';
  import type { GraphNode } from './lib/types/scenario';

  let selectedNode: GraphNode | null = null;
  let hoveredLayer: string | null = null;

  function handleNodeClick(node: GraphNode) {
    selectedNode = node;
  }

  function handleCloseDetails() {
    selectedNode = null;
  }

  function getConnectionStatusColor(status: string): string {
    switch (status) {
      case 'connected':
        return '#10b981';
      case 'connecting':
        return '#f59e0b';
      case 'error':
        return '#ef4444';
      default:
        return '#6b7280';
    }
  }

  onMount(async () => {
    // Fetch available scenarios
    await fetchScenarios();

    // Load initial scenario
    const scenarioId = $currentScenarioId;
    await fetchScenario(scenarioId);
    await fetchLogicalGraph(scenarioId);

    // Poll for graph updates
    const intervalId = setInterval(async () => {
      await fetchLogicalGraph($currentScenarioId);
    }, 5000); // Update every 5 seconds

    return () => {
      clearInterval(intervalId);
    };
  });
</script>

<main class="app-container">
  <!-- Pattern Overlay -->
  {#if $showPatternOverlay && $currentScenario}
    <PatternOverlay
      bind:visible={$showPatternOverlay}
      metadata={$currentScenario.metadata}
      patterns={$currentScenario.patterns}
      technologies={$currentScenario.technologies}
    />
  {/if}

  <!-- Top Navigation -->
  <header class="app-header">
    <div class="header-left">
      <h1 class="app-title">Architecture Patterns Visualizer</h1>
    </div>

    <div class="header-center">
      <ScenarioSwitcher />
    </div>

    <div class="header-right">
      <button
        class="help-button"
        on:click={() => $showPatternOverlay = true}
        title="Show pattern information"
      >
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z" clip-rule="evenodd" />
        </svg>
        Help
      </button>

      <div class="connection-status">
        <div
          class="status-dot"
          style="background-color: {getConnectionStatusColor($connectionStatus)}"
        ></div>
        <span class="status-text">{$connectionStatus}</span>
      </div>
    </div>
  </header>

  <!-- Main Content with Sidebar -->
  <div class="app-main">
    {#if $logicalGraph && $currentScenario}
      <!-- Layer Sidebar -->
      <LayerSidebar
        layers={$logicalGraph.layers}
        bind:hoveredLayer={hoveredLayer}
      />

      <!-- Visualization Area -->
      <div class="visualization-area">
        <LayoutRouter
          layoutConfig={$currentScenario.layout}
          layers={$logicalGraph.layers}
          nodes={$logicalGraph.components}
          edges={$logicalGraph.edges}
          health={$logicalGraph.health}
          onNodeClick={handleNodeClick}
          {hoveredLayer}
        />
      </div>
    {:else}
      <div class="loading-container">
        <div class="loading-spinner"></div>
        <p class="loading-text">Loading scenario...</p>
      </div>
    {/if}

    <!-- Node Details Panel -->
    <NodeDetails {selectedNode} onClose={handleCloseDetails} />
  </div>

  <!-- Stats Footer -->
  {#if $logicalGraph}
    <footer class="app-footer">
      <div class="stat-item">
        <span class="stat-value">{Object.keys($logicalGraph.components).length}</span>
        <span class="stat-label">components</span>
      </div>
      <div class="stat-item">
        <span class="stat-value">{Object.keys($logicalGraph.edges).length}</span>
        <span class="stat-label">flows</span>
      </div>
      <div class="stat-item">
        <span class="stat-value">{$logicalGraph.layers.length}</span>
        <span class="stat-label">layers</span>
      </div>
    </footer>
  {/if}
</main>

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    background: linear-gradient(135deg, #0a0e27 0%, #1a1f3a 100%);
  }

  .app-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 1.5rem;
    background: linear-gradient(to bottom, #1e293b 0%, #0f172a 100%);
    border-bottom: 1px solid rgba(71, 85, 105, 0.5);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
    gap: 1rem;
    flex-shrink: 0;
  }

  .header-left {
    flex: 1;
  }

  .app-title {
    font-size: 1.25rem;
    font-weight: 700;
    color: #e2e8f0;
    margin: 0;
    background: linear-gradient(135deg, #60a5fa 0%, #a78bfa 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .header-center {
    display: flex;
    justify-content: center;
  }

  .header-right {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 1.5rem;
  }

  .help-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: rgba(59, 130, 246, 0.1);
    border: 1px solid rgba(59, 130, 246, 0.3);
    border-radius: 8px;
    padding: 0.5rem 1rem;
    color: #60a5fa;
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .help-button:hover {
    background: rgba(59, 130, 246, 0.2);
    border-color: rgba(59, 130, 246, 0.5);
  }

  .connection-status {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: rgba(30, 41, 59, 0.6);
    border: 1px solid rgba(71, 85, 105, 0.3);
    border-radius: 8px;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    box-shadow: 0 0 8px currentColor;
  }

  .status-text {
    font-size: 0.875rem;
    color: #94a3b8;
    text-transform: capitalize;
  }

  .app-main {
    flex: 1;
    display: flex;
    position: relative;
    overflow: hidden;
  }

  .visualization-area {
    flex: 1;
    position: relative;
    overflow: hidden;
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    gap: 1rem;
  }

  .loading-spinner {
    width: 48px;
    height: 48px;
    border: 4px solid rgba(96, 165, 250, 0.2);
    border-top-color: #60a5fa;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .loading-text {
    font-size: 1rem;
    color: #94a3b8;
  }

  .app-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 3rem;
    padding: 0.75rem 1.5rem;
    background: linear-gradient(to bottom, #1e293b 0%, #0f172a 100%);
    border-top: 1px solid rgba(71, 85, 105, 0.5);
    box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.3);
    flex-shrink: 0;
  }

  .stat-item {
    display: flex;
    align-items: baseline;
    gap: 0.5rem;
  }

  .stat-value {
    font-size: 1.25rem;
    font-weight: 700;
    color: #60a5fa;
  }

  .stat-label {
    font-size: 0.875rem;
    color: #94a3b8;
  }
</style>
