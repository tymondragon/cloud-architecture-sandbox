import { writable, derived } from 'svelte/store';
import type { GraphNode, GraphEdge, GraphEvent } from '../types/graph';

// Raw graph state
export const nodes = writable<Record<string, GraphNode>>({});
export const edges = writable<Record<string, GraphEdge>>({});

// Apply graph event from WebSocket
export function applyEvent(event: GraphEvent) {
  switch (event.type) {
    case 'SYNC':
      // Full graph sync
      if (event.graph) {
        nodes.set(event.graph.nodes);
        edges.set(event.graph.edges);
      }
      break;

    case 'ADDED':
    case 'MODIFIED':
      if (event.node) {
        nodes.update((n) => ({ ...n, [event.node!.id]: event.node! }));
      }
      if (event.edge) {
        edges.update((e) => ({ ...e, [event.edge!.id]: event.edge! }));
      }
      break;

    case 'DELETED':
      if (event.node) {
        nodes.update((n) => {
          const { [event.node!.id]: _, ...rest } = n;
          return rest;
        });
      }
      if (event.edge) {
        edges.update((e) => {
          const { [event.edge!.id]: _, ...rest } = e;
          return rest;
        });
      }
      break;
  }
}

// Derived stats
export const stats = derived([nodes, edges], ([$nodes, $edges]) => ({
  nodeCount: Object.keys($nodes).length,
  edgeCount: Object.keys($edges).length,
}));
