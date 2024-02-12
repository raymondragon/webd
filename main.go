package main
import (
    "flag"
    "log"
    "net/http"
    "os"
    "strings"
    "sync"
    "github.com/google/uuid"
)
var (
    auth = flag.String("a", "/auth", "")
    bind = flag.String("b", ":8080", "")
    ipst = struct {
        sync.Map
        sync.Mutex
    }{}
)
func main() {
    flag.Parse()
    http.HandleFunc(*auth, func(w http.ResponseWriter, r *http.Request) {
        ip := strings.Split(r.RemoteAddr, ":")[0]
        uuid := uuid.New().String()
        w.Write([]byte("ID: " + uuid + "\n" + "IP: " + ip + "\n"))
        ipst.Lock()
        defer ipst.Unlock()
        _, ok := ipst.Load(ip)
        if !ok {
            ipst.Store(ip, true)
            ipfile("ip.list", ip)
        }
    })
    log.Fatal(http.ListenAndServe(*bind, nil))
}
func ipfile(filename string, ip string) {
    file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer file.Close()
    file.WriteString(ip + "\n")
}