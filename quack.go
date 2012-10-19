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

    finished := make(chan int)
    concurrency := 0

    for _,service := range os.Args[1:] {
        parts := strings.SplitN(service, ":", 2)

        host := parts[0]
        port, err := strconv.Atoi(parts[1])
        if err != nil {
            fmt.Fprintf(os.Stderr, "Ignoring invalid service: %s\n", service)
            continue
        }

        concurrency++
        go isUp(host, port, finished)
    }

    for i :=  0; i < concurrency; _, i = <- finished, i+1 { }

    os.Exit(0)
}

func isUp(host string, port int, done chan int) {
    service := fmt.Sprintf("%s:%d", host, port)

    conn, err := net.DialTimeout("tcp", service, TIMEOUT)
    if err != nil {
        goto ERROR
    }

    _, err = conn.Write([]byte("0x00"))
    if err != nil {
        goto ERROR
    }

    hostIsUp(service)
    goto FINAL


    ERROR:
        hostIsDown(service)

    FINAL:
        done <- 1

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
