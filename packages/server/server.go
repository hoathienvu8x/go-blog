package server

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

type Server struct {
    HostName    string  `json:"HostName"`
    UseHTTP     bool    `json:"UseHTTP"`
    UseHTTPS    bool    `json:"UseHTTPS"`
    HTTPPort    int     `json:"HTTPPort"`
    HTTPSPort   int     `json:"HTTPSPort"`
    CertFile    string  `json:"CertFile"`
    KeyFile     string  `json:"KeyFile"`
}

func Run(httpHandlers http.Handler, httpsHandlers http.Handler, s Server) {
    if s.UseHTTP && s.UseHTTPS {
        go func() {
            startHTTPS(httpsHandlers, s)
        }()
        startHTTP(httpHandlers, s)
    } else if s.UseHTTP {
        startHTTP(httpHandlers, s)
    } else if s.UseHTTPS {
        startHTTPS(httpsHandlers, s)
    } else {
        log.Println("Config file does not specify a listener to start")
    }
}

func startHTTP(handlers http.Handler, s Server) {
    fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "Running HTTP " + httpAddress(s))
    log.Fatal(http.ListenAndServe(httpAddress(s), handlers))
}

func startHTTPS(handlers http.Handler, s Server) {
    fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), "Running HTTPS " + httpsAddress(s))
    log.Fatal(http.ListenAndServeTLS(httpsAddress(s), s.CertFile, s.KeyFile, handlers))
}

func httpAddress(s Server) string {
    return s.HostName + ":" + fmt.Sprintf("%d", s.HTTPPort)
}

func httpsAddress(s Server) string {
    return s.HostName + ":" + fmt.Sprintf("%d", s.HTTPSPort)
}
