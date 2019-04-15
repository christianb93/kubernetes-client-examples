package main

import (
	"fmt"
	"time"

	"k8s.io/client-go/util/workqueue"
)

//
// This function will simply fill the queue with five
// words and then exit
func fillQueue(queue workqueue.Interface) {
	time.Sleep(time.Second)
	queue.Add("this")
	queue.Add("is")
	queue.Add("a")
	queue.Add("complete")
	queue.Add("sentence")
	fmt.Println("fillQueue completed, sending shutdown")
	queue.ShutDown()
}

//
// Read from the queue and print results
//
func readFromQueue(queue workqueue.Interface, stop chan int) {
	time.Sleep(3 * time.Second)
	for {
		item, shutdown := queue.Get()
		if shutdown {
			// signal that we are done
			stop <- -1
			return
		}
		fmt.Printf("Got item[shutdown = %t]: %s\n", shutdown, item)
		queue.Done(item)
	}
}

func main() {
	fmt.Println("Starting main")
	// Create a channel
	stop := make(chan int)
	// Create a queue.
	myQueue := workqueue.New()
	// Create our first worker thread.  This goroutine will
	// simply put five words into the queue after one second
	// has passed
	go fillQueue(myQueue)
	// Create second thread that will read from the queue
	go readFromQueue(myQueue, stop)
	fmt.Println("Goroutines started, now waiting for reader to complete")
	<-stop
	fmt.Println("Reader signaled completion, exiting")
}
