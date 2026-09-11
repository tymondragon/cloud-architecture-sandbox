<script lang="ts">
  import type { GraphNode, GraphEdge } from '../types/graph';

  export let nodes: Record<string, GraphNode>;
  export let edges: Record<string, GraphEdge>;
  export let onNodeClick: (node: GraphNode) => void;

  // Group nodes by namespace into columns
  $: columns = (() => {
    const namespaces: Record<string, GraphNode[]> = {};

    Object.values(nodes).forEach(node => {
      const ns = node.namespace || 'cluster';
      (namespaces[ns] ??= []).push(node);
    });

    return Object.entries(namespaces).map(([name, nodes]) => ({ name, nodes }));
  })();

  let containerEl: HTMLDivElement | undefined;
  const nodeEls: Record<string, HTMLElement> = {};

  interface EdgePath {
    id: string;
    d: string;
    kind: string;
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
        const x1 = fromRect.right - containerRect.left;
        const y1 = fromRect.top - containerRect.top + fromRect.height / 2;
        const x2 = toRect.left - containerRect.left;
        const y2 = toRect.top - containerRect.top + toRect.height / 2;
        const bend = Math.max(32, (x2 - x1) / 2);

        return {
          id: edge.id,
          d: `M ${x1} ${y1} C ${x1 + bend} ${y1}, ${x2 - bend} ${y2}, ${x2} ${y2}`,
          kind: edge.kind
        };
      })
      .filter((e): e is EdgePath => e !== null);

    svgSize = { width: containerEl.scrollWidth, height: containerEl.scrollHeight };
  }

  $: {
    nodes;
    edges;
    if (containerEl) {
      computeEdges();
    }
  }

  $: if (containerEl) {
    const ro = new ResizeObserver(() => computeEdges());
    ro.observe(containerEl);
  }

  function handleNodeClick(node: GraphNode) {
    onNodeClick(node);
  }

  function getNodeGradient(kind: string): string {
    const gradients: Record<string, string> = {
      'Deployment': 'from-cyan-500 to-blue-600',
      'StatefulSet': 'from-purple-500 to-pink-600',
      'Service': 'from-emerald-500 to-teal-600',
      'Pod': 'from-blue-400 to-cyan-500',
      'HTTPRoute': 'from-indigo-500 to-purple-600',
      'Gateway': 'from-violet-500 to-fuchsia-600',
    };
    return gradients[kind] || 'from-slate-500 to-slate-600';
  }

  function getNodeIcon(kind: string): string {
    const icons: Record<string, string> = {
      'Deployment': 'DE',
      'StatefulSet': 'ST',
      'Service': 'SV',
      'Pod': 'PO',
      'HTTPRoute': 'HR',
      'Gateway': 'GW',
    };
    return icons[kind] || kind.substring(0, 2).toUpperCase();
  }

  function getEdgeColor(kind: string): string {
    const colors: Record<string, string> = {
      'dataflow': '#f97316',    // Orange
      'routing': '#3b82f6',     // Blue
      'parentRef': '#a855f7',   // Purple
      'selection': '#6b7280',   // Gray
      'ownership': '#6b7280',   // Gray
    };
    return colors[kind] || '#6b7280';
  }
</script>

<div class="graph-container">
  <div class="graph-wrapper" bind:this={containerEl}>
    <svg class="edges-svg" width={svgSize.width} height={svgSize.height}>
      <defs>
        <!-- Glow filters for edges -->
        <filter id="glow-orange" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
        <filter id="glow-blue" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
        <filter id="glow-purple" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>

        <!-- Arrow markers for each edge type -->
        <marker id="arrow-orange" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
          <path d="M0 0L10 5L0 10z" fill="#f97316" />
        </marker>
        <marker id="arrow-blue" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
          <path d="M0 0L10 5L0 10z" fill="#3b82f6" />
        </marker>
        <marker id="arrow-purple" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
          <path d="M0 0L10 5L0 10z" fill="#a855f7" />
        </marker>
        <marker id="arrow-gray" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto">
          <path d="M0 0L10 5L0 10z" fill="#6b7280" />
        </marker>
      </defs>

      {#each edgePaths as edge (edge.id)}
        <path
          d={edge.d}
          stroke={getEdgeColor(edge.kind)}
          stroke-width="2"
          fill="none"
          opacity="0.6"
          filter={edge.kind === 'dataflow' ? 'url(#glow-orange)' : edge.kind === 'routing' ? 'url(#glow-blue)' : edge.kind === 'parentRef' ? 'url(#glow-purple)' : 'none'}
          marker-end={edge.kind === 'dataflow' ? 'url(#arrow-orange)' : edge.kind === 'routing' ? 'url(#arrow-blue)' : edge.kind === 'parentRef' ? 'url(#arrow-purple)' : 'url(#arrow-gray)'}
          class:pulse-edge={edge.kind === 'dataflow'}
        >
          <title>{edge.kind}</title>
        </path>
      {/each}
    </svg>

    <div class="columns-container">
      {#each columns as col (col.name)}
        <div class="column">
          <div class="column-header">{col.name}</div>
          {#each col.nodes as node (node.id)}
            <div
              bind:this={nodeEls[node.id]}
              class="node"
              onclick={() => handleNodeClick(node)}
              role="button"
              tabindex="0"
            >
              <div class="node-icon bg-gradient-to-br {getNodeGradient(node.kind)}">
                {getNodeIcon(node.kind)}
              </div>
              <div class="node-content">
                <div class="node-name">{node.name}</div>
                <div class="node-kind">{node.kind}</div>
              </div>
              <div class="node-glow bg-gradient-to-br {getNodeGradient(node.kind)}"></div>
            </div>
          {/each}
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
    min-width: max-content;
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
      opacity: 0.4;
      stroke-width: 2;
    }
    50% {
      opacity: 0.8;
      stroke-width: 2.5;
    }
  }

  .columns-container {
    position: relative;
    display: flex;
    gap: 4rem;
    align-items: flex-start;
    padding: 1rem;
  }

  .column {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    min-width: 220px;
  }

  .column-header {
    font-size: 0.875rem;
    font-weight: 600;
    color: #60a5fa;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 0.5rem;
    text-shadow: 0 0 10px rgba(96, 165, 250, 0.5);
  }

  .node {
    position: relative;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background: rgba(30, 41, 59, 0.6);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(71, 85, 105, 0.5);
    border-radius: 12px;
    padding: 1rem;
    min-width: 220px;
    cursor: pointer;
    transition: all 0.3s ease;
    overflow: hidden;
  }

  .node:hover {
    transform: translateY(-2px);
    border-color: rgba(96, 165, 250, 0.8);
    box-shadow:
      0 10px 30px rgba(0, 0, 0, 0.5),
      0 0 20px rgba(96, 165, 250, 0.3);
  }

  .node:hover .node-glow {
    opacity: 0.15;
  }

  .node-glow {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    opacity: 0;
    transition: opacity 0.3s ease;
    pointer-events: none;
    border-radius: 12px;
  }

  .node-icon {
    flex: none;
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 700;
    color: white;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
    position: relative;
    z-index: 1;
  }

  .node-content {
    flex: 1;
    min-width: 0;
    position: relative;
    z-index: 1;
  }

  .node-name {
    font-size: 0.875rem;
    font-weight: 600;
    color: #e2e8f0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-bottom: 0.25rem;
  }

  .node-kind {
    font-size: 0.75rem;
    color: #94a3b8;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
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
