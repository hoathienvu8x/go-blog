package main

import (
    "fmt"
    "encoding/json"
    "log"
    "os"
    "runtime"
    "io/ioutil"
    "os/exec"
    "strconv"
    "strings"
    "os/signal"
    "syscall"
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

var PIDFile = "/tmp/goblog.pid"

func savePID(pid int) {
    file, err := os.Create(PIDFile)
    if err != nil {
        log.Printf("Unable to create pid file : %v\n", err)
        os.Exit(1)
    }
    defer file.Close()
    _, err = file.WriteString(strconv.Itoa(pid))
    if err != nil {
        log.Printf("Unable to create pid file : %v\n", err)
        os.Exit(1)
    }
    file.Sync()
}

func init() {
    log.SetFlags(log.Lshortfile)
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
    if len(os.Args) != 2 {
        fmt.Printf("Usage : %s [start|stop] \n ", os.Args[0])
        os.Exit(0)
    }
    if strings.ToLower(os.Args[1]) == "main" {
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
        go func() {
            signalType := <- ch
            signal.Stop(ch)
            fmt.Println("Exit command received. Exiting...")
            fmt.Println("Received signal type : ", signalType)
            os.Remove(PIDFile)
            os.Exit(0)
        }()
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
    if strings.ToLower(os.Args[1]) == "start" {
        if _, err := os.Stat(PIDFile); err == nil {
            fmt.Println("Already running or pid file exist.")
            os.Exit(1)
        }
        cmd := exec.Command(os.Args[0],"main")
        cmd.Start()
        fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
        savePID(cmd.Process.Pid)
        os.Exit(0)
    }
    if strings.ToLower(os.Args[1]) == "stop" {
        if _, err := os.Stat(PIDFile); err == nil {
            data, err := ioutil.ReadFile(PIDFile)
            if err != nil {
                fmt.Println("Not running")
                os.Exit(1)
            }
            ProcessID, err := strconv.Atoi(string(data))
            if err != nil {
                fmt.Println("Unable to read and parse process id found in ", PIDFile)
                os.Exit(1)
            }
            process, err := os.FindProcess(ProcessID)
            if err != nil {
                fmt.Printf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
                os.Exit(1)
            }

            os.Remove(PIDFile)
            fmt.Printf("Kill pricess ID [%v] now.\n", ProcessID)
            
            err = process.Kill()
            if err != nil {
                fmt.Printf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
                os.Exit(1)
            } else {
                fmt.Printf("Kill process ID [%v]\n", ProcessID)
                os.Exit(0)
            }
        } else {
            fmt.Println("Not running.")
            os.Exit(1)
        }
    } else {
        fmt.Printf("Unknown command : %v\n", os.Args[1])
        os.Exit(1)
    }
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
