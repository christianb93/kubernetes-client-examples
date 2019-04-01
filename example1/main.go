package main

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := "/home/chr/.kube/config"
	fmt.Printf("%-20s  %-10s %s\n", "NAME", "ARCH", "OS")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	coreClient := clientset.CoreV1()
	nodeList, err := coreClient.Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	items := nodeList.Items
	for _, item := range items {
		fmt.Printf("%-20s  %-10s %s\n", item.Name,
			item.Status.NodeInfo.Architecture,
			item.Status.NodeInfo.OSImage)
	}
}
