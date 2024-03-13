package main

import (
    "flag"
    "log"

    "github.com/raymondragon/golib"
)

var rawURL = flag.String("url", "", "http(s)://host:port/path#directory")

func main() {
    flag.Parse()
    parsedURL, err := golib.URLParse(*rawURL)
    if err != nil {
        log.Printf("[WARN] %v", err)
    }
    webdavHandler := golib.WebdavHandler(parsedURL.Fragment, parsedURL.Path)
    tlsConfig := golib.TLSConfigInit()
    switch parsedURL.Scheme {
    case "http":
    case "https":
        tlsConfig, err = golib.TLSConfigGeneration(parsedURL.Hostname)
        if err != nil {
            log.Printf("[WARN] %v", err)
        }
    default:
        log.Fatalf("[ERRO] %v", parsedURL.Scheme)
    }
    log.Printf("[INFO] %v", *rawURL)
    if err := golib.ServeHTTP(parsedURL.Hostname, parsedURL.Port, webdavHandler, tlsConfig); err != nil {
        log.Fatalf("[ERRO] %v", err)
    }
}