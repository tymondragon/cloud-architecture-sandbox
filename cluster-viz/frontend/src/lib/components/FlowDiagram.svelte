<script lang="ts">
  import { SvelteFlow, Controls, Background, MiniMap, type NodeTypes } from '@xyflow/svelte';
  import '@xyflow/svelte/dist/style.css';
  import { nodes, edges } from '../stores/graph';
  import KubeNode from './KubeNode.svelte';
  import type { GraphNode } from '../types/graph';

  export let onNodeClick: (node: GraphNode) => void;

  const nodeTypes: NodeTypes = {
    kubeNode: KubeNode,
  };

  // Simple conversion without derived stores
  let flowNodes: any[] = [];
  let flowEdges: any[] = [];

  $: {
    flowNodes = Object.values($nodes).map((node: any, i: number) => ({
      id: node.id,
      type: 'kubeNode',
      data: node,
      position: { x: (i % 5) * 200, y: Math.floor(i / 5) * 150 }
    }));

    flowEdges = Object.values($edges).map((edge: any) => ({
      id: edge.id,
      source: edge.source,
      target: edge.target,
      label: edge.label,
      animated: edge.kind === 'dataflow',
      style: edge.kind === 'dataflow' ? 'stroke: #f97316; stroke-width: 2;' : 'stroke: #3b82f6;'
    }));
  }

  function handleNodeClick(event: CustomEvent) {
    const node = event.detail.node;
    if (node.data) {
      onNodeClick(node.data);
    }
  }
</script>

<div class="flow-container">
  <SvelteFlow
    nodes={flowNodes}
    edges={flowEdges}
    {nodeTypes}
    fitView
    on:nodeclick={handleNodeClick}
  >
    <Controls />
    <Background />
    <MiniMap />
  </SvelteFlow>
</div>

<style>
  .flow-container {
    width: 100%;
    height: 100%;
  }

  :global(.svelte-flow) {
    background-color: #f9fafb;
  }

  :global(.edge-dataflow) {
    animation: dash 1s linear infinite;
  }

  @keyframes dash {
    to {
      stroke-dashoffset: -10;
    }
  }
</style>
