package main
import (
    "flag"
    "log"
    "net/http"
    "golang.org/x/net/webdav"
)
var (
    bind = flag.String("b", ":8080", "")
    cust = flag.String("c", "/webd", "")
    dirt = flag.String("d", ".", "")
)
func main() {
    flag.Parse()
    log.Fatal(http.ListenAndServe(*bind, &webdav.Handler{
        FileSystem: webdav.Dir(*dirt),
        Prefix:     *cust,
        LockSystem: webdav.NewMemLS(),
    }))
}
