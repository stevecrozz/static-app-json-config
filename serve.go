package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "os/signal"
  "strings"
  "sync"
  "syscall"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
  http.NotFound(w, r)
  fmt.Fprintf(w, "Static app server only serves /static-apps/")
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  appname := strings.Replace(path, mountPath, "", 1)
  c := GetConfig()
  appconfig := c[appname]

  if appconfig == nil {
    http.NotFound(w, r)
    fmt.Fprintf(w, "No app named '%s'", appname)
  } else {
    // merge in stuff from headers
    json, _ := json.Marshal(appconfig)

    fmt.Fprint(w, string(json))
  }
}

func main() {
  loadConfig(true)
  s := make(chan os.Signal, 1)
  signal.Notify(s, syscall.SIGHUP)

  go func() {
    for {
      <-s
      loadConfig(false)
      log.Println("Reloaded")
    }
  }()

  http.HandleFunc(mountPath, staticHandler)
  http.HandleFunc("/", defaultHandler)
  http.ListenAndServe(":8080", nil)
}

type Config map[string]interface{}

var (
  config Config
  configLock = new(sync.RWMutex)
  mountPath string = "/static-apps/"
)

func loadConfig(fail bool){
  file, err := ioutil.ReadFile("config.json")
  if err != nil {
    log.Println("open config: ", err)
    if fail { os.Exit(1) }
  }

  temp := Config{}
  if err = json.Unmarshal(file, &temp); err != nil {
    log.Println("parse config: ", err)
    if fail { os.Exit(1) }
  }
  configLock.Lock()
  config = temp
  configLock.Unlock()
}

func GetConfig() Config {
  configLock.RLock()
  defer configLock.RUnlock()
  return config
}
