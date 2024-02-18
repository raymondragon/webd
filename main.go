package main
import (
    "flag"
    "log"
    "net/http"
    "os"
    "golang.org/x/net/webdav"
)
var (
    auth = flag.String("a", "/auth", "")
    bind = flag.String("b", ":8080", "")
    cust = flag.String("c", "/webd", "")
    dirt = flag.String("d", ".", "")
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
    http.Handle(*cust+"/", &webdav.Handler{
        FileSystem: webdav.Dir(*dirt),
        Prefix:     *cust,
        LockSystem: webdav.NewMemLS(),
    })
    log.Fatal(http.ListenAndServe(*bind, nil))
}