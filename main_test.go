package main
import (
    "flag"
    "log"
    "net"
    "net/http"
    "os"
    "sync"
)
var (
    auth = flag.String("a", "/auth", "")
    bind = flag.String("b", ":8080", "")
    mute = sync.Mutex{}
)
func main() {
    flag.Parse()
    file, err := os.OpenFile("IPlist", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("[ERR-0] ", err)
    }
    defer file.Close()
    http.HandleFunc(*auth, func(w http.ResponseWriter, r *http.Request) {
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            log.Println("[ERR-1] ", err)
        }
        w.Write([]byte(ip+"\n"))
        mute.Lock()
        defer mute.Unlock()
        if _, err := file.WriteString(ip + "\n"); err != nil {
            log.Println("[ERR-2] ", err)
        }
    })
    log.Fatal(http.ListenAndServe(*bind, nil))
}
