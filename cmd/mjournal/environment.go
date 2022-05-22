package main

import (
    "fmt"
    "path/filepath"
    "runtime"

    "github.com/gobuffalo/envy"
)

func GetEnvironmentPath() (string, error) {
    var appPath string
    var err error

    switch runtime.GOOS {
    case "windows":
        appPath = filepath.Join(envy.Get("APPDATA", "C:\\\\"), "mjournal")
    case "darwin":
        appPath = filepath.Join(envy.Get("HOME", "~"), ".config", "mjournal")
    case "linux":
        appPath = filepath.Join(envy.Get("HOME", "~"), ".config", "mjournal")
    default:
        err = fmt.Errorf("Unasupported Operating System - %s", runtime.GOOS)
    }
    return  appPath, err
}
