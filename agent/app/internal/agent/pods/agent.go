package pods

import (
	"context"
	"fmt"
	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"slices"
	"time"
)

type Chunk struct {
	kind      string
	name      string
	namespace string
	content   string
}

type Agent struct {
	kubeconfig                string
	excludedNamespaces        []string
	includedNamespaces        []string
	collectionIntervalSeconds int
	collectionDirectory       string
	client                    *kubernetes.Clientset
	readTimestamps            map[string]int64
}

func NewAgent(kubeconfig string, excludedNamespaces []string, collectionIntervalSeconds int, collectionDirectory string) *Agent {
	return &Agent{
		kubeconfig:                kubeconfig,
		excludedNamespaces:        excludedNamespaces,
		collectionIntervalSeconds: collectionIntervalSeconds,
		collectionDirectory:       collectionDirectory,
		readTimestamps:            make(map[string]int64),
	}
}

func (a *Agent) Start() {
	a.authenticate()
	a.fetchNamespaces()
	a.gatherLogs(10)
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

func (a *Agent) fetchNamespaces() {
	a.includedNamespaces = make([]string, 0)
	namespaces, err := a.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(fmt.Sprintf("Error fetching namespaces: %s", err.Error()))
	}

	for _, namespace := range namespaces.Items {
		if !slices.Contains(a.excludedNamespaces, namespace.Namespace) {
			a.includedNamespaces = append(a.includedNamespaces, namespace.Name)
		}
	}
}

func (a *Agent) gatherLogs(scrapeInterval int) {
	for {
		for _, namespace := range a.includedNamespaces {
			log.Println("Fetching logs for namespace: ", namespace)

			deployments, err := a.client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching Deployments: ", err)
				log.Println("Skipping iteration")
			} else {
				a.fetchDeploymentLogsSinceSeconds(namespace, deployments.Items)
			}

			statefulSets, err := a.client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching StatefulSets: ", err)
				log.Println("Skipping iteration")
			} else {
				a.fetchStatefulSetLogsSinceSeconds(namespace, statefulSets.Items)
			}

			daemonSets, err := a.client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Println("Error fetching DaemonSets: ", err)
				log.Println("Skipping iteration")
			} else {
				a.fetchDaemonSetLogsSinceSeconds(namespace, daemonSets.Items)
			}
		}

		log.Println("Sleeping for: ", scrapeInterval, " seconds")
		time.Sleep(time.Duration(scrapeInterval) * time.Second)
	}
}

func (a *Agent) calculateLastReadTimeDiff(name string) int64 {
	var readDiff int64 = 10
	now := time.Now()
	unix := now.Unix()
	previousRead, ok := a.readTimestamps[name]
	if ok {
		readDiff = unix - previousRead
	}
	return readDiff
}

func (a *Agent) setReadTimestamp(name string) {
	now := time.Now()
	a.readTimestamps[name] = now.Unix()
}

func (a *Agent) fetchDeploymentLogsSinceSeconds(namespace string, deployments []v2.Deployment) {
	for _, deployment := range deployments {
		name := deployment.Name
		selectors := deployment.Spec.Selector

		log.Println("Deployment: ", name)

		readDiff := a.calculateLastReadTimeDiff(name)
		a.fetchLogsSinceSeconds(selectors, namespace, &readDiff)
		a.setReadTimestamp(name)
	}
}

func (a *Agent) fetchStatefulSetLogsSinceSeconds(namespace string, statefulSets []v2.StatefulSet) {
	log.Println("Fetching logs from StatefulSets")
	for _, statefulSet := range statefulSets {
		name := statefulSet.Name
		selectors := statefulSet.Spec.Selector

		log.Println("StatefulSet: ", name)
		readDiff := a.calculateLastReadTimeDiff(name)
		a.fetchLogsSinceSeconds(selectors, namespace, &readDiff)
		a.setReadTimestamp(name)
	}
}

func (a *Agent) fetchDaemonSetLogsSinceSeconds(namespace string, daemonSets []v2.DaemonSet) {
	log.Println("Fetching logs from DaemonSets")
	for _, daemonSet := range daemonSets {
		name := daemonSet.Name
		selectors := daemonSet.Spec.Selector

		log.Println("DaemonSet: ", name)
		readDiff := a.calculateLastReadTimeDiff(name)
		a.fetchLogsSinceSeconds(selectors, namespace, &readDiff)
		a.setReadTimestamp(name)
	}
}

func (a *Agent) fetchLogsSinceSeconds(selector *metav1.LabelSelector, namespace string, sinceSeconds *int64) {
	pods, _ := a.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labels.Set(selector.MatchLabels).String()})
	for _, pod := range pods.Items {
		log.Println("pod: ", pod.Name)

		// todo - all containers
		logs := a.client.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{Container: pod.Spec.Containers[0].Name, SinceSeconds: sinceSeconds}).Do(context.TODO())

		l, _ := logs.Raw()
		log.Println("read log chunk: ", string(l))
	}
}

// TODO - leaving for reference
//	func (a *Agent) gatherLogs() {
//		for _, namespace := range a.includedNamespaces {
//			pods, _ := a.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
//			for _, pod := range pods.Items {
//				f := a.getLogFileForPod(namespace, pod.Name)
//				// TODO - gather logs from all containers
//				logs := a.client.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{Container: pod.Spec.Containers[0].Name}).Do(context.TODO())
//
//				res, err := logs.Raw()
//				if err != nil {
//					panic(err)
//				}
//
//				_, err = f.Write(res)
//				if err != nil {
//					panic(err)
//				}
//
//				err = f.Close()
//				if err != nil {
//					panic(err)
//				}
//			}
//		}
//	}
//
//	func (a *Agent) getLogFileForPod(namespace string, pod string) *os.File {
//		f, err := os.Create(a.collectionDirectory + "/" + namespace + "/" + pod + ".log")
//		if err != nil {
//			fmt.Println(err)
//		}
//		return f
//	}

//// TODO - fix error handling
//func (a *Agent) prepareDirectoryTree() {
//	err := os.Mkdir(a.collectionDirectory, 0755)
//	if err != nil {
//		fmt.Println(err)
//		//panic(err.Error())
//	}
//
//	for _, namespace := range a.includedNamespaces {
//		dir := path.Join(a.collectionDirectory, namespace)
//		err = os.Mkdir(dir, 0755)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//}
