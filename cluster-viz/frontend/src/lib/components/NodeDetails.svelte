<script lang="ts">
  import type { GraphNode } from '../types/graph';
  import { nodes, edges } from '../stores/graph';

  export let selectedNode: GraphNode | null;
  export let onClose: () => void;

  $: connectedEdges = selectedNode
    ? Object.values($edges).filter(
        (e) => e.source === selectedNode.id || e.target === selectedNode.id
      )
    : [];
</script>

{#if selectedNode}
  <div class="fixed right-0 top-10 bottom-0 w-96 bg-white border-l border-gray-200 shadow-xl overflow-y-auto">
    <div class="sticky top-0 bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
      <h2 class="text-lg font-semibold text-gray-900">Node Details</h2>
      <button
        on:click={onClose}
        class="text-gray-400 hover:text-gray-600"
        aria-label="Close"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <div class="p-4 space-y-4">
      <div>
        <h3 class="text-sm font-semibold text-gray-700 mb-1">Name</h3>
        <p class="text-sm text-gray-900">{selectedNode.name}</p>
      </div>

      <div>
        <h3 class="text-sm font-semibold text-gray-700 mb-1">Kind</h3>
        <p class="text-sm text-gray-900">{selectedNode.kind}</p>
      </div>

      {#if selectedNode.namespace}
        <div>
          <h3 class="text-sm font-semibold text-gray-700 mb-1">Namespace</h3>
          <p class="text-sm text-gray-900">{selectedNode.namespace}</p>
        </div>
      {/if}

      <div>
        <h3 class="text-sm font-semibold text-gray-700 mb-1">Status</h3>
        <p class="text-sm text-gray-900">{selectedNode.status}</p>
      </div>

      {#if Object.keys(selectedNode.labels).length > 0}
        <div>
          <h3 class="text-sm font-semibold text-gray-700 mb-1">Labels</h3>
          <div class="space-y-1">
            {#each Object.entries(selectedNode.labels) as [key, value]}
              <div class="text-xs bg-gray-100 px-2 py-1 rounded">
                <span class="font-medium">{key}:</span>
                <span class="text-gray-700">{value}</span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if Object.keys(selectedNode.metadata).length > 0}
        <div>
          <h3 class="text-sm font-semibold text-gray-700 mb-1">Metadata</h3>
          <div class="space-y-1">
            {#each Object.entries(selectedNode.metadata) as [key, value]}
              <div class="text-xs">
                <span class="font-medium text-gray-700">{key}:</span>
                <span class="text-gray-600 break-all">{value}</span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if connectedEdges.length > 0}
        <div>
          <h3 class="text-sm font-semibold text-gray-700 mb-1">Connected Edges ({connectedEdges.length})</h3>
          <div class="space-y-2">
            {#each connectedEdges as edge}
              <div class="text-xs bg-gray-50 px-2 py-1 rounded border border-gray-200">
                <div class="font-medium text-gray-900">{edge.kind}</div>
                <div class="text-gray-600">
                  {edge.source === selectedNode.id ? '→' : '←'}
                  {edge.source === selectedNode.id ? edge.target : edge.source}
                </div>
                {#if edge.label}
                  <div class="text-gray-500 italic">{edge.label}</div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}
