package khaos

import (
  "fmt"
  "math/rand"
  "time"

  "github.com/30x/khaos-monkey/utils"

  "k8s.io/client-go/1.4/kubernetes"
  "k8s.io/client-go/1.4/pkg/api"
)

const (
  killPodsEventStr = "kill-pods"
  drainNodeEventStr = "drain-node"
  targetDaemonsetEventStr = "target-daemonset"
)

// RunRandomKhaoticEvent randomly picks from the acceptable chaotic events and executes it
func RunRandomKhaoticEvent(clientset *kubernetes.Clientset, khaosConfig *utils.Config) (err error) {
  rand.Seed(time.Now().UTC().UnixNano())
  eventNdx := rand.Intn(len(khaosConfig.KhaoticEvents))

  switch khaosConfig.KhaoticEvents[eventNdx] {
  case killPodsEventStr:
    return KillRandomPod(clientset, khaosConfig)
  case drainNodeEventStr:
    return DrainNode(clientset, khaosConfig)
  case targetDaemonsetEventStr:
    fmt.Println("Targeting a daemonset with a banana")
    break
  }

  return
}

// KillRandomPod kills a random pod in the configured khaos-monkey namespace
func KillRandomPod(clientset *kubernetes.Clientset, khaosConfig *utils.Config) (err error) {
  rand.Seed(time.Now().UTC().UnixNano())

  pods, err := clientset.Core().Pods(khaosConfig.Namespace).List(api.ListOptions{})
  if err != nil { return err }

  numPods := len(pods.Items)
  deletingPod := pods.Items[rand.Intn(numPods)].Name
  tempInt := int64(0)
  if deletingPod != khaosConfig.Name {
    err = clientset.Core().Pods(khaosConfig.Namespace).Delete(deletingPod, &api.DeleteOptions{
      GracePeriodSeconds: &tempInt,
    })

    if err != nil { return err }

    fmt.Printf("Killed Pod: %s\n", deletingPod)
  }

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
  graceTime := int64(0)

  pods, err := podsInter.List(api.ListOptions{})
  fmt.Printf("Draining Node: %s\n", drainingNode.Name)

  for _, pod := range pods.Items {
    if pod.Spec.NodeName == drainingNode.Name && pod.Name != khaosConfig.Name {
      err = podsInter.Delete(pod.Name, &api.DeleteOptions{
        GracePeriodSeconds: &graceTime,
      })

      if err != nil { return err }

      fmt.Printf("Drained Pod: %s\n", pod.Name)
    }
  }

  return
}