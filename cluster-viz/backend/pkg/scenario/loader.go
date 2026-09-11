package scenario

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/models"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Loader loads scenario metadata from ConfigMaps or filesystem
type Loader struct {
	clientset     *kubernetes.Clientset
	namespace     string
	scenariosPath string
}

// NewLoader creates a new scenario loader
func NewLoader(clientset *kubernetes.Clientset, namespace, scenariosPath string) *Loader {
	return &Loader{
		clientset:     clientset,
		namespace:     namespace,
		scenariosPath: scenariosPath,
	}
}

// LoadFromConfigMap loads a scenario from a ConfigMap
func (l *Loader) LoadFromConfigMap(ctx context.Context, scenarioID string) (*models.Scenario, error) {
	configMapName := fmt.Sprintf("%s-architecture", scenarioID)

	cm, err := l.clientset.CoreV1().ConfigMaps(l.namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ConfigMap %s: %w", configMapName, err)
	}

	architectureYAML, ok := cm.Data["architecture.yaml"]
	if !ok {
		return nil, fmt.Errorf("ConfigMap %s missing architecture.yaml key", configMapName)
	}

	var scenario models.Scenario
	if err := yaml.Unmarshal([]byte(architectureYAML), &scenario); err != nil {
		return nil, fmt.Errorf("failed to parse architecture.yaml: %w", err)
	}

	return &scenario, nil
}

// LoadFromFilesystem loads a scenario from the local filesystem
func (l *Loader) LoadFromFilesystem(scenarioID string) (*models.Scenario, error) {
	architecturePath := filepath.Join(l.scenariosPath, fmt.Sprintf("%s-eda", scenarioID), "architecture.yaml")

	data, err := os.ReadFile(architecturePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read architecture.yaml: %w", err)
	}

	var scenario models.Scenario
	if err := yaml.Unmarshal(data, &scenario); err != nil {
		return nil, fmt.Errorf("failed to parse architecture.yaml: %w", err)
	}

	return &scenario, nil
}

// Load attempts to load from ConfigMap first, falls back to filesystem
func (l *Loader) Load(ctx context.Context, scenarioID string) (*models.Scenario, error) {
	// Try ConfigMap first (production)
	scenario, err := l.LoadFromConfigMap(ctx, scenarioID)
	if err == nil {
		return scenario, nil
	}

	// Fall back to filesystem (development)
	scenario, err = l.LoadFromFilesystem(scenarioID)
	if err != nil {
		return nil, fmt.Errorf("failed to load scenario %s from ConfigMap or filesystem: %w", scenarioID, err)
	}

	return scenario, nil
}

// ListScenariosFromConfigMaps discovers all scenarios from ConfigMaps
func (l *Loader) ListScenariosFromConfigMaps(ctx context.Context) ([]string, error) {
	configMaps, err := l.clientset.CoreV1().ConfigMaps(l.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=cluster-viz",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list ConfigMaps: %w", err)
	}

	var scenarioIDs []string
	for _, cm := range configMaps.Items {
		if _, ok := cm.Data["architecture.yaml"]; ok {
			// Extract scenario ID from ConfigMap name (e.g., "scenario-01-architecture" -> "scenario-01")
			if len(cm.Name) > len("-architecture") {
				scenarioID := cm.Name[:len(cm.Name)-len("-architecture")]
				scenarioIDs = append(scenarioIDs, scenarioID)
			}
		}
	}

	return scenarioIDs, nil
}

// ListScenariosFromFilesystem discovers all scenarios from filesystem
func (l *Loader) ListScenariosFromFilesystem() ([]string, error) {
	entries, err := os.ReadDir(l.scenariosPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenarios directory: %w", err)
	}

	var scenarioIDs []string
	for _, entry := range entries {
		if entry.IsDir() {
			architecturePath := filepath.Join(l.scenariosPath, entry.Name(), "architecture.yaml")
			if _, err := os.Stat(architecturePath); err == nil {
				// Extract scenario ID from directory name (e.g., "scenario-01-eda" -> "scenario-01")
				if len(entry.Name()) > len("-eda") {
					scenarioID := entry.Name()[:len(entry.Name())-len("-eda")]
					scenarioIDs = append(scenarioIDs, scenarioID)
				}
			}
		}
	}

	return scenarioIDs, nil
}

// ListScenarios attempts to list from ConfigMaps first, falls back to filesystem
func (l *Loader) ListScenarios(ctx context.Context) ([]string, error) {
	// Try ConfigMaps first
	scenarioIDs, err := l.ListScenariosFromConfigMaps(ctx)
	if err == nil && len(scenarioIDs) > 0 {
		return scenarioIDs, nil
	}

	// Fall back to filesystem
	scenarioIDs, err = l.ListScenariosFromFilesystem()
	if err != nil {
		return nil, fmt.Errorf("failed to list scenarios from ConfigMaps or filesystem: %w", err)
	}

	return scenarioIDs, nil
}

// CreateConfigMap creates a ConfigMap for a scenario (helper for deployment)
func (l *Loader) CreateConfigMap(ctx context.Context, scenario *models.Scenario) error {
	architectureYAML, err := yaml.Marshal(scenario)
	if err != nil {
		return fmt.Errorf("failed to marshal scenario: %w", err)
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-architecture", scenario.Metadata.ID),
			Namespace: l.namespace,
			Labels: map[string]string{
				"app":      "cluster-viz",
				"scenario": scenario.Metadata.ID,
			},
		},
		Data: map[string]string{
			"architecture.yaml": string(architectureYAML),
		},
	}

	_, err = l.clientset.CoreV1().ConfigMaps(l.namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create ConfigMap: %w", err)
	}

	return nil
}
