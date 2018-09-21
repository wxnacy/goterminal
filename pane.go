package terminal

import (
)

type Pane struct {
    Width, Height int
	CursorX, CursorY int
    xBegin, xEnd int
    PageWidth, PageHeight int
	PageOffsetX, PageOffsetY int
    E *Event
    Mode Mode
    cells [][]Cell
    viewCells [][]Cell
}
