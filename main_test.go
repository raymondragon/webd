package main
import (
    "flag"
    "log"
    "net"
    "net/http"
    "os"
)
var (
    auth = flag.String("a", "/auth", "")
    bind = flag.String("b", ":8080", "")
)
func main() {
    flag.Parse()
    http.HandleFunc(*auth, func(w http.ResponseWriter, r *http.Request) {
        ip, _, _ := net.SplitHostPort(r.RemoteAddr)
        w.Write([]byte(ip+"\n"))
        file, _ := os.OpenFile("IPlist", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        defer file.Close()
        file.WriteString(ip+"\n")
    })
    log.Fatal(http.ListenAndServe(*bind, nil))
}
