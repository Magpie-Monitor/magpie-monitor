package pods

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path"
	"slices"
)

type Agent struct {
	kubeconfig                string
	excludedNamespaces        []string
	includedNamespaces        []string
	collectionIntervalSeconds int
	collectionDirectory       string
	client                    *kubernetes.Clientset
}

func NewAgent(kubeconfig string, excludedNamespaces []string, collectionIntervalSeconds int, collectionDirectory string) *Agent {
	return &Agent{
		kubeconfig:                kubeconfig,
		excludedNamespaces:        excludedNamespaces,
		collectionIntervalSeconds: collectionIntervalSeconds,
		collectionDirectory:       collectionDirectory,
	}
}

func (a *Agent) Start() {
	a.authenticate()
	a.fetchNamespaces()
	a.test()
	//a.prepareDirectoryTree()
	//a.gatherLogs()
}

func (a *Agent) authenticate() {
	config, err := clientcmd.BuildConfigFromFlags("", a.kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	a.client = client
}

// TODO - fix error handling
func (a *Agent) prepareDirectoryTree() {
	err := os.Mkdir(a.collectionDirectory, 0755)
	if err != nil {
		fmt.Println(err)
		//panic(err.Error())
	}

	for _, namespace := range a.includedNamespaces {
		dir := path.Join(a.collectionDirectory, namespace)
		err = os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (a *Agent) test() {
	for _, namespace := range a.includedNamespaces {
		//pods, _ := a.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		//fmt.Println(pods)
		d, _ := a.client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})

		//for d, _
		if len(d.Items) > 0 {
			selectors := d.Items[0].Spec.Selector
			//fmt.Println(selectors)

			pods, _ := a.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(selectors.MatchLabels).String()})
			fmt.Println(len(pods.Items))

			//break
		}
		//pods, _ := a.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		//d.Items[0].Spec.
		//fmt.Println(d)
		//fmt.Println(d.Items[0].Spec.Template.GetName())
		//fmt.Println(d)
	}
}

func (a *Agent) gatherLogs() {
	for _, namespace := range a.includedNamespaces {
		pods, _ := a.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		for _, pod := range pods.Items {
			f := a.getLogFileForPod(namespace, pod.Name)
			// TODO - gather logs from all containers
			logs := a.client.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{Container: pod.Spec.Containers[0].Name}).Do(context.TODO())

			res, err := logs.Raw()
			if err != nil {
				panic(err)
			}

			_, err = f.Write(res)
			if err != nil {
				panic(err)
			}

			err = f.Close()
			if err != nil {
				panic(err)
			}
		}
	}
}

func (a *Agent) getLogFileForPod(namespace string, pod string) *os.File {
	f, err := os.Create(a.collectionDirectory + "/" + namespace + "/" + pod + ".log")
	if err != nil {
		fmt.Println(err)
	}
	return f
}

// TODO - fix error handling
func (a *Agent) fetchNamespaces() {
	a.includedNamespaces = make([]string, 0)
	namespaces, _ := a.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	for _, namespace := range namespaces.Items {
		if !slices.Contains(a.excludedNamespaces, namespace.Namespace) {
			a.includedNamespaces = append(a.includedNamespaces, namespace.Name)
		}
	}
}
