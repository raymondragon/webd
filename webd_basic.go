package main

import (
    "flag"
    "log"

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
    webdavHandler := golib.WebdavHandler(parsedURL.Fragment, parsedURL.Path)
    log.Printf("[INFO] %v", *rawURL)
    if err := golib.ServeHTTP(parsedURL.Hostname, parsedURL.Port, webdavHandler); err != nil {
        log.Fatalf("[ERRO] %v", err)
    }
}
