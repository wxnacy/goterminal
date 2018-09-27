package terminal

import (
    "github.com/nsf/termbox-go"
    "strings"
    "github.com/mattn/go-runewidth"
    "strconv"
)

type Cell struct {
    x, y int
    Ch rune
    Fg termbox.Attribute    // 文字颜色
    Bg termbox.Attribute    // 背景颜色
}


func newCell(ch rune) Cell {
    return Cell{
        Ch: ch,
        Fg: termbox.ColorDefault,
        Bg: termbox.ColorDefault,
    }
}

func newNilCell() Cell {
    return Cell{}
}

func (c *Cell) Width() int {
    return runewidth.RuneWidth(c.Ch)
}

func (c *Cell) setPosition(x, y int) {
    c.x = x
    c.y = y
}

func (c *Cell) Position() (int, int) {
    return c.x, c.y
}

func (c *Cell) IsEmpty() bool {
    if c.Ch > 0 {
        return false
    }
    return true
}

func getNextCell(cells []Cell, x, n int) Cell {
    count := 0
    var res Cell
    Loop:
    for i := 1; i < len(cells) - x; i++ {
        c := cells[x + i]
        if c.Ch > 0 {
            res = c
            count++
        }
        if count == n {
            break Loop
        }
    }
    return res
}

func getPrevCell(cells []Cell, x, n int) Cell {
    count := 0
    var res Cell
    if x <= 0 {
        return res
    }
    Loop:
    for i := 1; i <= x; i++ {
        c := cells[x - i]
        if c.Ch > 0 {
            res = c
            count++
            LogFile(strconv.Itoa(count))
        }
        if count == n {
            break Loop
        }
    }
    return res
}


func cellsWidth(cells []Cell) int {
    // length := 0
    // for _, d := range cells {
        // length += d.Width()
    // }
    // return length
    return len(cells)
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
    for _, d := range lines {
        ycells := stringToLineWithColor(d, fg, bg)
        cells = append(cells, ycells)
    }
    return cells
}

func stringToLineWithColor(s string, fg, bg termbox.Attribute) []Cell {
    cells := make([]Cell, 0)
    chs := []rune(s)
    for _, c := range chs {
        cell := Cell{
            Ch: c,
            Fg: fg,
            Bg: bg,
        }
        cells = append(cells, cell)
        for i := 1; i < cell.Width(); i++ {
            cells = append(cells, Cell{})
        }
    }
    return cells
}
