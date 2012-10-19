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

    if len(os.Args) == 1 {
        fmt.Fprintf(os.Stderr, "Usage: %s host:port [host:port ...]", os.Args[0])
        os.Exit(1)
    }

    for _,service := range os.Args[1:] {
        parts := strings.SplitN(service, ":", 2)

        host := parts[0]
        port, err := strconv.Atoi(parts[1])
        checkError(err)

        isUp(host, port)
    }

    os.Exit(0)
}

func isUp(host string, port int) {
    service := fmt.Sprintf("%s:%d", host, port)

    conn, err := net.DialTimeout("tcp", service, TIMEOUT)
    if err != nil {
        hostIsDown(service)
        return
    }

    _, err = conn.Write([]byte("0x00"))
    if err != nil {
        hostIsDown(service)
        return
    }

    hostIsUp(service)
}

func hostIsUp(service string) {
    fmt.Fprintf(os.Stderr, "%s is UP\n", service)
}

func hostIsDown(service string) {
    fmt.Fprintf(os.Stderr, "%s is DOWN\n", service)
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}
