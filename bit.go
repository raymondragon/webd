package main
import (
    "bufio"
    "flag"
    "io"
    "log"
    "net"
    "os"
    "strings"
    "time"
)
var (
    bin = flag.String("b", ":9000", "local IP:port")
    ips = flag.String("i", "IPlist", "path to list")
    tar = flag.String("t", "", "target IP:port")
)
func main() {
    flag.Parse()
    if *tar == "" {
        log.Fatal("Target service required.")
    }
    listenAndServe()
}
func listenAndServe() {
    listener, err := net.Listen("tcp", *bin)
    if err != nil {
        log.Fatal("Error listening: ", err)
    }
    defer listener.Close()
    go watchIPList()
    for {
        clientConn, err := listener.Accept()
        if err != nil {
            log.Fatal("Error accepting: ", err)
        }
        go handleClient(clientConn)
    }
}
func handleClient(clientConn net.Conn) {
    defer clientConn.Close()
    clientIP := clientConn.RemoteAddr().(*net.TCPAddr).IP.String()
    if !isIPInList(clientIP) {
        log.Printf("Client %s not in IPlist, closing connection.\n", clientIP)
        return
    }
    serverConn, err := net.Dial("tcp", *tar)
    if err != nil {
        log.Fatal("Error connecting to server: ", err)
    }
    defer serverConn.Close()
    go io.Copy(serverConn, clientConn)
    go io.Copy(clientConn, serverConn)
}
func watchIPList() {
    for range time.Tick(time.Second) {
        if err := loadIPList(); err != nil {
            log.Println("Error loading IPlist: ", err)
        }
    }
}
func loadIPList() error {
    file, err := os.Open(*ips)
    if err != nil {
        return err
    }
    defer file.Close()
    var newIPList []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        newIPList = append(newIPList, strings.TrimSpace(scanner.Text()))
    }
    if err := scanner.Err(); err != nil {
        return err
    }
    return nil
}
func isIPInList(ip string) bool {
    file, err := os.Open(*ips)
    if err != nil {
        log.Println("Error opening IPlist: ", err)
        return false
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if strings.TrimSpace(scanner.Text()) == ip {
            return true
        }
    }
    if err := scanner.Err(); err != nil {
        log.Println("Error scanning IPlist: ", err)
    }
    return false
}