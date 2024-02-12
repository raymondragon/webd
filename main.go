package main

import (
        "flag"
        "log"
        "net/http"
        "os"
        "strings"
        "sync"
)

var (
        auth = flag.String("a", "/auth", "")
        bind = flag.String("b", ":8080", "")
        seenIPs    = make(map[string]bool)
        seenIPsMux sync.Mutex
)

func main() {
        flag.Parse()
        http.HandleFunc(*auth, func(w http.ResponseWriter, r *http.Request) {
                ip := strings.Split(r.RemoteAddr, ":")[0]
                w.Write([]byte(ip + "\n"))
                seenIPsMux.Lock()
                defer seenIPsMux.Unlock()
                if !seenIPs[ip] {
                        seenIPs[ip] = true
                        ip2file("ip.list", ip)
                }
        })
        log.Fatal(http.ListenAndServe(*bind, nil))
}

func ip2file(filename string, ip string) {
        file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        defer file.Close()
        file.WriteString(ip + "\n")
}
