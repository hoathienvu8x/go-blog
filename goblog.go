package main

import (
    "encoding/json"
    "log"
    "os"
    "runtime"
    "goblog/packages/route"
    "goblog/packages/session"
    "goblog/packages/database"
    "goblog/packages/email"
    "goblog/packages/recaptcha"
    "goblog/packages/config"
    "goblog/packages/server"
    "goblog/packages/view"
    "goblog/packages/plugins"
)

func init() {
    log.SetFlags(log.Lshortfile)
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
    rootDir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    config.Load(rootDir + string(os.PathSeparator) + "conf" + string(os.PathSeparator) + "config.json", conf)
    session.Configure(conf.Session)
    recaptcha.Configure(conf.Recaptcha)
    database.Connect(conf.Database)
    view.Configure(conf.View)
    view.LoadTemplates(conf.Template.Root, conf.Template.Children)
    view.LoadPlugins(plugins.TemplateFuncMap(conf.View))
    server.Run(route.LoadHTTP(), route.LoadHTTPS(), conf.Server)
}

var conf = &configuration{}

type configuration struct {
    Database    database.Info   `json:"Database"`
    Email       email.SMTPInfo  `json:"Email"`
    Recaptcha   recaptcha.Info  `json:"Recaptcha"`
    Server      server.Server   `json:"Server"`
    Session     session.Session `json:"Session"`
    Template    view.Template   `json:"Template"`
    View        view.View       `json:"View"`
}

func (c *configuration) ParseJSON(b []byte) error {
    return json.Unmarshal(b, c)
}
