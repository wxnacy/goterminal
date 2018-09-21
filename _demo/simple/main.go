package main

import (
    "fmt"
    "github.com/wxnacy/goterminal"
    "github.com/nsf/termbox-go"
)

func onCh(ch rune) {
    fmt.Printf("you enter %s\n", string(ch))
}

func onKey(key terminal.Key) {
    fmt.Printf("you enter rune %d\n", key)
}

func main() {
    t, err := terminal.New()
    if err != nil {
        panic(err)
    }

    t.Run(onCh, onKey)
    fmt.Println("Hello World ")
}
