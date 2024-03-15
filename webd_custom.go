package main

import (
    "flag"
    "log"
    "os"

    "github.com/raymondragon/golib"
)

var rawURL = flag.String("url", "", "http://host:port/path#dir")

func main() {
    flag.Parse()
    if *rawURL == nil {
        flag.Usage()
        os.Exit(1)
    }
    parsedURL, err := golib.URLParse(*rawURL)
    if err != nil {
        log.Printf("[WARN] %v", err)
    }
    webdavHandler := golib.WebdavHandler(parsedURL.Fragment, parsedURL.Path)
    switch parsedURL.Scheme {
    case "http":
        log.Printf("[INFO] %v", *rawURL)
        if err := golib.ServeHTTP(parsedURL.Hostname, parsedURL.Port, webdavHandler); err != nil {
            log.Fatalf("[ERRO] %v", err)
        }
    default:
        log.Fatalf("[ERRO] %v", parsedURL.Scheme)
    }
}