package k8s

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClients creates both typed and dynamic Kubernetes clients
// It tries in-cluster config first, then falls back to kubeconfig for local dev
func NewClients() (*kubernetes.Clientset, dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// Not running in cluster, try kubeconfig
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, nil, err
			}
			kubeconfig = filepath.Join(home, ".kube", "config")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, err
		}
	}

	// Create typed client for core resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	// Create dynamic client for CRDs (Gateway API)
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return clientset, dynamicClient, nil
}
