package terminal

import (
    "github.com/nsf/termbox-go"
    "github.com/mattn/go-runewidth"
    "strconv"
)

type Pane struct {
    Width, Height int               // 窗口的宽高
    hasCursor bool                  // 是否有光标
	CursorX, CursorY int            // 光标坐标
    PositionX, PositionY int
    xBegin, xEnd int
    PageWidth, PageHeight int       // 内容的宽高
	PageOffsetX, PageOffsetY int    // 内容和窗口的偏移坐标
    E *Event
    Mode Mode
    cells [][]Cell
    viewCells [][]Cell
}

func newPane(w, h int) *Pane {
    p := &Pane{
        Width: w,
        Height: h,
        hasCursor: true,
        E: &Event{},
        cells: make([][]Cell, 0),
        viewCells: make([][]Cell, 0),
        Mode: ModeNormal,
        xBegin: DefaultXBegin,
        xEnd: DefaultXEnd,
    }

    return p
}

func printCells(cells [][]Cell) {
    chs := make([]rune, 0)

    for _, yd := range cells {
        for _, d := range yd {
            chs = append(chs, d.Ch)
        }
    }

    LogFile("cells", string(chs))
}

func (p *Pane) AppendCellFromString(s string) {
    p.AppendCellFromStringWithColor(s, ColorDefault, ColorDefault)
    printCells(p.cells)
    printCells(p.viewCells)
}

func (p *Pane) AppendCellFromStringWithColor(s string, fg, bg termbox.Attribute) {
    cells := stringToCellsWithColor(s, fg, bg)
    for _, d := range cells {
        p.cells = append(p.cells, d)
    }
    p.reset()
}


func (p *Pane) SetMode(m Mode) {
    p.Mode = m
}

func (p *Pane) SetCellBeforeCursor(ch rune, fg, bg termbox.Attribute) {

    cell := Cell{Ch: ch, Fg: fg, Bg: bg}
    runeW := runewidth.RuneWidth(ch)
    line := p.realLine()

    newLine := make([]Cell, 0)

    index := p.CursorX

    for i, d := range line {
        if i == index {
            newLine = append(newLine, cell)
            if cell.Width() > 1 {
                for i := 0; i < cell.Width() - 1; i++ {
                    newLine = append(newLine, newCell(0))
                }
            }
        }
        newLine = append(newLine, d)
    }

    if cellsWidth(line) == index {
        newLine = append(newLine, cell)
        if cell.Width() > 1 {
            for i := 0; i < cell.Width() - 1; i++ {
                newLine = append(newLine, newCell(0))
            }
        }
    }

    p.cells[p.CursorY + p.PageOffsetY] = newLine
    p.CursorX = p.CursorX + runeW

}

func (p *Pane) insert(ch rune) {

    cell := newCell(ch)
    runeW := runewidth.RuneWidth(ch)
    line := p.realLine()
    lineWidth := cellsWidth(line)

    newLine := make([]Cell, 0)

    index := p.CursorX

    if index == p.xBegin {
        newLine = append(newLine, cell)
        newLine = append(newLine, line...)
    }

    // for i, d := range line {
        // if i == index {
            // newLine = append(newLine, cell)
            // if cell.Width() > 1 {
                // for i := 0; i < cell.Width() - 1; i++ {
                    // newLine = append(newLine, newCell(0))
                // }
            // }
        // }
        // newLine = append(newLine, d)
    // }

    if index > p.xBegin && index < lineWidth {
        newLine = append(newLine, line[:index]...)
        newLine = append(newLine, cell)
        if cell.Width() > 1 {
            for i := 0; i < cell.Width() - 1; i++ {
                newLine = append(newLine, newNilCell())
            }
        }
        newLine = append(newLine, line[index:]...)
    }

    if lineWidth == index {
        newLine = append(newLine, line...)
        newLine = append(newLine, cell)
        if cell.Width() > 1 {
            for i := 0; i < cell.Width() - 1; i++ {
                newLine = append(newLine, newCell(0))
            }
        }
    }

    p.cells[p.CursorY + p.PageOffsetY] = newLine
    p.CursorX = p.CursorX + runeW

}

func (p *Pane) delete(length int) {
    if length <= 0 {
        return
    }
    newLine := make([]Cell, 0)
    line := p.realLine()

    targetCell := getPrevCell(line, p.CursorX, length)
    if targetCell.Ch == 0 {
        return
    }

    targetX, _ := targetCell.Position()

    newLine = append(newLine, line[0:targetX]...)
    newLine = append(newLine, line[p.CursorX:]...)

    p.cells[p.CursorY + p.PageOffsetY] = newLine

    p.CursorX = targetX

}

func (p *Pane) realLine() []Cell {
    return p.getRealLine(p.CursorY)
}

func (p *Pane) getRealLine(y int) []Cell {
    if len(p.cells) > 0 && y >= 0{
        return p.cells[y + p.PageOffsetY]
    }
    return make([]Cell, 0)
}

func (p *Pane) Line() []Cell {
    return p.GetLine(p.CursorY)
}

func (p *Pane) GetLine(y int) []Cell {
    if len(p.viewCells) > 0 && y < len(p.viewCells) && y >= 0{
        res := p.viewCells[y]
        return res
    }
    return make([]Cell, 0)
}

func (p *Pane) ResetPageSize() {
    p.PageHeight = len(p.cells)
}

// 重置显示的 cell 集合
func (p *Pane) ResetViewCells() {
    viewCells := make([][]Cell, 0)
    minLine := min(p.PageHeight, p.Height)
    for y := 0; y < minLine; y++ {
        newLine := make([]Cell, 0)
        index := min(y + p.PageOffsetY, p.PageHeight - 1)
        line := p.cells[index]
        for x := p.PageOffsetX; x < len(line); x++ {
            d := line[x]
            newLine = append(newLine, d)
        }
        viewCells = append(viewCells, newLine)
    }
    p.viewCells = viewCells
}

func (p *Pane) ResetCellPosition() {
    for y, line := range p.cells {
        for x, d := range line {
            d.x = x
            d.y = y
            p.cells[y][x] = d
        }
    }
}


func (p *Pane) nextCell() Cell {
    x := p.CursorX
    y := p.CursorY
    return getNextCell(p.getRealLine(y), x, 1)
    // res := Cell{}
    // Loop:
    // for i := 1; i < cellsWidth(p.Line()); i++ {
        // x += i
        // c := p.viewCells[y][x]
        // if c.Ch > 0 {
            // c.setPosition(x, y)
            // res = c
            // break Loop
        // }
    // }
    // return res
}

func (p *Pane) realOffsetX(x int) int {
    // line := p.GetLine(p.CursorY)
    offset := 0
    // cursorX := p.CursorX

    if x > 0 {
        // end := p.CursorX + x
        // if end >= len(line) {
            // end = len(line) - 1
        // }
        // rangeCells := line[p.CursorX:end]
        // for _, d := range rangeCells {
            // offset += d.Width()
        // }
        // for i := 0; i < x; i++ {
            // _ = i
            // c := p.nextCell(cursorX, p.CursorY)
            // cursorX, _ = c.Position()
            // LogFile("next ", string(c.Ch))
            // offset = cursorX
        // }
        return offset

    }

    if x < 0 {
        if p.CursorX == 0 {
            return p.CursorX
        }

        // begin := p.CursorX + x
        // if begin < 0 {
            // begin = 0
        // }

        // rangeCells := line[begin:p.CursorX]
        // for _, d := range rangeCells {
            // offset += d.Width()
        // }
        // return -offset
        return p.CursorX - 1

    }

    return p.CursorX
}

// 启动光标
func (p *Pane) MoveCursor(x, y int) {
    line := p.GetLine(p.CursorY)
    minWidth := min(cellsWidth(line), p.Width)
    minHeight := min(p.PageHeight, p.Height)

    cx := p.CursorX + x
    if cx >= minWidth {
        cx = minWidth - 1
        if p.Mode == ModeInsert {
            cx = minWidth
        }
    }

    cy := p.CursorY + y
    if cy >= minHeight {
        cy = minHeight - 1
        if p.Mode == ModeInsert {
            cy = minHeight
        }
    }

    LogFile(
        "move x", strconv.Itoa(x), "y", strconv.Itoa(y), "xbegin",
        strconv.Itoa(p.xBegin), "cx", strconv.Itoa(cx), "cy",
        strconv.Itoa(cy),
        "cursorx",
        strconv.Itoa(p.CursorX), "len(line)", strconv.Itoa(len(line)),
    )
    if cx >= p.xBegin {
        p.CursorX = cx
    }

    if cy >= 0 {
        p.CursorY = cy
    }

    if cx >= p.xBegin || cy >= 0 {
        LogFile(
            "move after cursorx", strconv.Itoa(p.CursorX), "cursory",
            strconv.Itoa(p.CursorY),
        )
        termbox.SetCursor(p.CursorX, p.CursorY)
        termbox.Flush()
    }
    if p.CursorY + 1 == p.Height || p.CursorY == 0{
        maxOffset := max(p.PageHeight, p.Height) - min(p.PageHeight, p.Height)
        p.PageOffsetY = min(p.PageOffsetY + y, maxOffset)
        if p.PageOffsetY < 0 {
            p.PageOffsetY = 0
        }
    }
}

// 渲染
func (p *Pane) Rendering() {
    termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

    p.reset()

    for y, yd := range p.viewCells {
        x := 0
        for _, d := range yd {
            termbox.SetCell(x, y, d.Ch, d.Fg, d.Bg)
            x += d.Width()
        }
    }

    // 矫正光标，使其最大出现在当前展示行的最后一个字付前
    line := p.Line()
    lineWidth := cellsWidth(line)
    switch p.Mode {
        case ModeInsert: {
            if p.CursorX >= lineWidth {
                p.CursorX = lineWidth
            }
        }
        case ModeNormal: {
            if p.CursorX >= lineWidth {
                minWidth := max(lineWidth, 1)
                p.CursorX = minWidth - 1
            }
        }
    }

    termbox.SetCursor(p.CursorX, p.CursorY)
    termbox.Flush()

}

func (p *Pane) reset() {
    p.ResetPageSize()
    p.ResetCellPosition()
    p.ResetViewCells()
}
