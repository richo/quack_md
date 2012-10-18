package main

import (
    "net"
    "os"
    "strings"
    "fmt"
    "strconv"
    "time"
)

const TIMEOUT = time.Duration(2) * time.Second

func main() {

    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
        os.Exit(1)
    }

    service := os.Args[1]

    parts := strings.SplitN(service, ":", 2)

    host := parts[0]
    port, err := strconv.Atoi(parts[1])
    checkError(err)

    if isUp(host, port) {
        fmt.Fprintf(os.Stderr, "%s is UP", service)
    } else {
        fmt.Fprintf(os.Stderr, "%s is DOWN", service)
        os.Exit(1)
    }
    os.Exit(0)
}

func isUp(host string, port int) bool {
    service := fmt.Sprintf("%s:%d", host, port)

    conn, err := net.DialTimeout("tcp", service, TIMEOUT)
    if err != nil {
        return false
    }

    _, err = conn.Write([]byte("0x00"))
    if err != nil {
        return false
    }

    return true
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}
