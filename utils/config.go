package utils

import (
  "fmt"
  "os"
  "time"
  "strings"
)

const (
  // DefaultKhaosInterval is the default time in seconds between khaos events
  DefaultKhaosInterval = "30s"
  // DefaultKhaoticEvents default set of khaotic events, this is all of them
  DefaultKhaoticEvents = "kill-pods,drain-node"
  // DefaultKhaosDuration is the default duration of time khaos-monkey will wreak havoc
  DefaultKhaosDuration = "10m"
)

type Config struct {
  Name string
  Namespace string
  KhaosInterval time.Duration
  KhaoticEvents []string
  KhaosDuration time.Duration
}

func NewConfig() (conf *Config, err error) {
  conf = &Config{}

  // get namespace to wreak havoc in
  var namespace string
  if namespace = os.Getenv("KHAOS_NAMESPACE"); namespace == "" {
    // shouldn't error our b.c we are using Downward API, but better safe than sorry
    return nil, fmt.Errorf("No namespace available in environment")
  }

  conf.Namespace = namespace

  // get self name so as to avoid suicide
  var name string
  if name = os.Getenv("KHAOS_MONKEY_NAME"); name == "" {
    // shouldn't error our b.c we are using Downward API, but better safe than sorry
    return nil, fmt.Errorf("No name available in environment")
  }

  conf.Name = name

  // get interval between khaos events, or default
  var intervalStr string
  if intervalStr = os.Getenv("KHAOS_INTERVAL"); intervalStr == "" {
    intervalStr = DefaultKhaosInterval
  }

  conf.KhaosInterval, err = time.ParseDuration(intervalStr)
  if err != nil { return nil, err }

  // get list of acceptable events
  var events string
  if events = os.Getenv("KHAOTIC_EVENTS"); events == "" || events == "all" {
    events = DefaultKhaoticEvents
  }

  conf.KhaoticEvents = strings.Split(events, ",")

  // get duration of khaos
  var durationStr string
  if durationStr = os.Getenv("KHAOS_DURATION"); durationStr == "" {
    durationStr = DefaultKhaosDuration
  }

  conf.KhaosDuration, err = time.ParseDuration(durationStr)
  if err != nil { return nil, err }

  return
}