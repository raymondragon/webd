package main

import (
    "log"
    "net/url"
    "os"

    "github.com/raymondragon/golib"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("[ERRO] Usage: http://name:port/path#dir")
    }
    rawURL := os.Args[1]
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        log.Printf("[WARN] %v", err)
    }
    webdavHandler := golib.WebdavHandler(parsedURL.Fragment, parsedURL.Path)
    log.Printf("[INFO] %v", rawURL)
    if err := golib.ServeHTTP(parsedURL.Host, webdavHandler); err != nil {
        log.Fatalf("[ERRO] %v", err)
    }
}