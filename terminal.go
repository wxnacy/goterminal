package terminal

import (
    "github.com/nsf/termbox-go"
    "strings"
    "os"
)

type Event struct {
    PreCh rune
    Ch rune
}

type Terminal struct {
    Width, Height    int
	CursorX, CursorY int
    xBegin, xEnd int
    PageWidth, PageHeight int
	PageOffsetX, PageOffsetY int
    E *Event
    Mode Mode
    cells [][]Cell
    viewCells [][]Cell
    panes []*Pane
}

// func (t *Terminal) AddPaneLandscape() *Pane {
    // t.panes = append(t.panes, p)
// }

// func (t *Terminal) AddPaneVertical(p *Pane) {
    // p.hasCursor = true
    // t.panes = append(t.panes, p)
// }


func New() (*Terminal, error){
    err := termbox.Init()
    if err != nil {
        return nil, err
    }

    w, h := termbox.Size()

    p := newPane(w, h)

    return &Terminal{
        Width: w,
        Height: h,
        E: &Event{},
        cells: make([][]Cell, 0),
        viewCells: make([][]Cell, 0),
        Mode: ModeNormal,
        xBegin: DefaultXBegin,
        xEnd: DefaultXEnd,
        panes: []*Pane{p},
    }, nil
}

func NewFromString(s string) (*Terminal, error){
    t, err := New()
    if err != nil {
        return nil, err
    }

    t.AppendCellFromString(s)

    return t, nil
}

func (t *Terminal) HasCursorPane() *Pane{
    for _, d := range t.panes {
        if d.hasCursor {
            return d
        }
    }
    return nil
}

// func (t *Terminal) ResetPageSize() {
    // t.PageHeight = len(t.cells)
// }



func (t *Terminal) SetLineRange(begin, end int) {
    t.xBegin = begin
    t.xEnd = end
    t.CursorX = begin
}

func (t *Terminal) Run(onCh func(ch rune), onKey func(key termbox.Key)) {

    for {
        t.Rendering()
        e := t.PollEvent()

        if e.Key == termbox.KeyEsc {
            os.Exit(0)
        }

        if e.Ch > 0 {
            onCh(e.Ch)
        } else {
            onKey(e.Key)
        }
    }

}

// func (t *Terminal) SetCells(cells [][]Cell) {
    // t.cells = cells
    // t.ResetViewCells()
    // t.ResetPageSize()
// }

// 渲染
func (t *Terminal) Rendering() {
    t.HasCursorPane().Rendering()
}

func (t *Terminal) SetMode(m Mode) {
    t.HasCursorPane().SetMode(m)
}

// 重置显示的 cell 集合
// func (t *Terminal) ResetViewCells() {
    // viewCells := make([][]Cell, 0)
    // minLine := min(t.PageHeight, t.Height)
    // for y := 0; y < minLine; y++ {
        // newLine := make([]Cell, 0)
        // index := min(y + t.PageOffsetY, t.PageHeight - 1)
        // line := t.cells[index]
        // for x := t.PageOffsetX; x < len(line); x++ {
            // d := line[x]
            // newLine = append(newLine, d)
        // }
        // viewCells = append(viewCells, newLine)
    // }
    // t.viewCells = viewCells
// }


func (t *Terminal) Insert(ch rune) {
    t.HasCursorPane().insert(ch)
}

func (t *Terminal) Delete(length int) {
    t.HasCursorPane().delete(length)
}


func (t *Terminal) moveCursor(x, y int) {
    p := t.HasCursorPane()
    p.MoveCursor(x, y)
}


func (t *Terminal) SetCursor(x, y int) {
    p := t.HasCursorPane()
    t.moveCursor(x - p.CursorX, y - p.CursorY)
}

func (t *Terminal) MoveCursorRight() {
    p := t.HasCursorPane()
    c := getNextCell(p.getRealLine(p.CursorY), p.CursorX, 1)
    if c.Ch > 0 {
        x, y := c.Position()
        t.SetCursor(x, y)
    }
}

func (t *Terminal) MoveCursorUp() {
    t.moveCursor(0, -1)
}

func (t *Terminal) MoveCursorDown() {
    t.moveCursor(0, 1)
}

func (t *Terminal) MoveCursorLeft() {
    p := t.HasCursorPane()
    if p.xBegin == p.CursorX {
        return
    }
    c := getPrevCell(p.getRealLine(p.CursorY), p.CursorX, 1)
    if c.Ch > 0 {
        x, y := c.Position()
        t.SetCursor(x, y)
    }
}

func (t *Terminal) MoveCursorToLineEnd() {
    p := t.HasCursorPane()
    line := p.Line()
    lineWidth := cellsWidth(line)

    switch p.Mode {
        case ModeInsert: {
            LogFile("mode insert")
            t.SetCursor(lineWidth, p.CursorY)
        }
        case ModeNormal: {
            LogFile("mode normal")
            t.SetCursor(lineWidth - 1, p.CursorY)
        }
    }
}

func (t *Terminal) MoveCursorToLineBegin() {
    p := t.HasCursorPane()
    t.SetCursor(0, p.CursorY)
}

func (t *Terminal) MoveCursorToFirstLine() {
    t.SetCursor(0, 0)
}

func (t *Terminal) MoveCursorToLastLine() {
    p := t.HasCursorPane()
    y := min(len(p.cells) - 1, p.Height - 1)
    t.SetCursor(0, y)
}

func (t *Terminal) MoveCursorToLastCell() {
    t.MoveCursorToLastLine()
    t.MoveCursorToLineEnd()
}

// func (t *Terminal) realLine() []Cell {
    // return t.getRealLine(t.CursorY)
// }

// func (t *Terminal) getRealLine(y int) []Cell {
    // if len(t.cells) > 0 {
        // return t.cells[y + t.PageOffsetY]
    // }
    // return make([]Cell, 0)
// }

// func (t *Terminal) Line() []Cell {
    // return t.GetLineByY(t.CursorY)
// }

// func (t *Terminal) GetLineByY(y int) []Cell {
    // if len(t.viewCells) > 0 && y < len(t.viewCells) && y >= 0{
        // res := t.viewCells[y]
        // return res
    // }
    // return make([]Cell, 0)
// }

// func (t *Terminal) LineWidth() int{
    // return len(t.Line())
// }


func (t *Terminal) Close() {
    termbox.Close()
}

func (t *Terminal) PollEvent() termbox.Event{
    for {
        switch e := termbox.PollEvent(); e.Type {
            case termbox.EventKey:
                t.E.PreCh = t.E.Ch
                t.E.Ch = e.Ch
                return e
            case termbox.EventResize:
                t.Resize(e.Width, e.Height)
        }
    }
}

func (t *Terminal) Resize(w, h int) {

    t.Width = w
    t.Height = h
    t.CursorX = min(t.Width - 1, t.CursorX)
    t.CursorY = min(t.Height - 1, t.CursorY)
    t.Rendering()

}

func (t *Terminal) AppendCellFromString(str string) {
    t.HasCursorPane().AppendCellFromString(str)
}

func (t *Terminal) AppendCellFromStringWithColor(str string, fg, bg termbox.Attribute ) {
    t.HasCursorPane().AppendCellFromStringWithColor(str, fg, bg)
}

// func (t *Terminal) AppendCells(cells [][]Cell) {
    // for _, d := range cells {
        // t.cells = append(t.cells, d)
    // }
    // t.ResetViewCells()
    // t.ResetPageSize()
// }

// func (t *Terminal) SetLineCells(y int, c []Cell) {
    // t.setRealLine(y, c)
// }

// func (t *Terminal) setRealLine(y int, c []Cell) {
    // t.cells[y + t.PageOffsetY] = c
    // t.ResetViewCells()
    // t.ResetPageSize()
// }

func (t *Terminal) ListenKeyBorad(e termbox.Event) {
    switch e.Key {
        case termbox.KeyEsc: {
            os.Exit(0)
        }
    }

    if e.Ch <= 0 {
        return
    }

    switch e.Ch {
        case 'q': {
            os.Exit(0)
        }
    }
}

func (t *Terminal) ListenKeyBoardLikeVim(e termbox.Event) {
    p := t.HasCursorPane()
    LogFile(
        "listen", string(e.Ch),
    )
    switch e.Key {
        case termbox.KeyArrowLeft:  {
            t.MoveCursorLeft()
        }
        case termbox.KeyArrowRight:  {
            t.MoveCursorRight()
        }
        case termbox.KeyArrowDown:  {
            t.MoveCursorDown()
        }
        case termbox.KeyArrowUp:  {
            t.MoveCursorUp()
        }
        case termbox.KeyCtrlE: {
            t.MoveCursorToLineEnd()
        }
        case termbox.KeyCtrlA: {
            t.MoveCursorToLineBegin()
        }
        case termbox.KeySpace: {
            switch t.Mode {
                case ModeInsert: {
                    t.Insert(' ')
                }
                case ModeNormal: {
                    t.MoveCursorRight()
                }
            }
        }
    }
    switch p.Mode {
        case ModeNormal: {
            LogFile(
                "listen", "normal", string(e.Ch),
            )
            switch e.Key {
                case termbox.KeyEsc: {
                    os.Exit(0)
                }
            }

            if e.Ch > 0{
                switch e.Ch {
                    case 'l': {
                        t.MoveCursorRight()
                    }
                    case 'h': {
                        t.MoveCursorLeft()
                    }
                    case 'j': {
                        t.MoveCursorDown()
                    }
                    case 'i': {
                        t.SetMode(ModeInsert)
                    }
                    case 'k': {
                        t.MoveCursorUp()
                    }
                    case 'g': {
                        if t.E.PreCh == 'g' {
                            t.MoveCursorToFirstLine()
                        }
                    }
                    case 'G': {
                        t.MoveCursorToLastLine()
                    }
                    case '0': {
                        t.MoveCursorToLineBegin()
                    }
                    case '$': {
                        t.MoveCursorToLineEnd()
                    }
                }
            }
        }
        case ModeInsert: {
            switch e.Key {
                case termbox.KeyBackspace2: {
                    t.Delete(1)
                }
                case termbox.KeyEsc: {
                    t.SetMode(ModeNormal)
                }
            }

            if e.Ch > 0 {
                t.Insert(e.Ch)
            }
        }
    }
}

// func (t *Terminal) SetCursorLineColor(fg, bg termbox.Attribute) {
    // t.SetLineColor(t.CursorY, fg, bg)
// }

// func (t *Terminal) SetLineColor(y int, fg, bg termbox.Attribute) {
    // line := t.getRealLine(y)
    // newLine := make([]Cell, 0)
    // for _, d := range line {
        // newLine = append(newLine, Cell{Ch: d.Ch, Fg: fg, Bg: bg})
        // LogFile("setcolor", strconv.Itoa(y), string(d.Ch))
    // }
    // t.cells[y + t.PageOffsetY] = newLine
// }

// func (t *Terminal) GetCursorLineColor() (termbox.Attribute, termbox.Attribute) {
    // return t.GetLineColor(t.CursorY)
// }

// func (t *Terminal) GetLineColor(y int) (termbox.Attribute, termbox.Attribute) {
    // line := t.getRealLine(y)
    // if len(line) > 0{
        // return line[0].Fg, line[0].Bg
    // }
    // return ColorDefault, ColorDefault
// }

// func (t *Terminal) SubLineString(lineNum, b, e int) string {
    // line := t.GetLineByY(lineNum)

    // chs := make([]rune, 0)
    // for i := b; i <= e; i++ {
        // chs = append(chs, line[i].Ch)
    // }
    // return string(chs)
// }

func LogFile(str ...string) {
    file, _ := os.OpenFile("wsh.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    file.WriteString(strings.Join(str, " ") + "\n")
}

