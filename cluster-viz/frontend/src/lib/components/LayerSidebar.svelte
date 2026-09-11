<script lang="ts">
  import type { Layer } from '../types/scenario';

  export let layers: Layer[];
  export let hoveredLayer: string | null = null;

  function handleLayerHover(layerId: string) {
    hoveredLayer = layerId;
  }

  function handleLayerLeave() {
    hoveredLayer = null;
  }

  // Sort layers by position
  $: sortedLayers = [...layers].sort((a, b) => a.position - b.position);
</script>

<aside class="layer-sidebar">
  <div class="sidebar-header">
    <h2 class="sidebar-title">Architecture Layers</h2>
  </div>

  <div class="layers-list">
    {#each sortedLayers as layer (layer.id)}
      <button
        class="layer-item"
        class:hovered={hoveredLayer === layer.id}
        style="--layer-color: {layer.color}"
        on:mouseenter={() => handleLayerHover(layer.id)}
        on:mouseleave={handleLayerLeave}
      >
        <div class="layer-indicator" style="background-color: {layer.color}"></div>
        <div class="layer-content">
          <h3 class="layer-name">{layer.name}</h3>
          <p class="layer-description">{layer.description}</p>
        </div>
      </button>
    {/each}
  </div>

  <div class="sidebar-footer">
    <p class="sidebar-hint">
      Hover over a layer to highlight its components
    </p>
  </div>
</aside>

<style>
  .layer-sidebar {
    width: 280px;
    height: 100%;
    background: linear-gradient(to bottom, #1e293b 0%, #0f172a 100%);
    border-right: 1px solid rgba(71, 85, 105, 0.5);
    display: flex;
    flex-direction: column;
    overflow-y: auto;
  }

  .sidebar-header {
    padding: 1.5rem 1.25rem 1rem;
    border-bottom: 1px solid rgba(71, 85, 105, 0.3);
  }

  .sidebar-title {
    font-size: 0.875rem;
    font-weight: 700;
    color: #e2e8f0;
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .layers-list {
    flex: 1;
    padding: 1rem 0.75rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .layer-item {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    background: rgba(30, 41, 59, 0.4);
    border: 1px solid rgba(71, 85, 105, 0.3);
    border-radius: 8px;
    padding: 1rem;
    cursor: pointer;
    transition: all 0.2s;
    text-align: left;
    width: 100%;
  }

  .layer-item:hover,
  .layer-item.hovered {
    background: rgba(30, 41, 59, 0.8);
    border-color: var(--layer-color);
    box-shadow: 0 0 12px rgba(0, 0, 0, 0.3);
    transform: translateX(4px);
  }

  .layer-indicator {
    width: 4px;
    min-height: 100%;
    border-radius: 2px;
    box-shadow: 0 0 8px currentColor;
  }

  .layer-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .layer-name {
    font-size: 0.9375rem;
    font-weight: 600;
    color: #e2e8f0;
    margin: 0;
  }

  .layer-description {
    font-size: 0.8125rem;
    line-height: 1.4;
    color: #94a3b8;
    margin: 0;
  }

  .sidebar-footer {
    padding: 1rem 1.25rem;
    border-top: 1px solid rgba(71, 85, 105, 0.3);
  }

  .sidebar-hint {
    font-size: 0.75rem;
    color: #64748b;
    margin: 0;
    text-align: center;
    font-style: italic;
  }

  /* Scrollbar styling */
  .layer-sidebar::-webkit-scrollbar {
    width: 6px;
  }

  .layer-sidebar::-webkit-scrollbar-track {
    background: rgba(15, 23, 42, 0.5);
  }

  .layer-sidebar::-webkit-scrollbar-thumb {
    background: rgba(71, 85, 105, 0.5);
    border-radius: 3px;
  }

  .layer-sidebar::-webkit-scrollbar-thumb:hover {
    background: rgba(71, 85, 105, 0.8);
  }
</style>
