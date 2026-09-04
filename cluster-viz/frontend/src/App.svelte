<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { connectWebSocket, disconnectWebSocket } from './lib/stores/websocket';
  import { nodes, edges } from './lib/stores/graph';
  import StatusBar from './lib/components/StatusBar.svelte';
  import SimpleGraph from './lib/components/SimpleGraph.svelte';
  import NodeDetails from './lib/components/NodeDetails.svelte';
  import type { GraphNode } from './lib/types/graph';

  let selectedNode: GraphNode | null = null;

  function handleNodeClick(node: GraphNode) {
    selectedNode = node;
  }

  function handleCloseDetails() {
    selectedNode = null;
  }

  onMount(() => {
    connectWebSocket();
  });

  onDestroy(() => {
    disconnectWebSocket();
  });
</script>

<main class="h-screen w-screen flex flex-col">
  <StatusBar />
  <div class="flex-1 relative">
    <SimpleGraph nodes={$nodes} edges={$edges} onNodeClick={handleNodeClick} />
    <NodeDetails {selectedNode} onClose={handleCloseDetails} />
  </div>
</main>
