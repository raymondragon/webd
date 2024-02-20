package main
import (
    "context"
    "flag"
    "log"
    "net/http"
    "golang.org/x/net/webdav"
)
var (
    bind = flag.String("b", ":1", "bind")
    dirt = flag.String("d", "./", "dirt")
    path = flag.String("p", "/1", "path")
)
type noRemoval struct {
    webdav.FileSystem
}
func (fs *noRemoval) RemoveAll(ctx context.Context, name string) error {
    return nil
}
func main() {
    flag.Parse()
    log.Printf("[LISTEN] %v%v [SERVE] %v\n", *bind, *path, *dirt)
    log.Fatal(http.ListenAndServe(*bind, &webdav.Handler{
        FileSystem: &noRemoval{webdav.Dir(*dirt)},
        Prefix:     *path,
        LockSystem: webdav.NewMemLS(),
    }))
}
