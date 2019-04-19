//
// Short demonstration for the usage of a shared informer
// to display pod events in a cluster
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

//
// A generic event handler for pod events
//
func handleEvent(pod *v1.Pod, eventType string) {
	fmt.Printf("Recorded event of type %s on pod %s\n", eventType, pod.Name)
}

//
// Our event handler for adding a pod
//
func onAdd(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		fmt.Printf("Error during conversion, this does not seem to be a pod\n")
		return
	}
	handleEvent(pod, "ADD")
}

//
// Our event handler for deletion of a pod
//
func onDelete(obj interface{}) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		fmt.Printf("Error during conversion, this does not seem to be a pod\n")
		return
	}
	handleEvent(pod, "DEL")
}

//
// Our event handler for modification of a pod
//
func onUpdate(old interface{}, new interface{}) {
	pod, ok := old.(*v1.Pod)
	if !ok {
		fmt.Printf("Error during conversion, this does not seem to be a pod\n")
		return
	}
	handleEvent(pod, "MOD")
}

//
// Create a channel that will be closed when a signal is received
//
func createSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		fmt.Printf("Signal handler: received signal %s\n", sig)
		close(stop)
	}()
	return stop
}

func main() {
	//
	// Create a clientset
	//
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		panic(err)
	}
	//
	// Create a Clientset
	//
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	//
	// Implement a signal handler - this will return a channel which
	// will be closed if a signal is received
	//
	stopCh := createSignalHandler()
	//
	// To create an informer for Pods, we will use a factory. This factory
	// expects two argument - a clientset and a resync time (after this time,
	// the cache will be rebuilt from scratch)
	//
	factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
	//
	// We can now ask our factory to create a pod informer for us - this
	// will be a shared informer that is listening to Pods
	//
	podInformer := factory.Core().V1().Pods().Informer()
	fmt.Println("Starting informer")
	//
	// Starting the factory will start all informers created
	// by this factory
	factory.Start(stopCh)
	fmt.Println("Informer running")
	//
	// Wait for informer to sync. We use the helper function
	// WaitForCacheSync which will also take care of signal
	// handling, i.e. it returns when stopCh is closed
	//
	if ok := cache.WaitForCacheSync(stopCh, podInformer.HasSynced); !ok {
		panic("Error while waiting for informer to sync")
	}
	//
	// The informers main loop is now running. We can now add our event handlers
	// and wait for the stopCh to be closed
	//
	fmt.Println("Informer synced, now adding event handlers and waiting for stop channel")
	podInformer.AddEventHandler(
		&cache.ResourceEventHandlerFuncs{
			AddFunc:    onAdd,
			DeleteFunc: onDelete,
			UpdateFunc: onUpdate,
		})
	<-stopCh
}
