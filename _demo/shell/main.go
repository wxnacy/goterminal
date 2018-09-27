package main

import (
    "github.com/wxnacy/goterminal"
    "github.com/nsf/termbox-go"
    "os"
    "strings"
    "strconv"
    "os/exec"
)


type History struct {
    index int
    cmds []string
}

func NewHistory() *History{
    return &History{
        index: -1,
        cmds: make([]string, 0),
    }
}

func (h *History) AddCmd(s string) {
    h.cmds = append(h.cmds, s)
    h.index++
}

func (h *History) GetPrevCmd() string {
    if h.index > 0 {
        c := h.cmds[h.index]
        h.index--
        return c
    }
    return ""
}

func (h *History) GetNextCmd() string {
    if h.index > 0 {
        h.index++
        c := h.cmds[h.index]
        return c
    }
    return ""
}


func main() {
    prefix := "wsh > $ "
    t, err := terminal.New()
    t.AppendCellFromString(prefix)
    t.SetMode(terminal.ModeInsert)
    t.SetLineRange(len(prefix), 0)
    if err != nil {
        panic(err)
    }

    h := NewHistory()

    for {
        t.Rendering()
        maxRemoveWidth := t.LineWidth() - len(prefix)
        LogFile("maxRemoveWidth" + strconv.Itoa(maxRemoveWidth))

        e := t.PollEvent()
        // t.ListenKeyBoardLikeVim(e)

        switch e.Key {
            case termbox.KeyEsc: {
                os.Exit(0)
            }
            case termbox.KeyBackspace: {
                LogFile("KeyBackspace")
                t.RemoveCellBeforeCursor(1)
            }
            case termbox.KeyCtrlW: {
                t.RemoveCellBeforeCursor(4)
            }
            case termbox.KeyArrowUp: {
                hcmd := h.GetPrevCmd()
                if hcmd != "" {
                    // t.AppendCellFromString(hcmd)
                    hchs := []rune(hcmd)
                    for _, d := range hchs {
                        t.SetCellBeforeCursor(
                            d, terminal.ColorDefault, terminal.ColorDefault,
                        )
                    }
                    // t.MoveCursorToLastCell()
                    t.MoveCursorToLineEnd()
                }
            }
            case termbox.KeyEnter: {
                bin := t.SubLineString(t.CursorY, len(prefix), t.LineWidth() - 1)
                LogFile("enter", bin + "--")
                h.AddCmd(bin)
                cmd := exec.Command("/bin/bash", "-c", bin)
                bytes, err := cmd.Output()
                if err != nil {
                    t.AppendCellFromString(err.Error())
                    t.AppendCellFromString(prefix)
                    t.MoveCursorToLastCell()
                } else {
                    t.AppendCellFromString(string(bytes))
                    t.AppendCellFromString(prefix)
                    t.MoveCursorToLastCell()
                }

                _ = bytes
            }
            default: {
                LogFile("default", string(e.Ch))
            }
        }

        if e.Ch > 0 {
            LogFile(string(e.Ch))
            t.SetCellBeforeCursor(e.Ch, termbox.ColorDefault, termbox.ColorDefault)
        }
    }

}

func LogFile(str ...string) {
    file, _ := os.OpenFile("wsh.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    file.WriteString(strings.Join(str, " ") + "\n")
}
func min(a, b int) int{
    if a > b {
        return b
    }
    return a
}
