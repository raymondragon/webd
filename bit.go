package main
import (
    "bufio"
    "flag"
    "io"
    "log"
    "net"
    "os"
    "strings"
)
var (
    bin = flag.String("b", ":9000", "local IP:port")
    ips = flag.String("i", "IPlist", "path to list")
    tar = flag.String("t", "", "target IP:port")
)
func main() {
    flag.Parse()
    if *tar == "" {
        log.Fatal("Target Service Required")
    }
    listener, err := net.Listen("tcp", *bin)
    if err != nil {
        log.Fatal("Error listening: ", err)
    }
    defer listener.Close()
    for {
        clientConn, err := listener.Accept()
        if err != nil {
            log.Println("Error accepting: ", err)
            continue
        }
        go handleClient(clientConn)
    }
}
func handleClient(clientConn net.Conn) {
    defer clientConn.Close()
    clientIP := clientConn.RemoteAddr().(*net.TCPAddr).IP.String()
    if !isIPInList(clientIP, *ips) {
        log.Printf("Client %s not in IPlist, closing connection.\n", clientIP)
        return
    }
    serverConn, err := net.Dial("tcp", *tar)
    if err != nil {
        log.Println("Error connecting to server: ", err)
        return
    }
    defer serverConn.Close()
    go io.CopyBuffer(serverConn, clientConn, nil)
    io.CopyBuffer(clientConn, serverConn, nil)
}
func isIPInList(ip string, iplist string) bool {
    file, err := os.Open(iplist)
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