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
    "net/url"
    "time"
    "golang.org/x/net/webdav"
)
var rawURL = flag.String("url", "", "")
type ParsedURL struct {
    Scheme   string
    Hostname string
    Port     string
    Path     string
    Fragment string
}
func main() {
    flag.Parse()
    parsedURL, err := urlParse(*rawURL)
    if err != nil {
        log.Fatalf("[ERRO-0 ] %v", err)
    }
    webd := &webdav.Handler{
        FileSystem: webdav.Dir(parsedURL.Fragment),
        Prefix:     parsedURL.Path,
        LockSystem: webdav.NewMemLS(),
    }
    switch parsedURL.Scheme {
    case "http":
        log.Printf("[INFO-0] %v", *rawURL)
        if err := http.ListenAndServe(parsedURL.Hostname+":"+parsedURL.Port, webd); err != nil {
            log.Fatalf("[ERRO-1] %v", err)
        }
    case "https":
        cert, err := generateCert(parsedURL.Hostname)
        if err != nil {
            log.Fatalf("[ERRO-2] %v", err)
        }
        serv := &http.Server{
            Addr:      parsedURL.Hostname+":"+parsedURL.Port,
            Handler:   webd,
            TLSConfig: &tls.Config{
                Certificates: []tls.Certificate{cert},
            },
        }
        log.Printf("[INFO-1] %v", *rawURL)
        if err := serv.ListenAndServeTLS("", ""); err != nil {
            log.Fatalf("[ERRO-3] %v", err)
        }
    default:
        log.Fatalf("[ERRO-4] %v", "Scheme Not Supported")
    }
}
func urlParse(rawURL string) (ParsedURL, error) {
    u, err := url.Parse(rawURL)
    if err != nil {
        return ParsedURL{}, err
    }
    return ParsedURL{
        Scheme:   u.Scheme,
        Hostname: u.Hostname(),
        Port:     u.Port(),
        Path:     u.Path,
        Fragment: u.Fragment,
    }, nil
}
func generateCert(orgName string) (tls.Certificate, error) {
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
            Organization: []string{orgName},
        },
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(10 * 365 * 24 * time.Hour),
        KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
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