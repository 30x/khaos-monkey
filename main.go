package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	emoji "gopkg.in/kyokomi/emoji.v1"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/rest"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	namespace := os.Getenv("MY_POD_NAMESPACE")
	emoji.Println(":see_no_evil: :hear_no_evil: :speak_no_evil:")

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.Core().Pods(namespace).List(api.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the %s namespace\n", len(pods.Items), namespace)
		numPods := len(pods.Items)
		deletingPod := pods.Items[rand.Intn(numPods)].Name
		tempInt := int64(0)
		clientset.Core().Pods(namespace).Delete(deletingPod, &api.DeleteOptions{
			GracePeriodSeconds: &tempInt,
		})
		fmt.Printf("Killed Pod: %s\n", deletingPod)
		time.Sleep(10 * time.Second)
	}

}
