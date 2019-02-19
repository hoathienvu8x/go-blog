package route

import (
    "net/http"
    "goblog/packages/controller"
    "goblog/packages/session"
    "goblog/packages/acl"
    hr "goblog/packages/httprouterwrapper"
    "goblog/packages/logrequest"
    "goblog/packages/pprofhandler"
    "goblog/packages/csrfbanana"
    "goblog/packages/alice"
    "goblog/packages/httprouter"

    "github.com/gorilla/context"
)

func Load() http.Handler {
    return middleware(routers())
}

func LoadHTTPS() http.Handler {
    return middleware(routers())
}

func LoadHTTP() http.Handler {
    return middleware(routers())
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
    http.Redirect(w, req, "https://" + req.Host, http.StatusMovedPermanently)
}

func routers() *httprouter.Router {
    r := httprouter.New()
    r.NotFound = alice.New().ThenFunc(controller.Error404)
    r.GET("/static/*filepath", hr.Handler(alice.New().ThenFunc(controller.Static)))
    // r.GET("/", hr.Handler(alice.New().ThenFunc(controller.IndexGET))
    r.GET("/debug/pprof/*pprof", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(pprofhandler.Handler)))
    return r
}

func middleware(h http.Handler) http.Handler {
    cs := csrfbanana.New(h, session.Store, session.Name)
    cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
    cs.ClearAfterUsage(true)
    cs.ExcludeRegexPaths([]string{"/static(.*)"})
    csrfbanana.TokenLength = 32
    csrfbanana.TokenName = "token"
    csrfbanana.SingleToken = false
    h = cs
    h = logrequest.Handler(h)
    h = context.ClearHandler(h)
    return h
}
