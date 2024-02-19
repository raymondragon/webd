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
func main() {
    flag.Parse()
    type CustomFileSystem struct {
        webdav.FileSystem
    }
    func (fs *CustomFileSystem) Delete(name string) error {
        return webdav.ErrForbidden
    }
    log.Fatal(http.ListenAndServe(*bind, &webdav.Handler{
        FileSystem: &CustomFileSystem{webdav.Dir(*dirt)},
        Prefix:     *pref,
        LockSystem: webdav.NewMemLS(),
    }))
}
