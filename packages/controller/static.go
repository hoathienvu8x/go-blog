package controller

import (
    "net/http"
    "strings"
)

func Static(w http.ResponseWriter, r *http.Request) {
    if strings.HasSuffix(r.URL.Path, "/") {
        Error404(w, r)
        return
    }
    http.ServeFile(w, r, r.URL.Path[1:])
}
