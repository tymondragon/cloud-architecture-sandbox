<script lang="ts">
  import { Handle, Position } from '@xyflow/svelte';
  import type { GraphNode } from '../types/graph';

  export let data: GraphNode;

  // Determine status color
  function getStatusColor(status: string): string {
    const s = status.toLowerCase();
    if (s === 'running' || s === 'active' || s === 'healthy' || s === 'ready' || s === 'accepted') {
      return 'bg-green-500';
    }
    if (s === 'pending' || s === 'degraded') {
      return 'bg-yellow-500';
    }
    if (s === 'error' || s === 'unavailable') {
      return 'bg-red-500';
    }
    if (s === 'terminated') {
      return 'bg-gray-500';
    }
    return 'bg-blue-500';
  }

  function getKindColor(kind: string): string {
    switch (kind) {
      case 'Namespace':
        return 'bg-purple-100 border-purple-400';
      case 'Deployment':
      case 'StatefulSet':
      case 'DaemonSet':
        return 'bg-blue-100 border-blue-400';
      case 'Pod':
        return 'bg-green-100 border-green-400';
      case 'Service':
        return 'bg-yellow-100 border-yellow-400';
      case 'Gateway':
        return 'bg-indigo-100 border-indigo-400';
      case 'HTTPRoute':
        return 'bg-pink-100 border-pink-400';
      default:
        return 'bg-gray-100 border-gray-400';
    }
  }
</script>

<div class="kube-node {getKindColor(data.kind)} border-2 rounded-lg shadow-md p-3 w-44">
  <Handle type="target" position={Position.Top} />

  <div class="flex items-center justify-between mb-2">
    <span class="text-xs font-semibold text-gray-700">{data.kind}</span>
    <div class="w-3 h-3 rounded-full {getStatusColor(data.status)}" title={data.status}></div>
  </div>

  <div class="text-sm font-medium text-gray-900 truncate" title={data.name}>
    {data.name}
  </div>

  {#if data.namespace}
    <div class="text-xs text-gray-600 truncate" title={data.namespace}>
      {data.namespace}
    </div>
  {/if}

  {#if data.metadata.replicas}
    <div class="text-xs text-gray-500 mt-1">
      {data.metadata.replicas}
    </div>
  {/if}

  <Handle type="source" position={Position.Bottom} />
</div>

<style>
  .kube-node {
    min-width: 180px;
    max-width: 180px;
  }
</style>
