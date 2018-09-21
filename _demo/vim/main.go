package main

import (
    "fmt"

    "github.com/wxnacy/goterminal"
    "github.com/nsf/termbox-go"
    "os"
    "strconv"
    "strings"
)

func LogFile(str ...string) {
    file, _ := os.OpenFile("wsh.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    file.WriteString(strings.Join(str, " ") + "\n")
}

type Mode int

const (
    ModeNormal Mode = iota
    ModeInsert
)

type Vim struct {
    Mode Mode
}

var vim = &Vim{Mode: ModeNormal}

func switchModeNormal(t *terminal.Terminal, e termbox.Event) {
    switch e.Key {
        case termbox.KeyEsc: {
            os.Exit(0)
        }
    }

    if e.Ch > 0 {
        switch e.Ch {
            case 'q': {
                os.Exit(0)
            }
            case 'l': {
                t.MoveCursor(1, 0)
            }
            case 'h': {
                t.MoveCursor(-1, 0)
            }
            case 'j': {
                t.MoveCursor(0, 1)
            }
            case 'i': {
                vim.Mode = ModeInsert
            }
            case 'k': {
                t.MoveCursor(0, -1)
            }
            case 'g': {
                if t.E.PreCh == 'g' {
                    t.MoveCursorToFirstLine()
                }
            }
            case 'G': {
                // t.SetCursor(0, t.Height - 1)
                t.MoveCursorToLastLine()
            }
            case '0': {
                t.SetCursor(0, t.CursorY)
            }
            case '$': {
                t.SetCursor(t.LineWidth() - 1, t.CursorY)
            }
            default: {

            }
        }
    }
}

func InitData() string {

    outs := make([]string, 0)
    for i := 0; i < 20; i++ {
        outs = append(outs, strconv.Itoa(i))
    }
    out := strings.Join(outs, "\n")

    out = `
This is Vim Mode

Use i insert ` + out
    return out
}


func main() {
    out := InitData()

    t, err := terminal.NewFromString(out)
    // t.SetMode(terminal.ModeInsert)
    if err != nil {
        panic(err)
    }
    defer t.Close()
    for {
        t.Rendering()
        e := t.PollEvent()
    switch e.Key {
        case termbox.KeyEsc: {
            os.Exit(0)
        }
    }
        t.ListenKeyBoardLikeVim(e)
        switch vim.Mode {
            case ModeNormal: {
                // switchModeNormal(t, e)
                // t.ListenKeyBoardLikeVim(e)
            }
            case ModeInsert: {
                switch e.Key {
                    case termbox.KeyEsc: {
                        vim.Mode = ModeNormal
                    }
                }
                if e.Ch > 0 {
                    t.SetCellBeforeCursor(
                        e.Ch,
                        termbox.ColorDefault,
                        termbox.ColorDefault,
                    )
                }
            }
        }
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

        fmt.Println("Hello World")

}

