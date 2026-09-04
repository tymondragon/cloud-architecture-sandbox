<script lang="ts">
  import { connectionStatus } from '../stores/websocket';
  import { stats } from '../stores/graph';

  function getStatusColor(status: string): string {
    switch (status) {
      case 'connected':
        return 'bg-green-500';
      case 'connecting':
        return 'bg-yellow-500';
      case 'error':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
    }
  }
</script>

<div class="bg-white border-b border-gray-200 px-4 py-2 flex items-center justify-between">
  <div class="flex items-center gap-6">
    <div class="flex items-center gap-2">
      <div class="w-3 h-3 rounded-full {getStatusColor($connectionStatus)}"></div>
      <span class="text-sm font-medium text-gray-700 capitalize">{$connectionStatus}</span>
    </div>

    <div class="text-sm text-gray-600">
      <span class="font-medium">{$stats.nodeCount}</span> nodes
    </div>

    <div class="text-sm text-gray-600">
      <span class="font-medium">{$stats.edgeCount}</span> edges
    </div>
  </div>

  <div class="flex items-center gap-4 text-xs text-gray-600">
    <div class="flex items-center gap-1">
      <div class="w-8 h-0.5 bg-gray-400" style="border-top: 2px dashed;"></div>
      <span>Ownership</span>
    </div>
    <div class="flex items-center gap-1">
      <div class="w-8 h-0.5 bg-blue-500"></div>
      <span>Routing</span>
    </div>
    <div class="flex items-center gap-1">
      <div class="w-8 h-0.5 bg-gray-400" style="border-top: 2px dotted;"></div>
      <span>Selection</span>
    </div>
    <div class="flex items-center gap-1">
      <div class="w-8 h-0.5 bg-orange-500 animate-pulse"></div>
      <span>Data Flow</span>
    </div>
    <div class="flex items-center gap-1">
      <div class="w-8 h-0.5 bg-purple-500"></div>
      <span>Parent Ref</span>
    </div>
  </div>
</div>
