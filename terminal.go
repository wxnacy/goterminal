package terminal

import (
    "github.com/nsf/termbox-go"
    "strings"
    "os"
    "strconv"
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
}

func New() (*Terminal, error){
    err := termbox.Init()
    if err != nil {
        return nil, err
    }

    w, h := termbox.Size()

    return &Terminal{
        Width: w,
        Height: h,
        E: &Event{},
        cells: make([][]Cell, 0),
        viewCells: make([][]Cell, 0),
        Mode: ModeNormal,
        xBegin: DefaultXBegin,
        xEnd: DefaultXEnd,
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

func (t *Terminal) ResetPageSize() {
    t.PageHeight = len(t.cells)
}

func (t *Terminal) MoveCursorToLastCell() {
    t.MoveCursorToLastLine()
    t.MoveCursorToLineEnd()
}


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

func (t *Terminal) SetCells(cells [][]Cell) {
    t.cells = cells
    t.ResetViewCells()
    t.ResetPageSize()
}

// 渲染
func (t *Terminal) Rendering() {
    termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

    t.ResetViewCells()

    // chs := make([]rune, 0)

    for y, yd := range t.viewCells {
        for x, d := range yd {
            termbox.SetCell(x, y, d.Ch, d.Fg, d.Bg)
            // chs = append(chs, d.Ch)
        }
    }

    // LogFile("content", string(chs))

    // 矫正光标，使其最大出现在当前展示行的最后一个字付前
    switch t.Mode {
        case ModeInsert: {
            if t.CursorX >= t.LineWidth() {
                // minWidth := max(t.LineWidth(), 1)
                t.CursorX = t.LineWidth()
            }
        }
        case ModeNormal: {

            if t.CursorX >= t.LineWidth() {
                minWidth := max(t.LineWidth(), 1)
                t.CursorX = minWidth - 1
            }
        }
    }

    termbox.SetCursor(t.CursorX, t.CursorY)
    termbox.Flush()
}

func (t *Terminal) SetMode(m Mode) {
    t.Mode = m
}

// 重置显示的 cell 集合
func (t *Terminal) ResetViewCells() {
    viewCells := make([][]Cell, 0)

    // chs := make([]rune, 0)
    minLine := min(t.PageHeight, t.Height)
    for y := 0; y < minLine; y++ {
        newLine := make([]Cell, 0)
        index := min(y + t.PageOffsetY, t.PageHeight - 1)
        line := t.cells[index]
        for x := t.PageOffsetX; x < len(line); x++ {
            d := line[x]
            newLine = append(newLine, d)
            // chs = append(chs, d.Ch)
        }
        viewCells = append(viewCells, newLine)
    }
    // LogFile("reset ", string(chs))
    t.viewCells = viewCells
}


func (t *Terminal) SetCellBeforeCursor(ch rune, fg, bg termbox.Attribute) {

    cell := Cell{Ch: ch, Fg: fg, Bg: bg}
    line := t.realLine()

    newLine := make([]Cell, 0)

    // if t.CursorX >= len(line) {
        // newLine = append(newLine, line[:t.CursorX]...)
        // newLine = append(newLine, cell)
    // } else {
    // }

    newLine = append(newLine, line[:t.CursorX]...)
    newLine = append(newLine, cell)
    newLine = append(newLine, line[t.CursorX:]...)

    LogFile(strconv.Itoa(t.CursorX))

    t.cells[t.CursorY + t.PageOffsetY] = newLine
    // copy(t.cells[t.CursorY + t.PageOffsetY] , newLine)
    t.CursorX++

}

func (t *Terminal) RemoveCellBeforeCursor(length int) {
    if length <= 0 {
        return
    }
    newLine := make([]Cell, 0)
    line := t.realLine()

    index := len(line) - min(len(line), length)
    LogFile("line width ", strconv.Itoa(len(line)), "index", strconv.Itoa(index))
    index = max(index, t.xBegin)
    if index >= t.xBegin {
        newLine = append(line[0:index])
        // newLine = ArrayRemove(line, index)

        t.cells[t.CursorY + t.PageOffsetY] = newLine
        if len(newLine) != len(line) {
            t.CursorX--
        }
    }
    LogFile("cursor x ", strconv.Itoa(t.CursorX))
}

// 启动光标
func (t *Terminal) MoveCursor(x, y int) {
    line := t.GetLineByY(t.CursorY)
    minWidth := min(len(line), t.Width)
    minHeight := min(t.PageHeight, t.Height)

    cx := t.CursorX + x
    if cx >= minWidth {
        cx = minWidth - 1
        if t.Mode == ModeInsert {
            cx = minWidth
        }
    }

    cy := t.CursorY + y
    if cy >= minHeight {
        cy = minHeight - 1
        if t.Mode == ModeInsert {
            cy = minHeight
        }
    }

    LogFile(
        "move x", strconv.Itoa(x), "y", strconv.Itoa(y), "xbegin",
        strconv.Itoa(t.xBegin), "cx", strconv.Itoa(cx), "cursorx",
        strconv.Itoa(t.CursorX), "len(line)", strconv.Itoa(len(line)),
    )
    if cx >= t.xBegin {
        t.CursorX = cx
    }

    if cy >= 0 {
        t.CursorY = cy
    }

    if cx >= t.xBegin || cy >= 0 {
        LogFile("move")
        termbox.SetCursor(t.CursorX, t.CursorY)
        termbox.Flush()
    }
    if t.CursorY + 1 == t.Height || t.CursorY == 0{
        maxOffset := max(t.PageHeight, t.Height) - min(t.PageHeight, t.Height)
        t.PageOffsetY = min(t.PageOffsetY + y, maxOffset)
        if t.PageOffsetY < 0 {
            t.PageOffsetY = 0
        }
    }

}

func (t *Terminal) SetCursor(x, y int) {
    t.MoveCursor(x - t.CursorX, y - t.CursorY)
}

func (t *Terminal) MoveCursorToLineEnd() {
    switch t.Mode {
        case ModeInsert: {
            LogFile("mode insert")
            t.SetCursor(t.LineWidth(), t.CursorY)
        }
        case ModeNormal: {
            LogFile("mode normal")
            t.SetCursor(t.LineWidth() - 1, t.CursorY)
        }
    }
}

func (t *Terminal) MoveCursorToLineBegin() {
    t.SetCursor(0, t.CursorY)
}

func (t *Terminal) MoveCursorToFirstLine() {
    t.SetCursor(0, 0)
}

func (t *Terminal) MoveCursorToLastLine() {
    y := min(len(t.cells) - 1, t.Height - 1)
    t.SetCursor(0, y)
}

func (t *Terminal) realLine() []Cell {
    return t.getRealLine(t.CursorY)
}

func (t *Terminal) getRealLine(y int) []Cell {
    if len(t.cells) > 0 {
        return t.cells[y + t.PageOffsetY]
    }
    return make([]Cell, 0)
}

func (t *Terminal) Line() []Cell {
    return t.GetLineByY(t.CursorY)
}

func (t *Terminal) GetLineByY(y int) []Cell {
    if len(t.viewCells) > 0 && y < len(t.viewCells) && y >= 0{
        res := t.viewCells[y]
        // LogFile("getline", strconv.Itoa(y), "linewidth", strconv.Itoa(len(res)))
        return res
    }
    return make([]Cell, 0)
}

func (t *Terminal) LineWidth() int{
    return len(t.Line())
}


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
    t.AppendCellFromStringWithColor(str, ColorDefault, ColorDefault)
}

func (t *Terminal) AppendCellFromStringWithColor(str string, fg, bg termbox.Attribute ) {
    cells := stringToCellsWithColor(str, fg, bg)
    for _, d := range cells {
        t.cells = append(t.cells, d)
    }
    t.ResetViewCells()
    t.ResetPageSize()
}

func (t *Terminal) AppendCells(cells [][]Cell) {
    for _, d := range cells {
        t.cells = append(t.cells, d)
    }
    t.ResetViewCells()
    t.ResetPageSize()
}

func (t *Terminal) SetLineCells(y int, c []Cell) {
    t.setRealLine(y, c)
}

func (t *Terminal) setRealLine(y int, c []Cell) {
    t.cells[y + t.PageOffsetY] = c
    t.ResetViewCells()
    t.ResetPageSize()
}

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
    // e := t.PollEvent()
    LogFile(
        "listen", string(e.Ch),
    )
    switch e.Key {
        case termbox.KeyArrowLeft:  {
            t.MoveCursor(-1, 0)
        }
        case termbox.KeyArrowRight:  {
            t.MoveCursor(1, 0)
        }
        case termbox.KeyArrowDown:  {
            t.MoveCursor(0, 1)
        }
        case termbox.KeyArrowUp:  {
            t.MoveCursor(0, -1)
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
                    t.SetCellBeforeCursor(' ', ColorDefault, ColorDefault)
                }
                case ModeNormal: {
                    t.MoveCursor(1, 0)
                }
            }
        }
    }
    switch t.Mode {
        case ModeNormal: {
            LogFile(
                "listen", "normal", string(e.Ch),
            )

            switch e.Key {
                case termbox.KeyEsc: {
                    os.Exit(0)
                }
                case termbox.KeyArrowLeft:  {
                    t.MoveCursor(-1, 0)
                }
                case termbox.KeyArrowRight:  {
                    t.MoveCursor(1, 0)
                }
                case termbox.KeyArrowDown:  {
                    t.MoveCursor(0, 1)
                }
                case termbox.KeyArrowUp:  {
                    t.MoveCursor(0, -1)
                }
                case termbox.KeyCtrlE: {
                    t.MoveCursorToLineEnd()
                }
                case termbox.KeyCtrlA: {
                    t.MoveCursorToLineBegin()
                }
            }

            if e.Ch > 0{
                switch e.Ch {
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
                        t.SetMode(ModeInsert)
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
                        t.MoveCursorToLastLine()
                    }
                    case '0': {
                        t.SetCursor(0, t.CursorY)
                    }
                    case '$': {
                        t.SetCursor(t.LineWidth() - 1, t.CursorY)
                    }
                }
            }
        }
        case ModeInsert: {
            switch e.Key {
                case termbox.KeyBackspace2: {
                    t.RemoveCellBeforeCursor(1)
                }
            }
        }
    }
}

func (t *Terminal) SetCursorLineColor(fg, bg termbox.Attribute) {
    t.SetLineColor(t.CursorY, fg, bg)
}

func (t *Terminal) SetLineColor(y int, fg, bg termbox.Attribute) {
    line := t.getRealLine(y)
    newLine := make([]Cell, 0)
    for _, d := range line {
        newLine = append(newLine, Cell{Ch: d.Ch, Fg: fg, Bg: bg})
        LogFile("setcolor", strconv.Itoa(y), string(d.Ch))
    }
    t.cells[y + t.PageOffsetY] = newLine
}

func (t *Terminal) GetCursorLineColor() (termbox.Attribute, termbox.Attribute) {
    return t.GetLineColor(t.CursorY)
}

func (t *Terminal) GetLineColor(y int) (termbox.Attribute, termbox.Attribute) {
    line := t.getRealLine(y)
    if len(line) > 0{
        return line[0].Fg, line[0].Bg
    }
    return ColorDefault, ColorDefault
}

func (t *Terminal) SubLineString(lineNum, b, e int) string {
    line := t.GetLineByY(lineNum)

    chs := make([]rune, 0)
    for i := b; i <= e; i++ {
        chs = append(chs, line[i].Ch)
    }
    return string(chs)
}

func LogFile(str ...string) {
    file, _ := os.OpenFile("wsh.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    file.WriteString(strings.Join(str, " ") + "\n")
}

