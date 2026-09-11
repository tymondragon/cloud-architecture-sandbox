<script lang="ts">
  import type { Layer, GraphNode, GraphEdge, ComponentHealth } from '../../types/scenario';

  export let layers: Layer[];
  export let nodes: Record<string, GraphNode>;
  export let edges: Record<string, GraphEdge>;
  export let health: Record<string, ComponentHealth>;
  export let onNodeClick: (node: GraphNode) => void;
  export let hoveredLayer: string | null = null;
  export let compactCards: boolean = true;

  // Group nodes by layer
  $: nodesByLayer = (() => {
    const grouped: Record<string, GraphNode[]> = {};
    Object.values(nodes).forEach(node => {
      const layerId = node.namespace; // namespace is actually layer ID
      (grouped[layerId] ??= []).push(node);
    });
    return grouped;
  })();

  // Sort layers by position
  $: sortedLayers = [...layers].sort((a, b) => a.position - b.position);

  let containerEl: HTMLDivElement | undefined;
  const nodeEls: Record<string, HTMLElement> = {};

  interface EdgePath {
    id: string;
    d: string;
    color: string;
    animated: boolean;
  }

  let edgePaths: EdgePath[] = [];
  let svgSize = { width: 0, height: 0 };

  function computeEdges() {
    if (!containerEl) return;
    const containerRect = containerEl.getBoundingClientRect();

    edgePaths = Object.values(edges)
      .map((edge): EdgePath | null => {
        const fromEl = nodeEls[edge.source];
        const toEl = nodeEls[edge.target];
        if (!fromEl || !toEl) return null;

        const fromRect = fromEl.getBoundingClientRect();
        const toRect = toEl.getBoundingClientRect();

        // Start from bottom center of source node
        const x1 = fromRect.left - containerRect.left + fromRect.width / 2;
        const y1 = fromRect.bottom - containerRect.top;

        // End at top center of target node
        const x2 = toRect.left - containerRect.left + toRect.width / 2;
        const y2 = toRect.top - containerRect.top;

        // Vertical control points for smooth curves
        const controlY1 = y1 + (y2 - y1) / 3;
        const controlY2 = y1 + (2 * (y2 - y1)) / 3;

        return {
          id: edge.id,
          d: `M ${x1} ${y1} C ${x1} ${controlY1}, ${x2} ${controlY2}, ${x2} ${y2}`,
          color: edge.metadata.color || '#6b7280',
          animated: edge.metadata.animated === 'true'
        };
      })
      .filter((e): e is EdgePath => e !== null);

    svgSize = { width: containerEl.scrollWidth, height: containerEl.scrollHeight };
  }

  $: {
    nodes;
    edges;
    if (containerEl) {
      setTimeout(computeEdges, 100);
    }
  }

  $: if (containerEl) {
    const ro = new ResizeObserver(() => computeEdges());
    ro.observe(containerEl);
  }

  function handleNodeClick(node: GraphNode) {
    onNodeClick(node);
  }

  function getNodeIcon(labels: Record<string, string>): string {
    return labels.icon || '📦';
  }

  function getHealthStatus(componentId: string): ComponentHealth | undefined {
    return health[componentId];
  }

  function getHealthColor(status: string): string {
    switch (status) {
      case 'healthy':
        return '#10b981';
      case 'degraded':
        return '#f59e0b';
      case 'unhealthy':
        return '#ef4444';
      default:
        return '#6b7280';
    }
  }

  function isLayerHighlighted(layerId: string): boolean {
    return hoveredLayer === layerId;
  }
</script>

<div class="layout-container" bind:this={containerEl}>
  <svg class="edges-svg" width={svgSize.width} height={svgSize.height}>
    <defs>
      <filter id="edge-glow-v" x="-50%" y="-50%" width="200%" height="200%">
        <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
        <feMerge>
          <feMergeNode in="coloredBlur"/>
          <feMergeNode in="SourceGraphic"/>
        </feMerge>
      </filter>

      <marker id="arrowhead-v" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
        <path d="M0 0L10 5L0 10z" fill="currentColor" />
      </marker>
    </defs>

    {#each edgePaths as edge (edge.id)}
      <path
        d={edge.d}
        stroke={edge.color}
        stroke-width="2"
        fill="none"
        opacity="0.7"
        filter={edge.animated ? 'url(#edge-glow-v)' : 'none'}
        marker-end="url(#arrowhead-v)"
        style="color: {edge.color}"
        class:pulse-edge={edge.animated}
      >
        <title>{edge.id}</title>
      </path>
    {/each}
  </svg>

  <div class="components-flow">
    {#each sortedLayers as layer (layer.id)}
      {@const layerNodes = nodesByLayer[layer.id] || []}
      {#if layerNodes.length > 0}
        <div class="layer-section">
          {#each layerNodes as node (node.id)}
            {@const nodeHealth = getHealthStatus(node.id)}
            {@const highlighted = isLayerHighlighted(layer.id)}
            <div
              bind:this={nodeEls[node.id]}
              class="component-card"
              class:compact={compactCards}
              class:highlighted={highlighted}
              on:click={() => handleNodeClick(node)}
              role="button"
              tabindex="0"
            >
              {#if nodeHealth}
                <div
                  class="health-indicator"
                  style="background-color: {getHealthColor(nodeHealth.status)}"
                  title={nodeHealth.status}
                ></div>
              {/if}

              <div class="card-icon">
                {getNodeIcon(node.labels)}
              </div>

              <div class="card-content">
                <h4 class="card-name">{node.name}</h4>
                <p class="card-role">{node.labels.role || node.kind}</p>

                {#if !compactCards && node.metadata.description}
                  <p class="card-description">{node.metadata.description}</p>
                {/if}

                {#if nodeHealth?.replicas}
                  <div class="card-replicas">
                    <span class="replicas-badge">{nodeHealth.replicas}</span>
                  </div>
                {/if}
              </div>

              {#if node.metadata.actualTechnology}
                <div class="tech-badge" title="Implemented with {node.metadata.actualTechnology}">
                  💡 {node.metadata.actualTechnology}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      {/if}
    {/each}
  </div>
</div>

<style>
  .layout-container {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
  }

  .edges-svg {
    position: absolute;
    top: 0;
    left: 0;
    pointer-events: none;
    overflow: visible;
  }

  .pulse-edge {
    animation: pulse-glow 2s ease-in-out infinite;
  }

  @keyframes pulse-glow {
    0%, 100% {
      opacity: 0.5;
      stroke-width: 2;
    }
    50% {
      opacity: 1;
      stroke-width: 3;
    }
  }

  .components-flow {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 3rem;
    padding: 2rem;
    max-height: 100%;
    justify-content: center;
  }

  .layer-section {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .component-card {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    background: rgba(30, 41, 59, 0.8);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(71, 85, 105, 0.5);
    border-radius: 16px;
    padding: 1.5rem;
    min-width: 260px;
    max-width: 320px;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .component-card.compact {
    padding: 1rem;
    min-width: 220px;
    gap: 0.5rem;
  }

  .component-card:hover,
  .component-card.highlighted {
    transform: translateY(-4px);
    border-color: rgba(96, 165, 250, 0.8);
    box-shadow:
      0 12px 32px rgba(0, 0, 0, 0.5),
      0 0 24px rgba(96, 165, 250, 0.3);
  }

  .health-indicator {
    position: absolute;
    top: 12px;
    right: 12px;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    box-shadow: 0 0 8px currentColor;
  }

  .card-icon {
    font-size: 2.5rem;
    line-height: 1;
  }

  .component-card.compact .card-icon {
    font-size: 2rem;
  }

  .card-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.375rem;
    text-align: center;
    width: 100%;
  }

  .card-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #e2e8f0;
    margin: 0;
  }

  .component-card.compact .card-name {
    font-size: 1rem;
  }

  .card-role {
    font-size: 0.8125rem;
    font-weight: 500;
    color: #60a5fa;
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .component-card.compact .card-role {
    font-size: 0.75rem;
  }

  .card-description {
    font-size: 0.8125rem;
    line-height: 1.5;
    color: #cbd5e1;
    margin: 0;
  }

  .card-replicas {
    margin-top: 0.25rem;
  }

  .replicas-badge {
    display: inline-block;
    background: rgba(59, 130, 246, 0.2);
    color: #60a5fa;
    padding: 0.25rem 0.625rem;
    border-radius: 10px;
    font-size: 0.6875rem;
    font-weight: 600;
    border: 1px solid rgba(59, 130, 246, 0.3);
  }

  .tech-badge {
    margin-top: 0.25rem;
    padding: 0.375rem 0.625rem;
    background: rgba(168, 85, 247, 0.1);
    border: 1px solid rgba(168, 85, 247, 0.3);
    border-radius: 6px;
    font-size: 0.6875rem;
    color: #c4b5fd;
    text-align: center;
  }

  .component-card.compact .tech-badge {
    font-size: 0.625rem;
    padding: 0.25rem 0.5rem;
  }
</style>
