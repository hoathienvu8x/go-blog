package controller

import (
    "net/http"
    "goblog/packages/session"
    "goblog/packages/view"
)


func HomePage(w http.ResponseWriter, r *http.Request) {
    session := session.Instance(r)
    if session.Values["id"] != nil {
        v := view.New(r)
        v.Name = "auth"
        v.Vars["first_name"] = session.Values["first_name"]
        v.Render(w)
    } else {
        // https://github.com/josephspurrier/gowebapp/tree/master/vendor/app/controller
    }
}
