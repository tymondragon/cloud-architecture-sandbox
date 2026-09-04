package k8s

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/models"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	// GatewayGVR is the GroupVersionResource for Gateway API Gateway
	GatewayGVR = schema.GroupVersionResource{
		Group:    "gateway.networking.k8s.io",
		Version:  "v1",
		Resource: "gateways",
	}

	// HTTPRouteGVR is the GroupVersionResource for Gateway API HTTPRoute
	HTTPRouteGVR = schema.GroupVersionResource{
		Group:    "gateway.networking.k8s.io",
		Version:  "v1",
		Resource: "httproutes",
	}
)

// ResourceWatcher watches Kubernetes resources and emits graph events
type ResourceWatcher struct {
	clientset     *kubernetes.Clientset
	dynamicClient dynamic.Interface
	graphBuilder  *GraphBuilder
	events        chan models.GraphEvent
	stopCh        chan struct{}
	excludeNs     map[string]bool
}

// NewResourceWatcher creates a new resource watcher
func NewResourceWatcher(clientset *kubernetes.Clientset, dynamicClient dynamic.Interface, graphBuilder *GraphBuilder) *ResourceWatcher {
	// Parse excluded namespaces from env var
	excludeNs := make(map[string]bool)
	excludeList := os.Getenv("EXCLUDE_NAMESPACES")
	if excludeList == "" {
		excludeList = "kube-system,kube-public,kube-node-lease,local-path-storage"
	}
	for _, ns := range strings.Split(excludeList, ",") {
		excludeNs[strings.TrimSpace(ns)] = true
	}

	return &ResourceWatcher{
		clientset:     clientset,
		dynamicClient: dynamicClient,
		graphBuilder:  graphBuilder,
		events:        make(chan models.GraphEvent, 100),
		stopCh:        make(chan struct{}),
		excludeNs:     excludeNs,
	}
}

// Events returns the event channel
func (rw *ResourceWatcher) Events() <-chan models.GraphEvent {
	return rw.events
}

// Start begins watching all resource types
func (rw *ResourceWatcher) Start(ctx context.Context) {
	log.Println("Starting resource watchers...")

	// Watch core resources
	go rw.watchNamespaces(ctx)
	go rw.watchPods(ctx)
	go rw.watchServices(ctx)
	go rw.watchEndpoints(ctx)
	go rw.watchDeployments(ctx)
	go rw.watchStatefulSets(ctx)
	go rw.watchDaemonSets(ctx)

	// Watch Gateway API resources
	go rw.watchGateways(ctx)
	go rw.watchHTTPRoutes(ctx)

	log.Println("All watchers started")
}

// Stop stops all watchers
func (rw *ResourceWatcher) Stop() {
	close(rw.stopCh)
	close(rw.events)
}

// shouldExcludeNamespace checks if a namespace should be excluded
func (rw *ResourceWatcher) shouldExcludeNamespace(ns string) bool {
	return rw.excludeNs[ns]
}

// watchNamespaces watches Namespace resources
func (rw *ResourceWatcher) watchNamespaces(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.clientset.CoreV1().Namespaces().List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.clientset.CoreV1().Namespaces().Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &corev1.Namespace{}, 10*time.Minute)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns := obj.(*corev1.Namespace)
			if rw.shouldExcludeNamespace(ns.Name) {
				return
			}
			node := rw.namespaceToNode(ns)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ns := newObj.(*corev1.Namespace)
			if rw.shouldExcludeNamespace(ns.Name) {
				return
			}
			node := rw.namespaceToNode(ns)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}
		},
		DeleteFunc: func(obj interface{}) {
			ns := obj.(*corev1.Namespace)
			nodeID := BuildNodeID("Namespace", "", ns.Name)
			rw.graphBuilder.DeleteNode(nodeID)
			rw.events <- models.GraphEvent{Type: models.EventTypeDeleted, Node: &models.GraphNode{ID: nodeID}}
		},
	})

	go informer.Run(rw.stopCh)
}

// watchPods watches Pod resources
func (rw *ResourceWatcher) watchPods(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.clientset.CoreV1().Pods(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.clientset.CoreV1().Pods(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &corev1.Pod{}, 10*time.Minute)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			if rw.shouldExcludeNamespace(pod.Namespace) {
				return
			}
			node := rw.podToNode(pod)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			// Rebuild edges after pod changes (affects selection and dataflow)
			go rw.rebuildEdgesAsync()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*corev1.Pod)
			if rw.shouldExcludeNamespace(pod.Namespace) {
				return
			}
			node := rw.podToNode(pod)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			nodeID := BuildNodeID("Pod", pod.Namespace, pod.Name)
			rw.graphBuilder.DeleteNode(nodeID)
			rw.events <- models.GraphEvent{Type: models.EventTypeDeleted, Node: &models.GraphNode{ID: nodeID}}

			go rw.rebuildEdgesAsync()
		},
	})

	go informer.Run(rw.stopCh)
}

// watchServices watches Service resources
func (rw *ResourceWatcher) watchServices(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.clientset.CoreV1().Services(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.clientset.CoreV1().Services(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &corev1.Service{}, 10*time.Minute)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svc := obj.(*corev1.Service)
			if rw.shouldExcludeNamespace(svc.Namespace) {
				return
			}
			node := rw.serviceToNode(svc)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			svc := newObj.(*corev1.Service)
			if rw.shouldExcludeNamespace(svc.Namespace) {
				return
			}
			node := rw.serviceToNode(svc)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		DeleteFunc: func(obj interface{}) {
			svc := obj.(*corev1.Service)
			nodeID := BuildNodeID("Service", svc.Namespace, svc.Name)
			rw.graphBuilder.DeleteNode(nodeID)
			rw.events <- models.GraphEvent{Type: models.EventTypeDeleted, Node: &models.GraphNode{ID: nodeID}}

			go rw.rebuildEdgesAsync()
		},
	})

	go informer.Run(rw.stopCh)
}

// watchEndpoints watches Endpoints resources (for future use)
func (rw *ResourceWatcher) watchEndpoints(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.clientset.CoreV1().Endpoints(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.clientset.CoreV1().Endpoints(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &corev1.Endpoints{}, 10*time.Minute)
	go informer.Run(rw.stopCh)
}

// watchDeployments watches Deployment resources
func (rw *ResourceWatcher) watchDeployments(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.clientset.AppsV1().Deployments(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.clientset.AppsV1().Deployments(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &appsv1.Deployment{}, 10*time.Minute)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			if rw.shouldExcludeNamespace(deploy.Namespace) {
				return
			}
			node := rw.deploymentToNode(deploy)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			deploy := newObj.(*appsv1.Deployment)
			if rw.shouldExcludeNamespace(deploy.Namespace) {
				return
			}
			node := rw.deploymentToNode(deploy)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		DeleteFunc: func(obj interface{}) {
			deploy := obj.(*appsv1.Deployment)
			nodeID := BuildNodeID("Deployment", deploy.Namespace, deploy.Name)
			rw.graphBuilder.DeleteNode(nodeID)
			rw.events <- models.GraphEvent{Type: models.EventTypeDeleted, Node: &models.GraphNode{ID: nodeID}}

			go rw.rebuildEdgesAsync()
		},
	})

	go informer.Run(rw.stopCh)
}

// watchStatefulSets watches StatefulSet resources
func (rw *ResourceWatcher) watchStatefulSets(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.clientset.AppsV1().StatefulSets(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.clientset.AppsV1().StatefulSets(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &appsv1.StatefulSet{}, 10*time.Minute)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			sts := obj.(*appsv1.StatefulSet)
			if rw.shouldExcludeNamespace(sts.Namespace) {
				return
			}
			node := rw.statefulSetToNode(sts)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			sts := newObj.(*appsv1.StatefulSet)
			if rw.shouldExcludeNamespace(sts.Namespace) {
				return
			}
			node := rw.statefulSetToNode(sts)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		DeleteFunc: func(obj interface{}) {
			sts := obj.(*appsv1.StatefulSet)
			nodeID := BuildNodeID("StatefulSet", sts.Namespace, sts.Name)
			rw.graphBuilder.DeleteNode(nodeID)
			rw.events <- models.GraphEvent{Type: models.EventTypeDeleted, Node: &models.GraphNode{ID: nodeID}}

			go rw.rebuildEdgesAsync()
		},
	})

	go informer.Run(rw.stopCh)
}

// watchDaemonSets watches DaemonSet resources
func (rw *ResourceWatcher) watchDaemonSets(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.clientset.AppsV1().DaemonSets(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.clientset.AppsV1().DaemonSets(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &appsv1.DaemonSet{}, 10*time.Minute)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ds := obj.(*appsv1.DaemonSet)
			if rw.shouldExcludeNamespace(ds.Namespace) {
				return
			}
			node := rw.daemonSetToNode(ds)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			ds := newObj.(*appsv1.DaemonSet)
			if rw.shouldExcludeNamespace(ds.Namespace) {
				return
			}
			node := rw.daemonSetToNode(ds)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		DeleteFunc: func(obj interface{}) {
			ds := obj.(*appsv1.DaemonSet)
			nodeID := BuildNodeID("DaemonSet", ds.Namespace, ds.Name)
			rw.graphBuilder.DeleteNode(nodeID)
			rw.events <- models.GraphEvent{Type: models.EventTypeDeleted, Node: &models.GraphNode{ID: nodeID}}

			go rw.rebuildEdgesAsync()
		},
	})

	go informer.Run(rw.stopCh)
}

// watchGateways watches Gateway API Gateway resources
func (rw *ResourceWatcher) watchGateways(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.dynamicClient.Resource(GatewayGVR).Namespace(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.dynamicClient.Resource(GatewayGVR).Namespace(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &unstructured.Unstructured{}, 10*time.Minute)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			if rw.shouldExcludeNamespace(u.GetNamespace()) {
				return
			}
			node := rw.gatewayToNode(u)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			u := newObj.(*unstructured.Unstructured)
			if rw.shouldExcludeNamespace(u.GetNamespace()) {
				return
			}
			node := rw.gatewayToNode(u)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		DeleteFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			nodeID := BuildNodeID("Gateway", u.GetNamespace(), u.GetName())
			rw.graphBuilder.DeleteNode(nodeID)
			rw.events <- models.GraphEvent{Type: models.EventTypeDeleted, Node: &models.GraphNode{ID: nodeID}}

			go rw.rebuildEdgesAsync()
		},
	})

	go informer.Run(rw.stopCh)
}

// watchHTTPRoutes watches Gateway API HTTPRoute resources
func (rw *ResourceWatcher) watchHTTPRoutes(ctx context.Context) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return rw.dynamicClient.Resource(HTTPRouteGVR).Namespace(corev1.NamespaceAll).List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return rw.dynamicClient.Resource(HTTPRouteGVR).Namespace(corev1.NamespaceAll).Watch(ctx, options)
		},
	}

	informer := cache.NewSharedInformer(lw, &unstructured.Unstructured{}, 10*time.Minute)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			if rw.shouldExcludeNamespace(u.GetNamespace()) {
				return
			}
			node := rw.httpRouteToNode(u)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			u := newObj.(*unstructured.Unstructured)
			if rw.shouldExcludeNamespace(u.GetNamespace()) {
				return
			}
			node := rw.httpRouteToNode(u)
			eventType := rw.graphBuilder.AddOrUpdateNode(node)
			rw.events <- models.GraphEvent{Type: eventType, Node: node}

			go rw.rebuildEdgesAsync()
		},
		DeleteFunc: func(obj interface{}) {
			u := obj.(*unstructured.Unstructured)
			nodeID := BuildNodeID("HTTPRoute", u.GetNamespace(), u.GetName())
			rw.graphBuilder.DeleteNode(nodeID)
			rw.events <- models.GraphEvent{Type: models.EventTypeDeleted, Node: &models.GraphNode{ID: nodeID}}

			go rw.rebuildEdgesAsync()
		},
	})

	go informer.Run(rw.stopCh)
}

// rebuildEdgesAsync debounces edge rebuilding
var rebuildTimer *time.Timer
var rebuildMutex sync.Mutex

func (rw *ResourceWatcher) rebuildEdgesAsync() {
	rebuildMutex.Lock()
	defer rebuildMutex.Unlock()

	// Debounce: only rebuild after 500ms of no changes
	if rebuildTimer != nil {
		rebuildTimer.Stop()
	}

	rebuildTimer = time.AfterFunc(500*time.Millisecond, func() {
		edges := rw.graphBuilder.RebuildEdges()
		// Emit edge events (simplified: just send as modified)
		for _, edge := range edges {
			rw.events <- models.GraphEvent{Type: models.EventTypeModified, Edge: edge}
		}
	})
}

// Converter functions

func (rw *ResourceWatcher) namespaceToNode(ns *corev1.Namespace) *models.GraphNode {
	return &models.GraphNode{
		ID:        BuildNodeID("Namespace", "", ns.Name),
		Kind:      "Namespace",
		Name:      ns.Name,
		Namespace: "",
		Status:    string(ns.Status.Phase),
		Labels:    ns.Labels,
		Metadata:  make(map[string]string),
	}
}

func (rw *ResourceWatcher) podToNode(pod *corev1.Pod) *models.GraphNode {
	status := string(pod.Status.Phase)

	// Check container statuses for more detail
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Waiting != nil {
			status = "Pending"
			break
		}
		if cs.State.Terminated != nil {
			status = "Terminated"
			break
		}
		if !cs.Ready {
			status = "Error"
		}
	}

	metadata := make(map[string]string)

	// Store owner references
	if len(pod.OwnerReferences) > 0 {
		var ownerIDs []string
		for _, owner := range pod.OwnerReferences {
			ownerIDs = append(ownerIDs, BuildOwnerID(pod.Namespace, owner))
		}
		metadata["ownerReferences"] = strings.Join(ownerIDs, ",")
	}

	// Extract env vars for dataflow edge detection
	metadata["envVars"] = ExtractEnvVars(pod.Spec.Containers)

	return &models.GraphNode{
		ID:        BuildNodeID("Pod", pod.Namespace, pod.Name),
		Kind:      "Pod",
		Name:      pod.Name,
		Namespace: pod.Namespace,
		Status:    status,
		Labels:    pod.Labels,
		ParentID:  BuildNodeID("Namespace", "", pod.Namespace),
		Metadata:  metadata,
	}
}

func (rw *ResourceWatcher) serviceToNode(svc *corev1.Service) *models.GraphNode {
	metadata := make(map[string]string)

	// Store selector for selection edge detection
	if len(svc.Spec.Selector) > 0 {
		var pairs []string
		for k, v := range svc.Spec.Selector {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}
		metadata["selector"] = strings.Join(pairs, ",")
	}

	return &models.GraphNode{
		ID:        BuildNodeID("Service", svc.Namespace, svc.Name),
		Kind:      "Service",
		Name:      svc.Name,
		Namespace: svc.Namespace,
		Status:    "Active",
		Labels:    svc.Labels,
		ParentID:  BuildNodeID("Namespace", "", svc.Namespace),
		Metadata:  metadata,
	}
}

func (rw *ResourceWatcher) deploymentToNode(deploy *appsv1.Deployment) *models.GraphNode {
	status := "Healthy"
	if deploy.Status.ReadyReplicas < deploy.Status.Replicas {
		status = "Degraded"
	}
	if deploy.Status.ReadyReplicas == 0 {
		status = "Unavailable"
	}

	metadata := map[string]string{
		"replicas":      fmt.Sprintf("%d/%d", deploy.Status.ReadyReplicas, deploy.Status.Replicas),
	}

	return &models.GraphNode{
		ID:        BuildNodeID("Deployment", deploy.Namespace, deploy.Name),
		Kind:      "Deployment",
		Name:      deploy.Name,
		Namespace: deploy.Namespace,
		Status:    status,
		Labels:    deploy.Labels,
		ParentID:  BuildNodeID("Namespace", "", deploy.Namespace),
		Metadata:  metadata,
	}
}

func (rw *ResourceWatcher) statefulSetToNode(sts *appsv1.StatefulSet) *models.GraphNode {
	status := "Healthy"
	if sts.Status.ReadyReplicas < sts.Status.Replicas {
		status = "Degraded"
	}
	if sts.Status.ReadyReplicas == 0 {
		status = "Unavailable"
	}

	metadata := map[string]string{
		"replicas": fmt.Sprintf("%d/%d", sts.Status.ReadyReplicas, sts.Status.Replicas),
	}

	return &models.GraphNode{
		ID:        BuildNodeID("StatefulSet", sts.Namespace, sts.Name),
		Kind:      "StatefulSet",
		Name:      sts.Name,
		Namespace: sts.Namespace,
		Status:    status,
		Labels:    sts.Labels,
		ParentID:  BuildNodeID("Namespace", "", sts.Namespace),
		Metadata:  metadata,
	}
}

func (rw *ResourceWatcher) daemonSetToNode(ds *appsv1.DaemonSet) *models.GraphNode {
	status := "Healthy"
	if ds.Status.NumberReady < ds.Status.DesiredNumberScheduled {
		status = "Degraded"
	}
	if ds.Status.NumberReady == 0 {
		status = "Unavailable"
	}

	metadata := map[string]string{
		"replicas": fmt.Sprintf("%d/%d", ds.Status.NumberReady, ds.Status.DesiredNumberScheduled),
	}

	return &models.GraphNode{
		ID:        BuildNodeID("DaemonSet", ds.Namespace, ds.Name),
		Kind:      "DaemonSet",
		Name:      ds.Name,
		Namespace: ds.Namespace,
		Status:    status,
		Labels:    ds.Labels,
		ParentID:  BuildNodeID("Namespace", "", ds.Namespace),
		Metadata:  metadata,
	}
}

func (rw *ResourceWatcher) gatewayToNode(u *unstructured.Unstructured) *models.GraphNode {
	status := "Unknown"

	// Extract status conditions
	if conditions, found, _ := unstructured.NestedSlice(u.Object, "status", "conditions"); found {
		for _, c := range conditions {
			cond := c.(map[string]interface{})
			if cond["type"] == "Programmed" && cond["status"] == "True" {
				status = "Ready"
				break
			}
		}
	}

	return &models.GraphNode{
		ID:        BuildNodeID("Gateway", u.GetNamespace(), u.GetName()),
		Kind:      "Gateway",
		Name:      u.GetName(),
		Namespace: u.GetNamespace(),
		Status:    status,
		Labels:    u.GetLabels(),
		ParentID:  BuildNodeID("Namespace", "", u.GetNamespace()),
		Metadata:  make(map[string]string),
	}
}

func (rw *ResourceWatcher) httpRouteToNode(u *unstructured.Unstructured) *models.GraphNode {
	status := "Unknown"
	metadata := make(map[string]string)

	// Extract parentRefs
	if parentRefs, found, _ := unstructured.NestedSlice(u.Object, "spec", "parentRefs"); found {
		var parentIDs []string
		for _, p := range parentRefs {
			parent := p.(map[string]interface{})
			name := parent["name"].(string)
			namespace := u.GetNamespace()
			if ns, ok := parent["namespace"].(string); ok {
				namespace = ns
			}
			parentIDs = append(parentIDs, BuildNodeID("Gateway", namespace, name))
		}
		metadata["parentRefs"] = strings.Join(parentIDs, ",")
	}

	// Extract backendRefs
	if rules, found, _ := unstructured.NestedSlice(u.Object, "spec", "rules"); found {
		var backendIDs []string
		for _, r := range rules {
			rule := r.(map[string]interface{})
			if backendRefs, ok := rule["backendRefs"].([]interface{}); ok {
				for _, b := range backendRefs {
					backend := b.(map[string]interface{})
					name := backend["name"].(string)
					namespace := u.GetNamespace()
					if ns, ok := backend["namespace"].(string); ok {
						namespace = ns
					}
					backendIDs = append(backendIDs, BuildNodeID("Service", namespace, name))
				}
			}
		}
		metadata["backendRefs"] = strings.Join(backendIDs, ",")
	}

	// Extract status conditions
	if conditions, found, _ := unstructured.NestedSlice(u.Object, "status", "parents"); found && len(conditions) > 0 {
		parent := conditions[0].(map[string]interface{})
		if conds, ok := parent["conditions"].([]interface{}); ok {
			for _, c := range conds {
				cond := c.(map[string]interface{})
				if cond["type"] == "Accepted" && cond["status"] == "True" {
					status = "Accepted"
					break
				}
			}
		}
	}

	return &models.GraphNode{
		ID:        BuildNodeID("HTTPRoute", u.GetNamespace(), u.GetName()),
		Kind:      "HTTPRoute",
		Name:      u.GetName(),
		Namespace: u.GetNamespace(),
		Status:    status,
		Labels:    u.GetLabels(),
		ParentID:  BuildNodeID("Namespace", "", u.GetNamespace()),
		Metadata:  metadata,
	}
}
