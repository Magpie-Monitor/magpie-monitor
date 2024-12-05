package tests

import (
	v2 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InitializeDeployments(deployments []string) []v2.Deployment {
	var res []v2.Deployment

	for _, deploy := range deployments {
		res = append(res,
			v2.Deployment{
				metav1.TypeMeta{},
				metav1.ObjectMeta{Name: deploy},
				v2.DeploymentSpec{},
				v2.DeploymentStatus{},
			})
	}

	return res
}

func InitializeStatefulSets(statefulSets []string) []v2.StatefulSet {
	var res []v2.StatefulSet

	for _, sts := range statefulSets {
		res = append(res,
			v2.StatefulSet{
				metav1.TypeMeta{},
				metav1.ObjectMeta{Name: sts},
				v2.StatefulSetSpec{},
				v2.StatefulSetStatus{},
			})
	}

	return res
}

func InitializeDaemonSets(daemonSets []string) []v2.DaemonSet {
	var res []v2.DaemonSet

	for _, ds := range daemonSets {
		res = append(res,
			v2.DaemonSet{
				metav1.TypeMeta{},
				metav1.ObjectMeta{Name: ds},
				v2.DaemonSetSpec{},
				v2.DaemonSetStatus{},
			})
	}

	return res
}
