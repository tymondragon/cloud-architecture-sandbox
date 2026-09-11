<script lang="ts">
  import type { LayoutConfig, Layer, GraphNode, GraphEdge, ComponentHealth } from '../types/scenario';
  import VerticalLayout from './layouts/VerticalLayout.svelte';
  import HorizontalLayout from './layouts/HorizontalLayout.svelte';

  export let layoutConfig: LayoutConfig;
  export let layers: Layer[];
  export let nodes: Record<string, GraphNode>;
  export let edges: Record<string, GraphEdge>;
  export let health: Record<string, ComponentHealth>;
  export let onNodeClick: (node: GraphNode) => void;
  export let hoveredLayer: string | null = null;

  $: compactCards = layoutConfig.options.compactCards ?? true;
</script>

<div class="layout-router">
  {#if layoutConfig.type === 'vertical'}
    <VerticalLayout
      {layers}
      {nodes}
      {edges}
      {health}
      {onNodeClick}
      {hoveredLayer}
      {compactCards}
    />
  {:else if layoutConfig.type === 'horizontal'}
    <HorizontalLayout
      {layers}
      {nodes}
      {edges}
      {health}
      {onNodeClick}
      {hoveredLayer}
      {compactCards}
    />
  {:else if layoutConfig.type === 'radial'}
    <div class="placeholder">Radial layout coming soon</div>
  {:else if layoutConfig.type === 'force-graph'}
    <div class="placeholder">Force-graph layout coming soon</div>
  {:else if layoutConfig.type === 'grid'}
    <div class="placeholder">Grid layout coming soon</div>
  {:else}
    <div class="placeholder">Unknown layout type: {layoutConfig.type}</div>
  {/if}
</div>

<style>
  .layout-router {
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #94a3b8;
    font-size: 1.125rem;
  }
</style>
