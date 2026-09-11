<script lang="ts">
  import type { Layer, GraphNode, GraphEdge, ComponentHealth } from '../../types/scenario';

  export let layers: Layer[];
  export let nodes: Record<string, GraphNode>;
  export let edges: Record<string, GraphEdge>;
  export let health: Record<string, ComponentHealth>;
  export let onNodeClick: (node: GraphNode) => void;
  export let hoveredLayer: string | null = null;
  export let compactCards: boolean = true;

  // TODO: Implement horizontal pipeline layout
  // For now, show a placeholder
</script>

<div class="horizontal-layout">
  <div class="placeholder">
    <h2>Horizontal Layout</h2>
    <p>Pipeline visualization (coming soon)</p>
  </div>
</div>

<style>
  .horizontal-layout {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .placeholder {
    text-align: center;
    color: #94a3b8;
  }

  .placeholder h2 {
    font-size: 1.5rem;
    color: #e2e8f0;
    margin-bottom: 0.5rem;
  }
</style>
