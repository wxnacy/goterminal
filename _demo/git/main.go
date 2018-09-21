package main

import (
    "fmt"
    "github.com/wxnacy/goterminal"
    // "github.com/nsf/termbox-go"
    "os/exec"
)



func main() {
    t, err := terminal.New()
    if err != nil {
        panic(err)
    }

    gst := exec.Command("/bin/bash", "-c", "git status -s")
    bytes, err := gst.Output()
    if err != nil {
        panic(err)
    }

    t.AppendCellFromString(string(bytes))

    for {
        t.Rendering()
        e := t.PollEvent()
        t.ListenKeyBorad(e)
    }

    fmt.Println("Hello World ")
}
