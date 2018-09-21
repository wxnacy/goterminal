package terminal

// import (
    // "github.com/nsf/termbox-go"
    // "strings"
    // "os"
// )

// type Key termbox.Key
// type Attribute termbox.Attribute
// type Mode uint8

// const (
    // ModeNormal Mode = iota
    // ModeInsert
// )

// type Event struct {
    // PreCh rune
    // Ch rune
// }


// type Cell struct {
    // Ch rune
    // Fg termbox.Attribute    // 文字颜色
    // Bg termbox.Attribute    // 背景颜色
// }

// func NewCell(ch rune) *Cell {
    // return &Cell{
        // Ch: ch,
        // Fg: termbox.ColorDefault,
        // Bg: termbox.ColorDefault,
    // }
// }

// func stringToCells(s string) [][]*Cell {
    // cells := make([][]*Cell, 0)
    // lines := strings.Split(s, "\n")
    // for _, d := range lines {
        // ycells := make([]*Cell, 0)
        // chs := []rune(d)
        // for _, c := range chs {
            // cell := &Cell{
                // Ch: c,
                // Fg: termbox.ColorDefault,
                // Bg: termbox.ColorDefault,
            // }
            // ycells = append(ycells, cell)
        // }
        // cells = append(cells, ycells)
    // }
    // return cells
// }

// type Terminal struct {
    // Width, Height    int
	// CursorX, CursorY int
    // PageWidth, PageHeight int
	// PageOffsetX, PageOffsetY int
    // E *Event
    // Mode Mode
    // cells [][]*Cell
    // viewCells [][]*Cell

// }

// func New() (*Terminal, error){
    // err := termbox.Init()
    // if err != nil {
        // return nil, err
    // }

    // w, h := termbox.Size()

    // return &Terminal{
        // Width: w,
        // Height: h,
        // E: &Event{},
        // Mode: ModeNormal,
    // }, nil
// }

// func NewFromString(s string) (*Terminal, error){
    // err := termbox.Init()
    // if err != nil {
        // return nil, err
    // }

    // w, h := termbox.Size()

    // cells := stringToCells(s)
    // py := len(cells)

    // return &Terminal{
        // Width: w,
        // Height: h,
        // E: &Event{},
        // cells: cells,
        // PageHeight: py,
        // Mode: ModeNormal,
    // }, nil
// }

// func (t *Terminal) Run(onCh func(ch rune), onKey func(key termbox.Key)) {

    // for {
        // t.Rendering()
        // e := t.PollEvent()

        // if e.Key == termbox.KeyEsc {
            // os.Exit(0)
        // }

        // if e.Ch > 0 {
            // onCh(e.Ch)
        // } else {
            // onKey(e.Key)
        // }
    // }

// }

// func (t *Terminal) SetCell(x, y int, ch rune, fg, bg termbox.Attribute) {

    // t.cells[y][x].Ch = ch
    // t.cells[y][x].Fg = fg
    // t.cells[y][x].Bg = bg

    // termbox.SetCell(x, y, ch, fg, bg)
    // termbox.SetCursor(x, y)
    // termbox.Flush()
// }

// func (t *Terminal) SetCellOnCursor(ch rune, fg, bg termbox.Attribute) {
    // // cell := &Cell{Ch: ch, Fg: fg, Bg: bg}

    // // t.cells[t.CursorY + t.PageOffsetY] = append(t.realLine()[:t.CursorX], cell)
    // // line := make([]*Cell, 0)
    // // line = append(t.realLine()[:t.CursorX], cell)
    // // t.cells[t.CursorY + t.PageOffsetY] = line
    // // copy(t.cells[t.CursorY + t.PageOffsetY], )
    // t.CursorX++
// }

// func (t *Terminal) SetLineCells(y int, cells []*Cell) {
    // // t.cells[y] = make([]*Cell, 0)
    // line := t.cells[y + t.PageOffsetY]

    // for x, d := range cells {
        // // item := line[x]
        // if x + 1 >= len(line) {
            // line = append(line, d)
        // } else {
            // line[x] = d
        // }
    // }
// }

// // 渲染
// func (t *Terminal) Rendering() {
    // termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

    // t.ResetViewCells()

    // for y, yd := range t.viewCells {
        // for x, d := range yd {
            // termbox.SetCell(x, y, d.Ch, d.Fg, d.Bg)
        // }
    // }

    // // 矫正光标，使其最大出现在当前展示行的最后一个字付前
    // switch t.Mode {
        // case ModeInsert: {
            // if t.CursorX >= t.LineWidth() {
                // // minWidth := max(t.LineWidth(), 1)
                // t.CursorX = t.LineWidth()
            // }
        // }
        // case ModeNormal: {

            // if t.CursorX >= t.LineWidth() {
                // minWidth := max(t.LineWidth(), 1)
                // t.CursorX = minWidth - 1
            // }
        // }
    // }

    // termbox.SetCursor(t.CursorX, t.CursorY)
    // termbox.Flush()
// }

// func (t *Terminal) SetMode(m Mode) {
    // t.Mode = m
// }

// // 重置显示的 cell 集合
// func (t *Terminal) ResetViewCells() {
    // viewCells := make([][]*Cell, 0)

    // minLine := min(t.PageHeight, t.Height)
    // for y := 0; y < minLine; y++ {
        // newLine := make([]*Cell, 0)
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


// func (t *Terminal) SetCellBeforeCursor(ch rune, fg, bg termbox.Attribute) {
    // // viewCells := t.viewCells
    // // ycells := t.Line()
    // // x := t.CursorX
    // // y := t.CursorY

        // // for i := x; i < len(ycells); i++ {
            // // next := i + 1
            // // if next == len(ycells) {
                // // t.cells[y] = append(viewCells[y], ycells[i])
            // // } else {
                // // ycells[next] = ycells[i]
            // // }
        // // }
    // // if len(ycells) < x {
        // // viewCells[y] = append(viewCells[y], &Cell{Ch: ch, Fg: fg, Bg: bg})
    // // } else {
        // // viewCells[y][x] = &Cell{Ch: ch, Fg: fg, Bg: bg}
    // // }
    // // t.CursorX++

    // // line := t.realLine()
    // // x := t.CursorX

    // // line = append(line, NewCell(' '))

    // // for i := len(line) - 1; i > x; i-- {
        // // pre := max(i - 1, x + 1)
        // // line[i] = line[pre]
    // // }

    // // line[x] = &Cell{Ch: ch, Fg: fg, Bg: bg}

    // // cell := &Cell{Ch: ch, Fg: fg, Bg: bg}
    // // line := t.realLine()

    // // newLine := make([]*Cell, 0)
    // // newLine = append(newLine, line[:t.CursorX]...)
    // // newLine = append(newLine, cell)
    // // newLine = append(newLine, line...)
    // // t.cells[t.CursorY + t.PageOffsetY] = newLine
    // // t.cells[t.CursorY + t.PageOffsetY] = append(line[:t.CursorX], cell)
    // // t.cells[t.CursorY + t.PageOffsetY] = appelinend(line[:t.CursorX + 1], line[t.CursorX + 1:]...)
    // // line = append(line[:t.CursorX], cell)
    // // t.cells[t.CursorY + t.PageOffsetY] = newLine
    // t.CursorX++



// }

// // 启动光标
// func (t *Terminal) MoveCursor(x, y int) {
    // line := t.viewCells[t.CursorY]
    // minWidth := min(len(line), t.Width)
    // minHeight := min(t.PageHeight, t.Height)

    // cx := t.CursorX + x
    // if cx >= minWidth {
        // cx = minWidth - 1
        // if t.Mode == ModeInsert {
            // cx = minWidth
        // }
    // }

    // cy := t.CursorY + y
    // if cy >= minHeight {
        // cy = minHeight - 1
        // if t.Mode == ModeInsert {
            // cy = minHeight
        // }
    // }

    // if cx >= 0 {
        // t.CursorX = cx
    // }

    // if cy >= 0 {
        // t.CursorY = cy
    // }

    // if cx >=0 || cy >= 0 {
        // LogFile("move")
        // termbox.SetCursor(t.CursorX, t.CursorY)
        // termbox.Flush()
    // }
    // if t.CursorY + 1 == t.Height || t.CursorY == 0{
        // maxOffset := max(t.PageHeight, t.Height) - min(t.PageHeight, t.Height)
        // t.PageOffsetY = min(t.PageOffsetY + y, maxOffset)
        // if t.PageOffsetY < 0 {
            // t.PageOffsetY = 0
        // }
    // }

// }

// func (t *Terminal) SetCursor(x, y int) {
    // t.MoveCursor(x - t.CursorX, y - t.CursorY)
// }

// func (t *Terminal) realLine() []*Cell {
    // return t.cells[t.CursorY + t.PageOffsetY]
// }

// func (t *Terminal) Line() []*Cell {
    // return t.viewCells[t.CursorY]
// }

// func (t *Terminal) LineWidth() int{
    // return len(t.Line())
// }


// func (t *Terminal) Close() {
    // termbox.Close()
// }

// func (t *Terminal) PollEvent() termbox.Event{
    // for {
        // switch e := termbox.PollEvent(); e.Type {
            // case termbox.EventKey:
                // t.E.PreCh = t.E.Ch
                // t.E.Ch = e.Ch
                // return e
            // case termbox.EventResize:
                // t.Resize(e.Width, e.Height)
        // }
    // }
// }

// func (t *Terminal) Resize(w, h int) {

    // t.Width = w
    // t.Height = h
    // t.CursorX = min(t.Width - 1, t.CursorX)
    // t.CursorY = min(t.Height - 1, t.CursorY)
    // t.Rendering()

// }

// func LogFile(str ...string) {
    // file, _ := os.OpenFile("wsh.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    // file.WriteString(strings.Join(str, " ") + "\n")
// }

// func min(a, b int) int{
    // if a > b {
        // return b
    // }
    // return a
// }

// func max(a, b int) int{
    // if a > b {
        // return a
    // }
    // return b
// }
