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
    addr = flag.String("a", ":1", "addr")
    dirt = flag.String("d", "./", "dirt")
    path = flag.String("p", "/x", "path")
)
func main() {
    flag.Parse()
    cert, err := generateCert()
    if err != nil {
        log.Fatalf("[ERR-00] %v", err)
    }
    serv := &http.Server{
        Addr:      *addr,
        Handler:   &webdav.Handler{
            FileSystem: webdav.Dir(*dirt),
            Prefix:     *path,
            LockSystem: webdav.NewMemLS(),
        },
        TLSConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
        },
    }
    log.Printf("[LISTEN] %v%v [SERVE] %v", *addr, *path, *dirt)
    if err := serv.ListenAndServeTLS("", ""); err != nil {
        log.Fatalf("[ERR-01] %v", err)
    }
}
func generateCert() (tls.Certificate, error) {
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
                          Organization: []string{"webd"},
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