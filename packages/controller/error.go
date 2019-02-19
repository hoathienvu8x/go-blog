package controller

import (
    "fmt"
    "net/http"
)

func Error404(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "404 Not Found")
}

func Error500(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprint(w, "500 Internal Server Error")
}

func InvalidToken(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type","text/html")
    w.WriteHeader(http.StatusForbidden)
    fmt.Fprint(w, `You token <strong>expired</strong>, click <a href="javascript:void(0);" onclick="location.replace(document.referrer)">here</a> to try again.`)
}
