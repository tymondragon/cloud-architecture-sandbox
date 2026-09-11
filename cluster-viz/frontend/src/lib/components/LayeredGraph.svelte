<script lang="ts">
  import type { Layer } from '../types/scenario';
  import type { GraphNode, GraphEdge, ComponentHealth } from '../types/scenario';

  export let layers: Layer[];
  export let nodes: Record<string, GraphNode>;
  export let edges: Record<string, GraphEdge>;
  export let health: Record<string, ComponentHealth>;
  export let onNodeClick: (node: GraphNode) => void;

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
    kind: string;
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
          kind: edge.kind,
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
    layers;
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
</script>

<div class="graph-container">
  <div class="graph-wrapper" bind:this={containerEl}>
    <svg class="edges-svg" width={svgSize.width} height={svgSize.height}>
      <defs>
        <!-- Glow filters -->
        <filter id="edge-glow" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>

        <!-- Arrow markers -->
        <marker id="arrowhead-orange" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
          <path d="M0 0L10 5L0 10z" fill="#f97316" />
        </marker>
        <marker id="arrowhead-blue" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
          <path d="M0 0L10 5L0 10z" fill="#3b82f6" />
        </marker>
        <marker id="arrowhead-purple" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
          <path d="M0 0L10 5L0 10z" fill="#a855f7" />
        </marker>
        <marker id="arrowhead-gray" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
          <path d="M0 0L10 5L0 10z" fill="#6b7280" />
        </marker>
      </defs>

      {#each edgePaths as edge (edge.id)}
        <path
          d={edge.d}
          stroke={edge.color}
          stroke-width="2"
          fill="none"
          opacity="0.7"
          filter={edge.animated ? 'url(#edge-glow)' : 'none'}
          marker-end="url(#arrowhead-{edge.color === '#f97316' ? 'orange' : edge.color === '#3b82f6' ? 'blue' : edge.color === '#a855f7' ? 'purple' : 'gray'})"
          class:pulse-edge={edge.animated}
        >
          <title>{edge.kind}</title>
        </path>
      {/each}
    </svg>

    <div class="layers-container">
      {#each sortedLayers as layer (layer.id)}
        <div class="layer" style="border-left-color: {layer.color}">
          <div class="layer-header">
            <h3 class="layer-name">{layer.name}</h3>
            <p class="layer-description">{layer.description}</p>
          </div>

          <div class="layer-components">
            {#each nodesByLayer[layer.id] || [] as node (node.id)}
              {@const nodeHealth = getHealthStatus(node.id)}
              <div
                bind:this={nodeEls[node.id]}
                class="component-card"
                on:click={() => handleNodeClick(node)}
                role="button"
                tabindex="0"
              >
                <!-- Health indicator -->
                {#if nodeHealth}
                  <div
                    class="health-indicator"
                    style="background-color: {getHealthColor(nodeHealth.status)}"
                    title={nodeHealth.status}
                  ></div>
                {/if}

                <div class="component-icon">
                  {getNodeIcon(node.labels)}
                </div>

                <div class="component-content">
                  <h4 class="component-name">{node.name}</h4>
                  <p class="component-role">{node.labels.role || node.kind}</p>

                  {#if node.metadata.description}
                    <p class="component-description">{node.metadata.description}</p>
                  {/if}

                  {#if nodeHealth?.replicas}
                    <div class="component-replicas">
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
        </div>
      {/each}
    </div>
  </div>
</div>

<style>
  .graph-container {
    height: 100%;
    overflow: auto;
    background: linear-gradient(135deg, #0a0e27 0%, #1a1f3a 100%);
    padding: 2rem;
  }

  .graph-wrapper {
    position: relative;
    min-width: 100%;
    min-height: 100%;
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

  .layers-container {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 3rem;
    padding: 1rem;
  }

  .layer {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding: 1.5rem;
    background: rgba(30, 41, 59, 0.3);
    border-left: 4px solid;
    border-radius: 12px;
  }

  .layer-header {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .layer-name {
    font-size: 1.125rem;
    font-weight: 600;
    color: #e2e8f0;
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .layer-description {
    font-size: 0.875rem;
    color: #94a3b8;
    margin: 0;
  }

  .layer-components {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    justify-content: center;
    align-items: flex-start;
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
    min-width: 280px;
    max-width: 320px;
    cursor: pointer;
    transition: all 0.3s ease;
    overflow: visible;
  }

  .component-card:hover {
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

  .component-icon {
    font-size: 3rem;
    line-height: 1;
  }

  .component-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    text-align: center;
    width: 100%;
  }

  .component-name {
    font-size: 1.25rem;
    font-weight: 600;
    color: #e2e8f0;
    margin: 0;
  }

  .component-role {
    font-size: 0.875rem;
    font-weight: 500;
    color: #60a5fa;
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .component-description {
    font-size: 0.875rem;
    line-height: 1.5;
    color: #cbd5e1;
    margin: 0;
  }

  .component-replicas {
    margin-top: 0.5rem;
  }

  .replicas-badge {
    display: inline-block;
    background: rgba(59, 130, 246, 0.2);
    color: #60a5fa;
    padding: 0.25rem 0.75rem;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 600;
    border: 1px solid rgba(59, 130, 246, 0.3);
  }

  .tech-badge {
    margin-top: 0.5rem;
    padding: 0.5rem 0.75rem;
    background: rgba(168, 85, 247, 0.1);
    border: 1px solid rgba(168, 85, 247, 0.3);
    border-radius: 8px;
    font-size: 0.75rem;
    color: #c4b5fd;
    text-align: center;
  }

  /* Scrollbar styling */
  .graph-container::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }

  .graph-container::-webkit-scrollbar-track {
    background: rgba(15, 23, 42, 0.5);
  }

  .graph-container::-webkit-scrollbar-thumb {
    background: rgba(71, 85, 105, 0.5);
    border-radius: 4px;
  }

  .graph-container::-webkit-scrollbar-thumb:hover {
    background: rgba(71, 85, 105, 0.8);
  }
</style>
