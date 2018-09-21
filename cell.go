package terminal

import (
    "github.com/nsf/termbox-go"
    "strings"
)

type Cell struct {
    Ch rune
    Fg termbox.Attribute    // 文字颜色
    Bg termbox.Attribute    // 背景颜色
}

func NewCell(ch rune) Cell {
    return Cell{
        Ch: ch,
        Fg: termbox.ColorDefault,
        Bg: termbox.ColorDefault,
    }
}

func stringToCells(s string) [][]Cell {
    return stringToCellsWithColor(s, ColorDefault, ColorDefault)
}
func StringToCells(s string) [][]Cell {
    return stringToCellsWithColor(s, ColorDefault, ColorDefault)
}

func StringToCellsWithColor(s string, fg, bg termbox.Attribute) [][]Cell {
    return stringToCellsWithColor(s, fg, bg)
}

func stringToCellsWithColor(s string, fg, bg termbox.Attribute) [][]Cell {
    cells := make([][]Cell, 0)
    lines := strings.Split(s, "\n")
    // chss := make([]rune, 0)
    for _, d := range lines {
        ycells := make([]Cell, 0)
        chs := []rune(d)
        for _, c := range chs {
            cell := Cell{
                Ch: c,
                Fg: fg,
                Bg: bg,
            }
            // chss = append(chs, c)
            ycells = append(ycells, cell)
        }
        cells = append(cells, ycells)
    }
    // LogFile("stringToCells", s, string(chss))
    return cells
}
