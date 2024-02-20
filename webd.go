package main
import (
    "context"
    "flag"
    "log"
    "net/http"
    "golang.org/x/net/webdav"
)
var (
    bind = flag.String("b", ":80", "bind-to")
    dirt = flag.String("d", ".", "directory")
    pref = flag.String("p", "/web", "prefix")
)
type noRemoval struct {
    webdav.FileSystem
}
func (fs *noRemoval) RemoveAll(ctx context.Context, name string) error {
    return nil
}
func main() {
    flag.Parse()
    log.Printf("[LISTEN] %v%v [DIRECTORY] %v\n", *bind, *dirt, *pref)
    log.Fatal(http.ListenAndServe(*bind, &webdav.Handler{
        FileSystem: &noRemoval{webdav.Dir(*dirt)},
        Prefix:     *pref,
        LockSystem: webdav.NewMemLS(),
    }))
}
