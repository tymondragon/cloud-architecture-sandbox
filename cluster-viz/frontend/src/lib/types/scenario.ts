// Scenario metadata types matching the backend models

export interface Scenario {
  metadata: ScenarioMetadata;
  patterns: ScenarioPatterns;
  technologies: Technology[];
  layout: LayoutConfig;
  layers: Layer[];
  components: Component[];
  flows: Flow[];
  telemetry: TelemetryConfig;
  dataFlows: DataFlow[];
  testDataGenerator: TestDataGenerator;
}

export interface ScenarioMetadata {
  id: string;
  title: string;
  description: string;
}

export interface ScenarioPatterns {
  primary: Pattern[];
  secondary: Pattern[];
}

export interface Pattern {
  name: string;
  description: string;
  benefits: string[];
  color: string;
  icon: string;
}

export interface Technology {
  name: string;
  role: string;
  description: string;
}

export interface LayoutConfig {
  type: 'vertical' | 'horizontal' | 'radial' | 'force-graph' | 'grid' | 'custom';
  options: LayoutOptions;
}

export interface LayoutOptions {
  spacing?: string;
  fitToViewport?: boolean;
  compactCards?: boolean;
  centerComponent?: string;
  columns?: number;
  positions?: Record<string, { x: number; y: number }>;
}

export interface Layer {
  id: string;
  name: string;
  description: string;
  color: string;
  position: number;
}

export interface Component {
  id: string;
  name: string;
  type: string;
  layer: string;
  role: string;
  description: string;
  icon: string;
  implementation: ComponentImplementation;
}

export interface ComponentImplementation {
  kubernetes: K8sImplementation;
  abstractName?: string;
  actualTechnology?: string;
}

export interface K8sImplementation {
  deployment?: string;
  statefulset?: string;
  service?: string;
  cluster?: string;
  namespace: string;
}

export interface Flow {
  from: string;
  to: string;
  type: string;
  label: string;
  description: string;
  pattern: string;
  edgeStyle: EdgeStyle;
}

export interface EdgeStyle {
  color: string;
  style: string;
  animated: boolean;
}

export interface TelemetryConfig {
  collectors: CollectorConfig[];
}

export interface CollectorConfig {
  type: string;
  component: string;
  config: Record<string, any>;
  metrics: string[];
}

export interface DataFlow {
  id: string;
  name: string;
  description: string;
  steps: DataFlowStep[];
}

export interface DataFlowStep {
  component: string;
  action: string;
  description: string;
  duration: number;
  telemetry?: TelemetryRef;
  parallel?: boolean;
}

export interface TelemetryRef {
  collector: string;
  metric: string;
}

export interface TestDataGenerator {
  endpoint: string;
  method: string;
  samplePayloads: SamplePayload[];
}

export interface SamplePayload {
  name: string;
  description: string;
  data: Record<string, any>;
  expectedFlow: string;
}

export interface ScenarioSummary {
  id: string;
  title: string;
  description: string;
  patterns: string[];
}

// Logical graph types (response from /api/graph)
export interface LogicalGraph {
  scenario: ScenarioMetadata;
  patterns: ScenarioPatterns;
  layers: Layer[];
  components: Record<string, GraphNode>;
  edges: Record<string, GraphEdge>;
  health: Record<string, ComponentHealth>;
}

export interface GraphNode {
  id: string;
  kind: string;
  name: string;
  namespace: string;
  status: string;
  labels: Record<string, string>;
  parentId: string;
  metadata: Record<string, string>;
}

export interface GraphEdge {
  id: string;
  source: string;
  target: string;
  kind: string;
  label: string;
  metadata: Record<string, string>;
}

export interface ComponentHealth {
  componentId: string;
  status: string; // "healthy", "degraded", "unhealthy", "unknown"
  replicas?: string;
  message?: string;
}
