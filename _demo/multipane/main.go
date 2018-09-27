package main

import (
    "github.com/wxnacy/goterminal"
    "os"
    "strconv"
    "strings"
)

func LogFile(str ...string) {
    file, _ := os.OpenFile("wsh.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    file.WriteString(strings.Join(str, " ") + "\n")
}

func InitData() string {

    outs := make([]string, 0)
    for i := 0; i < 30; i++ {
        outs = append(outs, strconv.Itoa(i))
    }
    out := strings.Join(outs, "\n")

    out = strings.Repeat("0123456789", 20) + `
This is Vim Mode

测试汉字abc

Use i insert ` + out
    return out
}


func main() {
    out := InitData()

    t, err := terminal.NewFromString(out)
    if err != nil {
        panic(err)
    }
    defer t.Close()
    for {
        t.Rendering()
        e := t.PollEvent()
        t.ListenKeyBoardLikeVim(e)
        LogFile(strconv.Itoa(t.Width) , strconv.Itoa(t.Height))
        LogFile(
            string(e.Ch),
            strconv.Itoa(t.CursorX),
            strconv.Itoa(t.CursorY),
            strconv.Itoa(t.PageOffsetX),
            strconv.Itoa(t.PageOffsetY),
        )
        LogFile(string(e.Key))
    }

}

