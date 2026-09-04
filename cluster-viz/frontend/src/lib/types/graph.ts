// TypeScript mirrors of Go models

export interface GraphNode {
  id: string;
  kind: string;
  name: string;
  namespace: string;
  status: string;
  labels: Record<string, string>;
  parentId?: string;
  metadata: Record<string, string>;
}

export interface GraphEdge {
  id: string;
  source: string;
  target: string;
  kind: 'ownership' | 'routing' | 'selection' | 'dataflow' | 'parentRef';
  label?: string;
}

export interface Graph {
  nodes: Record<string, GraphNode>;
  edges: Record<string, GraphEdge>;
}

export type EventType = 'SYNC' | 'ADDED' | 'MODIFIED' | 'DELETED';

export interface GraphEvent {
  type: EventType;
  node?: GraphNode;
  edge?: GraphEdge;
  graph?: Graph;
}
