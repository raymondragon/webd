package main
import (
    "flag"
    "log"
    "net/http"
    "os"
    "strings"
    "github.com/google/uuid"
)
var (
    auth = flag.String("a", "/auth", "")
    bind = flag.String("b", ":8080", "")
)
func main() {
    flag.Parse()
    http.HandleFunc(*auth, func(w http.ResponseWriter, r *http.Request) {
        ip := strings.Split(r.RemoteAddr, ":")[0]
        w.Write([]byte("ID: " + uuid.New().String() + "\n" + "IP: " + ip + "\n"))
        file, _ := os.OpenFile("ip.list", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        defer file.Close()
        file.WriteString(ip + "\n")
    })
    log.Fatal(http.ListenAndServe(*bind, nil))
}