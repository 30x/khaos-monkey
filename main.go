package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/30x/khaos-monkey/utils"

	emoji "gopkg.in/kyokomi/emoji.v1"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/rest"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	khaosConfig, err := utils.NewConfig()
	if err != nil {
		panic(err.Error())
	}

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
		pods, err := clientset.Core().Pods(khaosConfig.Namespace).List(api.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the %s namespace\n", len(pods.Items), khaosConfig.Namespace)
		numPods := len(pods.Items)
		deletingPod := pods.Items[rand.Intn(numPods)].Name
		tempInt := int64(0)
		if deletingPod != khaosConfig.Name {
			clientset.Core().Pods(khaosConfig.Namespace).Delete(deletingPod, &api.DeleteOptions{
				GracePeriodSeconds: &tempInt,
			})
			fmt.Printf("Killed Pod: %s\n", deletingPod)
		}

		time.Sleep(khaosConfig.KhaosInterval * time.Second)
	}

}
