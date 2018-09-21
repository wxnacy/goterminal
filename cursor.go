package terminal

import (
)

type Cursor struct {
    X, Y int
    terminal *Terminal
}

func NewCursor(x, y int) *Cursor{
    c := &Cursor{
        X: x,
        Y: y,
    }
    return c
}

// func (c *Cursor) Move
