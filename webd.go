package main

import (
    "flag"
    "log"
    "net"
    "net/http"

    "golang/x/webdav"
    "github.com/raymondragon/golib"
)

var rawURL = flag.String("url", "", "http(s)://host:port/path#dir")

func main() {
    flag.Parse()
    if *rawURL == "" {
        flag.Usage()
        log.Fatalf("[ERRO] %v", "Invalid Flag")
    }
    parsedURL, err := golib.URLParse(*rawURL)
    if err != nil {
        log.Printf("[WARN] %v", err)
    }
    http.Handle(parsedURL.Path, &webdav.Handler{
        FileSystem: webdav.Dir(parsedURL.Fragment),
        LockSystem: webdav.NewMemLS(),
    })
    log.Printf("[INFO] %v", *rawURL)
    if err := http.ListenAndServe(net.JoinHostPort(parsedURL.Hostname, parsedURL.Port), nil); err != nil {
        log.Fatalf("[ERRO] %v", err)
    }
}
