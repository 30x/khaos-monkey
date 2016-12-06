package khaos

import (
  "fmt"
  "math/rand"
  "time"

  "github.com/30x/khaos-monkey/utils"
  emoji "gopkg.in/kyokomi/emoji.v1"

  "k8s.io/client-go/1.4/kubernetes"
  "k8s.io/client-go/1.4/pkg/api"
)

const (
  killPodsEventStr = "kill-pods"
  drainNodeEventStr = "drain-node"
)

// RunRandomKhaoticEvent randomly picks from the acceptable chaotic events and executes it
func RunRandomKhaoticEvent(clientset *kubernetes.Clientset, khaosConfig *utils.Config) (err error) {
  rand.Seed(time.Now().UTC().UnixNano())
  eventNdx := rand.Intn(len(khaosConfig.KhaoticEvents))

  emoji.Printf(":fire: %s :fire:\n", khaosConfig.KhaoticEvents[eventNdx])
  switch khaosConfig.KhaoticEvents[eventNdx] {
  case killPodsEventStr:
    return KillRandomPod(clientset, khaosConfig)
  case drainNodeEventStr:
    return DrainNode(clientset, khaosConfig)
  }

  return
}

// KillRandomPod kills a random pod in the configured khaos-monkey namespace
func KillRandomPod(clientset *kubernetes.Clientset, khaosConfig *utils.Config) (err error) {
  rand.Seed(time.Now().UTC().UnixNano())

  pods, err := clientset.Core().Pods(khaosConfig.Namespace).List(api.ListOptions{})
  if err != nil { return err }

  numPods := len(pods.Items)

  if numPods == 1 {
    fmt.Println("this monkey is the only pod! skipping this event...")
    return
  }

  // select pods until it selects one other than itself
  deletingPod := pods.Items[rand.Intn(numPods)].Name
  for deletingPod == khaosConfig.Name {
    deletingPod = pods.Items[rand.Intn(numPods)].Name
  }

  // delete the selected pod
  gracePeriod := int64(0)
  err = clientset.Core().Pods(khaosConfig.Namespace).Delete(deletingPod, &api.DeleteOptions{
    GracePeriodSeconds: &gracePeriod,
  })

  if err != nil { return err }

  fmt.Printf("Killed Pod: %s\n", deletingPod)

  return
}

// DrainNode targets a random node and drains it of all pods
func DrainNode(clientset *kubernetes.Clientset, khaosConfig *utils.Config) (err error) {
  rand.Seed(time.Now().UTC().UnixNano())

  podsInter := clientset.Core().Pods(api.NamespaceAll)
  nodes, err := clientset.Core().Nodes().List(api.ListOptions{})
  if err != nil { return err }

  numNodes := len(nodes.Items)
  drainingNode := nodes.Items[rand.Intn(numNodes)]
  if drainingNode.Spec.Unschedulable { // for now, we just skip if selected node is already cordoned
    fmt.Printf("%s is already cordoned, skipping this event\n", drainingNode.Name)
    return
  }

  graceTime := int64(0)

  pods, err := podsInter.List(api.ListOptions{})
  fmt.Printf("Draining Node: %s\n", drainingNode.Name)

  // make node unschedulable
  drainingNode.Spec.Unschedulable = true;
  _, err = clientset.Core().Nodes().Update(&drainingNode)
  if err != nil { return err }

  // kill pods on target node
  for _, pod := range pods.Items {
    if pod.Spec.NodeName == drainingNode.Name && pod.Name != khaosConfig.Name {
      err = clientset.Core().Pods(pod.ObjectMeta.Namespace).Delete(pod.Name, &api.DeleteOptions{
        GracePeriodSeconds: &graceTime,
      })

      if err != nil { return err }

      fmt.Printf("Drained Pod: %s\n", pod.Name)
    }
  }

  return
}