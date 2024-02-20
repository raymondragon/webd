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
    auth = flag.String("a", "/auth", "auth")
    bind = flag.String("b", ":1000", "bind")
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
            http.Error(w, "[ERR-1]", http.StatusInternalServerError)
            return
        }
        if _, err := w.Write([]byte(ip+"\n")); err != nil {
            log.Println("[ERR-2] ", err)
            http.Error(w, "[ERR-2]", http.StatusInternalServerError)
            return
        }
        mute.Lock()
        defer mute.Unlock()
        if _, err := file.WriteString(ip + "\n"); err != nil {
            log.Println("[ERR-3] ", err)
            http.Error(w, "[ERR-3]", http.StatusInternalServerError)
            return
        }
    })
    log.Fatal(http.ListenAndServe(*bind, nil))
}
