package utils

import (
  "fmt"
  "os"
  "time"
)

const (
  // DefaultKhaosInterval is the default time in seconds between khaos events
  DefaultKhaosInterval = 30
)

type Config struct {
  Name string
  Namespace string
  KhaosInterval time.Duration
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
  if intervalStr := os.Getenv("KHAOS_INTERVAL"); intervalStr == "" {
    conf.KhaosInterval = DefaultKhaosInterval
  } else {
    conf.KhaosInterval, err = time.ParseDuration(intervalStr)
    if err != nil { return nil, err }
  }

  return
}