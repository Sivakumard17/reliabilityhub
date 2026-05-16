package k8sclient

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps the Kubernetes client
type Client struct {
	cs  kubernetes.Interface
	log *zap.Logger
}

// New creates a Kubernetes client.
// Uses in-cluster config when running in a pod,
// falls back to kubeconfig for local development.
func New(kubeconfigPath string, log *zap.Logger) (*Client, error) {
	var cfg *rest.Config
	var err error

	if kubeconfigPath != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			// Fall back to default kubeconfig
			cfg, err = clientcmd.BuildConfigFromFlags("", "")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("build k8s config: %w", err)
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create k8s client: %w", err)
	}

	log.Info("kubernetes client initialized")
	return &Client{cs: cs, log: log}, nil
}

// RestartDeployment performs a rolling restart by patching an annotation.
// This is equivalent to: kubectl rollout restart deployment/<name>
func (c *Client) RestartDeployment(ctx context.Context, namespace, name string) error {
	patch := fmt.Sprintf(
		`{"spec":{"template":{"metadata":{"annotations":{"reliabilityhub.io/restartedAt":"%s"}}}}}`,
		time.Now().UTC().Format(time.RFC3339),
	)

	_, err := c.cs.AppsV1().Deployments(namespace).Patch(
		ctx,
		name,
		types.MergePatchType,
		[]byte(patch),
		metav1.PatchOptions{},
	)
	if err != nil {
		return fmt.Errorf("restart deployment %s/%s: %w", namespace, name, err)
	}

	c.log.Info("deployment restarted",
		zap.String("namespace", namespace),
		zap.String("deployment", name),
	)
	return nil
}

// ScaleDeployment sets the replica count for a deployment.
func (c *Client) ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error {
	scale, err := c.cs.AppsV1().Deployments(namespace).GetScale(
		ctx, name, metav1.GetOptions{},
	)
	if err != nil {
		return fmt.Errorf("get scale %s/%s: %w", namespace, name, err)
	}

	scale.Spec.Replicas = replicas
	_, err = c.cs.AppsV1().Deployments(namespace).UpdateScale(
		ctx, name, scale, metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("scale deployment %s/%s: %w", namespace, name, err)
	}

	c.log.Info("deployment scaled",
		zap.String("namespace", namespace),
		zap.String("deployment", name),
		zap.Int32("replicas", replicas),
	)
	return nil
}

// GetDeployment returns a deployment by name.
func (c *Client) GetDeployment(ctx context.Context, namespace, name string) (*appsv1.Deployment, error) {
	return c.cs.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

// ListDeployments returns all deployments in a namespace.
func (c *Client) ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	list, err := c.cs.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// RollbackDeployment rolls back to the previous revision.
func (c *Client) RollbackDeployment(ctx context.Context, namespace, name string) error {
	// Get current deployment
	deploy, err := c.GetDeployment(ctx, namespace, name)
	if err != nil {
		return fmt.Errorf("get deployment: %w", err)
	}

	// Annotate to trigger rollback via undo annotation
	if deploy.Annotations == nil {
		deploy.Annotations = map[string]string{}
	}
	deploy.Annotations["reliabilityhub.io/rollback"] = time.Now().UTC().Format(time.RFC3339)

	// Decrement revision to trigger rollback
	patch := `{"spec":{"template":{"metadata":{"annotations":{"reliabilityhub.io/rollback":"` +
		time.Now().UTC().Format(time.RFC3339) + `"}}}}}`

	_, err = c.cs.AppsV1().Deployments(namespace).Patch(
		ctx, name, types.MergePatchType, []byte(patch), metav1.PatchOptions{},
	)
	if err != nil {
		return fmt.Errorf("rollback deployment %s/%s: %w", namespace, name, err)
	}

	c.log.Info("deployment rollback triggered",
		zap.String("namespace", namespace),
		zap.String("deployment", name),
	)
	return nil
}

// Ping verifies the k8s API is reachable
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cs.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	return err
}
