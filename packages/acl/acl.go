package acl

import (
    "net/http"
    "goblog/packages/session"
)

func DisallowAuth(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sess := session.Instance(r)
        if sess.Values["id"] != nil {
            http.Redirect(w, r, "/", http.StatusFound)
            return
        }
        h.ServeHTTP(w, r)
    })
}

func DisallowAnon(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sess := session.Instance(r)
        if sess.Values["id"] == nil {
            http.Redirect(w, r, "/", http.StatusFound)
            return
        }
        h.ServeHTTP(w, r)
    })
}
