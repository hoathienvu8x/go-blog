package config

import (
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
)

type Parser interface {
    ParseJSON([]byte) error
}

func Load(configFile string, p Parser) {
    var err     error
    var absPath string
    var input   = io.ReadCloser(os.Stdin)
    if absPath, err = filepath.Abs(configFile); err != nil {
        log.Fatalln(err)
    }
    if input, err = os.Open(absPath); err != nil {
        log.Fatalln(err)
    }
    jsonBytes, err := ioutil.ReadAll(input)
    input.Close()
    if err != nil {
        log.Fatalln(err)
    }
    if err := p.ParseJSON(jsonBytes); err != nil {
        log.Fatalln("Could not parse %q: %v", configFile, err)
    }
}
