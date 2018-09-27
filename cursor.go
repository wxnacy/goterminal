package terminal

import (
)

type Cursor struct {
    X, Y int
    Width, Height int
    xBegin, xEnd int
}

func newCursor(p *Pane) *Cursor{
    c := &Cursor{
        X: p.PositionX,
        Y: p.PositionY,
        Width: p.Width,
        Height: p.Height,
        xBegin: p.xBegin,
        xEnd: p.xEnd,
    }
    return c
}

