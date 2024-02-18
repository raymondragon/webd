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
    webd := &webdav.Handler{
        FileSystem: webdav.Dir(*dirt),
        Prefix:     *cust,
        LockSystem: webdav.NewMemLS(),
    }
    log.Fatal(http.ListenAndServe(*bind, webd))
}