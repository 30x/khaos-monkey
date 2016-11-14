package main

import (
	"time"

	"github.com/30x/khaos-monkey/utils"
	"github.com/30x/khaos-monkey/khaos"

	emoji "gopkg.in/kyokomi/emoji.v1"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/rest"
)

func main() {
	khaosConfig, err := utils.NewConfig()
	if err != nil { panic(err.Error()) }

	emoji.Println(":see_no_evil: :hear_no_evil: :speak_no_evil:")

	config, err := rest.InClusterConfig()
	if err != nil { panic(err.Error()) }

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil { panic(err.Error()) }

	// event loop
	for {
		// put the monkey to sleep first
		time.Sleep(khaosConfig.KhaosInterval * time.Second)

		// wake up and wreak havoc
		err = khaos.RunRandomKhaoticEvent(clientset, khaosConfig)
		if err != nil { panic(err.Error()) }
	}

}
