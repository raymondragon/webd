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
    webd := &webdav.Handler{
        FileSystem: webdav.Dir(*dirt),
        Prefix:     *path,
        LockSystem: webdav.NewMemLS(),
    }
    cert, err := generateCert()
    if err != nil {
        log.Fatalf("[ERR-0] %v", err)
    }
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
    }
    serv := &http.Server{
        Addr:      *addr,
        Handler:   webd,
        TLSConfig: tlsConfig,
    }
    log.Printf("[LISTEN] %v%v [SERVE] %v", *addr, *path, *dirt)
    if err := serv.ListenAndServeTLS("", ""); err != nil {
        log.Fatalf("[ERR-1] %v", err)
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
        Subject: pkix.Name{Organization: []string{"webd"},},
        NotBefore: time.Now(),
        NotAfter: time.Now().Add(10*365*24*time.Hour),
        KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
    }
    crtDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
    if err != nil {
        return tls.Certificate{}, err
    }
    crtPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
    keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
    return tls.X509KeyPair(crtPEM, keyPEM)
}