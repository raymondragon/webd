package main

import (
    "flag"
    "log"

    "github.com/raymondragon/golib"
)

var rawURL = flag.String("url", "", "http(s)://host:port/path#dir")

func main() {
    flag.Parse()
    if *rawURL == nil {
        flag.Usage()
        log.Fatalf("[ERRO] %v", "Invalid Flag")
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
    case "https":
        tlsConfig, err := golib.TLSConfigApplication(parsedURL.Hostname)
        if err != nil {
            log.Printf("[WARN] %v", err)
            tlsConfig, err = golib.TLSConfigGeneration(parsedURL.Hostname)
            if err != nil {
                log.Printf("[WARN] %v", err)
            }
        }
        log.Printf("[INFO] %v", *rawURL)
        if err := golib.ServeHTTPS(parsedURL.Hostname, parsedURL.Port, webdavHandler, tlsConfig); err != nil {
            log.Fatalf("[ERRO] %v", err)
        }
    default:
        log.Fatalf("[ERRO] Invalid Scheme: %v", parsedURL.Scheme)
    }
}
