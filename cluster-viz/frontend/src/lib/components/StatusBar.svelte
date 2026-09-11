<script lang="ts">
  import { connectionStatus } from '../stores/websocket';
  import { stats } from '../stores/graph';

  function getStatusColor(status: string): string {
    switch (status) {
      case 'connected':
        return 'bg-emerald-500';
      case 'connecting':
        return 'bg-yellow-500';
      case 'error':
        return 'bg-red-500';
      default:
        return 'bg-slate-500';
    }
  }
</script>

<div class="status-bar">
  <div class="status-bar-left">
    <div class="status-indicator">
      <div class="status-dot {getStatusColor($connectionStatus)}"></div>
      <span class="status-text capitalize">{$connectionStatus}</span>
    </div>

    <div class="stat-item">
      <span class="stat-value">{$stats.nodeCount}</span>
      <span class="stat-label">nodes</span>
    </div>

    <div class="stat-item">
      <span class="stat-value">{$stats.edgeCount}</span>
      <span class="stat-label">edges</span>
    </div>
  </div>

  <div class="legend">
    <div class="legend-item">
      <div class="legend-line ownership"></div>
      <span>Ownership</span>
    </div>
    <div class="legend-item">
      <div class="legend-line routing"></div>
      <span>Routing</span>
    </div>
    <div class="legend-item">
      <div class="legend-line selection"></div>
      <span>Selection</span>
    </div>
    <div class="legend-item">
      <div class="legend-line dataflow"></div>
      <span>Data Flow</span>
    </div>
    <div class="legend-item">
      <div class="legend-line parentref"></div>
      <span>Parent Ref</span>
    </div>
  </div>
</div>

<style>
  .status-bar {
    background: linear-gradient(to bottom, #1e293b 0%, #0f172a 100%);
    border-bottom: 1px solid rgba(71, 85, 105, 0.5);
    padding: 0.75rem 1.5rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }

  .status-bar-left {
    display: flex;
    align-items: center;
    gap: 2rem;
  }

  .status-indicator {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .status-dot {
    width: 0.625rem;
    height: 0.625rem;
    border-radius: 50%;
    box-shadow: 0 0 8px currentColor;
  }

  .status-text {
    font-size: 0.875rem;
    font-weight: 600;
    color: #e2e8f0;
  }

  .stat-item {
    display: flex;
    align-items: baseline;
    gap: 0.375rem;
    font-size: 0.875rem;
  }

  .stat-value {
    font-weight: 700;
    color: #60a5fa;
  }

  .stat-label {
    color: #94a3b8;
  }

  .legend {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.75rem;
    color: #cbd5e1;
  }

  .legend-line {
    width: 2rem;
    height: 2px;
    border-radius: 1px;
  }

  .legend-line.ownership {
    background: #6b7280;
    border-top: 2px dashed #6b7280;
  }

  .legend-line.routing {
    background: #3b82f6;
    box-shadow: 0 0 4px #3b82f6;
  }

  .legend-line.selection {
    background: #6b7280;
    border-top: 2px dotted #6b7280;
  }

  .legend-line.dataflow {
    background: #f97316;
    box-shadow: 0 0 6px #f97316;
    animation: pulse-line 2s ease-in-out infinite;
  }

  .legend-line.parentref {
    background: #a855f7;
    box-shadow: 0 0 4px #a855f7;
  }

  @keyframes pulse-line {
    0%, 100% {
      opacity: 0.6;
    }
    50% {
      opacity: 1;
    }
  }
</style>
