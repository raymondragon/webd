package main

import (
    "flag"
    "log"
    "net/url"

    "github.com/raymondragon/golib"
)

var rawURL = flag.String("url", "", "http://name:port/path#dir")

func main() {
    flag.Parse()
    if *rawURL == "" {
        flag.Usage()
        log.Fatalf("[ERRO] %v", "Invalid Flag")
    }
    parsedURL, err := url.Parse(*rawURL)
    if err != nil {
        log.Printf("[WARN] %v", err)
    }
    webdavHandler := golib.WebdavHandler(parsedURL.Fragment, parsedURL.Path)
    log.Printf("[INFO] %v", *rawURL)
    if err := golib.ServeHTTP(parsedURL.Host, webdavHandler); err != nil {
        log.Fatalf("[ERRO] %v", err)
    }
}