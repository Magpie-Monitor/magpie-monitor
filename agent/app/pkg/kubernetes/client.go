package kubernetes

import (
	"context"
	"flag"
	"log"
	"path/filepath"
	"slices"
	"time"

	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewKubernetesApiClient(runningMode string) KubernetesApiClient {
	config := getClientConfig(runningMode)

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Client{client: client}
}

func getClientConfig(runningMode string) *rest.Config {
	var config *rest.Config

	if runningMode == "local" {
		var kubeconfig *string

		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}

		c, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Println("Failed to create kubernetes API client from kubeconfig")
			panic(err.Error())
		}

		return c
	}

	c, err := rest.InClusterConfig()
	if err != nil {
		log.Println("Failed to create kubernetes API client from ServiceAccount token")
		panic(err.Error())
	}
	config = c

	return config
}

type KubernetesApiClient interface {
	GetDeployments(namespace string) ([]v2.Deployment, error)
	GetStatefulSets(namespace string) ([]v2.StatefulSet, error)
	GetDaemonSets(namespace string) ([]v2.DaemonSet, error)
	GetNamespaces(excludedNamespaces []string) ([]string, error)
	GetPods(selector *metav1.LabelSelector, namespace string) ([]v1.Pod, error)
	GetContainerLogsSinceTime(podName, containerName, namespace string, sinceTime time.Time, timestamps bool) (string, error)
}

type Client struct {
	client *kubernetes.Clientset
}

func (c *Client) GetDeployments(namespace string) ([]v2.Deployment, error) {
	deployments, err := c.client.AppsV1().
		Deployments(namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return deployments.Items, nil
}

func (c *Client) GetStatefulSets(namespace string) ([]v2.StatefulSet, error) {
	statefulSets, err := c.client.AppsV1().
		StatefulSets(namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return statefulSets.Items, nil
}

func (c *Client) GetDaemonSets(namespace string) ([]v2.DaemonSet, error) {
	daemonSets, err := c.client.AppsV1().
		DaemonSets(namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return daemonSets.Items, nil
}

func (c *Client) GetNamespaces(excludedNamespaces []string) ([]string, error) {
	var includedNamespaces []string

	namespaces, err := c.client.CoreV1().
		Namespaces().
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces.Items {
		if !slices.Contains(excludedNamespaces, namespace.Name) {
			includedNamespaces = append(includedNamespaces, namespace.Name)
		}
	}

	return includedNamespaces, nil
}

func (c *Client) GetPods(selector *metav1.LabelSelector, namespace string) ([]v1.Pod, error) {
	pods, err := c.client.CoreV1().
		Pods(namespace).
		List(
			context.TODO(),
			metav1.ListOptions{LabelSelector: labels.Set(selector.MatchLabels).String()},
		)

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func (c *Client) GetContainerLogsSinceTime(podName, containerName, namespace string, sinceTime time.Time, timestamps bool) (string, error) {
	logsByte, err := c.client.CoreV1().
		Pods(namespace).
		GetLogs(
			podName,
			&v1.PodLogOptions{
				Container:  containerName,
				SinceTime:  &metav1.Time{Time: sinceTime},
				Timestamps: timestamps,
			},
		).
		Do(context.TODO()).
		Raw()

	if err != nil {
		return "", err
	}

	return string(logsByte), nil
}
