package main
import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "flag"
    "log"
    "math/big"
    "net/http"
    "time"
    "golang.org/x/net/webdav"
)
var (
    add = flag.String("add", ":8443", "Server address")
    dir = flag.String("dir", "./", "Working directory")
    org = flag.String("org", "ORG", "Orgnization name")
    pre = flag.String("pre", "/webd", "Webdav prefix")
    tls = flag.Bool("tls", false, "Enable TLS webdav")
)
func main() {
    flag.Parse()
    webd := &webdav.Handler{
        FileSystem: webdav.Dir(*dir),
        Prefix:     *pre,
        LockSystem: webdav.NewMemLS(),
    }
    if !tls {
        log.Printf("[LISTEN] %v%v [SERVE] %v [TLS] OFF", *add, *pre, *dir)
        if err := http.ListenAndServe(*add, webd); err != nil {
            log.Fatalf("[ERR-00] %v", err)
        }
    }
    cert, err := generateCert(*org)
    if err != nil {
        log.Fatalf("[ERR-01] %v", err)
    }
    serv := &http.Server{
        Addr:      *add,
        Handler:   webd,
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
        },
    }
    log.Printf("[LISTEN] %v%v [SERVE] %v [TLS] ON", *add, *pre, *dir)
    if err := serv.ListenAndServeTLS("", ""); err != nil {
        log.Fatalf("[ERR-02] %v", err)
    }
}
func generateCert(name String) (tls.Certificate, error) {
    priv, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return tls.Certificate{}, err
    }
    serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
    if err != nil {
        return tls.Certificate{}, err
    }
    template := x509.Certificate{
        SerialNumber: serialNumber,
        Subject:      pkix.Name{
                          Organization: []string{name},
                      },
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(10*365*24*time.Hour),
        KeyUsage:     x509.KeyUsageKeyEncipherment|x509.KeyUsageDigitalSignature,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    }
    crtDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
    if err != nil {
        return tls.Certificate{}, err
    }
    crtPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: crtDER})
    keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
    return tls.X509KeyPair(crtPEM, keyPEM)
}