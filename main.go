package main
import (
    "flag"
    "log"
    "net/http"
    "golang.org/x/net/webdav"
)
var (
    addr = flag.String("a", ":1", "addr")
    dirt = flag.String("d", "./", "dirt")
    path = flag.String("p", "/x", "path")
)
func main() {
    flag.Parse()
    webd := &webdav.Handler{
        FileSystem: webdav.Dir(*dirt),
        Prefix:     *path,
        LockSystem: webdav.NewMemLS(),
    }
    log.Printf("[LISTEN] %v%v [SERVE] %v", *addr, *path, *dirt)
    if err := http.ListenAndServe(*addr, webd); err != nil {
        log.Fatalf("[ERR-0] %v", err)
    }
}