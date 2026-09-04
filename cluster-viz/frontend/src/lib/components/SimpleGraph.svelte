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
        const bend = Math.max(24, (x2 - x1) / 2);

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
    nodes; // Track nodes
    edges; // Track edges
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

  function statusColor(status: string): string {
    if (status === 'Healthy') return 'bg-green-100 border-green-300';
    if (status === 'Unavailable') return 'bg-red-100 border-red-300';
    return 'bg-gray-100 border-gray-300';
  }

  function edgeColor(kind: string): string {
    if (kind === 'dataflow') return 'stroke-orange-500';
    if (kind === 'routing') return 'stroke-blue-500';
    if (kind === 'parentRef') return 'stroke-purple-500';
    return 'stroke-gray-400';
  }
</script>

<div class="overflow-auto border border-gray-200 rounded-lg p-6 bg-gray-50 h-full">
  <div class="relative min-w-max" bind:this={containerEl}>
    <svg class="absolute top-0 left-0 pointer-events-none overflow-visible" width={svgSize.width} height={svgSize.height}>
      <defs>
        <marker id="arrow" viewBox="0 0 10 10" refX="8.5" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
          <path d="M0 0L10 5L0 10z" fill="currentColor" />
        </marker>
      </defs>
      {#each edgePaths as edge (edge.id)}
        <path
          d={edge.d}
          class="{edgeColor(edge.kind)} fill-none stroke-2 opacity-60"
          class:animate-pulse={edge.kind === 'dataflow'}
          marker-end="url(#arrow)"
        >
          <title>{edge.kind}</title>
        </path>
      {/each}
    </svg>

    <div class="relative flex gap-14 items-start">
      {#each columns as col (col.name)}
        <div class="flex flex-col gap-4">
          <div class="text-sm font-semibold text-gray-700 mb-2">{col.name}</div>
          {#each col.nodes as node (node.id)}
            <div
              bind:this={nodeEls[node.id]}
              class="relative flex items-center gap-2 {statusColor(node.status)} border rounded-lg p-3 min-w-[12rem] max-w-[16rem] cursor-pointer hover:shadow-md transition-shadow"
              onclick={() => handleNodeClick(node)}
              role="button"
              tabindex="0"
            >
              <div class="flex-none w-8 h-8 bg-blue-500 text-white rounded flex items-center justify-center text-xs font-bold">
                {node.kind.substring(0, 2).toUpperCase()}
              </div>
              <div class="min-w-0 flex-1">
                <div class="text-sm font-medium text-gray-900 truncate">{node.name}</div>
                <div class="text-xs text-gray-600 truncate">{node.kind}</div>
              </div>
            </div>
          {/each}
        </div>
      {/each}
    </div>
  </div>
</div>
