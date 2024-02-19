package main
import (
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
type norm struct {
    webdav.FileSystem
}
func (fs *norm) RemoveAll(ctx context.Context, name string) error {
    return webdav.ErrForbidden
}
func main() {
    flag.Parse()
    log.Fatal(http.ListenAndServe(*bind, &webdav.Handler{
        FileSystem: &norm{webdav.Dir(*dirt)},
        Prefix:     *pref,
        LockSystem: webdav.NewMemLS(),
    }))
}
