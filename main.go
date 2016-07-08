package main

import (
  "fmt"
  "log"
  "strings"
  "time"
  "net/url"
  "net/http"
  "net/http/httputil"
  "github.com/tpbowden/ingress-router/router"
  "github.com/tpbowden/ingress-router/service"
)

var r *router.Router

func handler(w http.ResponseWriter, req *http.Request) {
  log.Printf("Starting request for %s", req.Host)
  dnsName := strings.Split(req.Host, ":")[0]

  s, ok := r.Route(dnsName)

  if ok {
    url, err := url.Parse(s.Url())
    if err != nil {
      fmt.Fprint(w, "Failed to route to service")
      log.Print("Failed to route to service")
    } else {
      log.Printf("Routing to %s", s.Url())
      proxy := httputil.NewSingleHostReverseProxy(url)
      proxy.ServeHTTP(w, req)

    }
  } else {
    fmt.Fprint(w, "Failed to route to service")
    log.Print("Failed to route to service")
  }
}

func updateServices() {
  log.Print("Updating routes")
  services := service.LoadAll()
  r.UpdateTable(services)
}


func main() {
  log.Print("Starting router")
  r = router.NewRouter()
  http.HandleFunc("/", handler)

  go updateServices()

  ticker := time.NewTicker(10 * time.Second)
  quit := make(chan struct{})
  go func() {
    for {
      select {
      case <- ticker.C:
        updateServices()
      case <- quit:
        ticker.Stop()
        return
      }
    }
  }()

  http.ListenAndServe("0.0.0.0:8080", nil)
}
