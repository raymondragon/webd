package main

import (
    "flag"
    "log"
    "net/http"

    "golang.org/x/net/webdav"
    "github.com/raymondragon/golib"
)

var rawURL = flag.String("url", "", "http://host:port/path#dir")

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
    if err := golib.ServeHTTP(parsedURL.Hostname, parsedURL.Port, nil); err != nil {
        log.Fatalf("[ERRO] %v", err)
    }
}