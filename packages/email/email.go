package email

import (
    "encoding/base64"
    "fmt"
    "net/smtp"
)

var (
    e   SMTPInfo
)

type SMTPInfo struct {
    UserName    string
    Password    string
    HostName    string
    Port        int
    From        string
}

func Configure(c SMTPInfo) {
    e = c
}

func LoadConfig() SMTPInfo {
    return e
}

func SendMail(to, subject, body string) error {
    auth := smtp.PlainAuth("", e.UserName, e.Password, e.HostName)
    header := make(map[string]string)
    header["From"] = e.From
    header["To"] = to
    header["Subject"] = subject
    header["MIME-Version"] = "1.0"
    header["Content-Type"] = `text/plain; charset="utf-8"`
    header["Content-Transfer-Encoding"] = "base64"
    message := ""
    for k, v := range header {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
    err := smtp.SendMail(fmt.Sprintf("%s:%d", e.HostName, e.Port), auth, e.From, []string{to}, []byte(message))
    return err
}
