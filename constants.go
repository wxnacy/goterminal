package terminal

import (
    "github.com/nsf/termbox-go"
)

type Mode uint8

const (
    ModeNormal Mode = iota
    ModeInsert
)

const (
    DefaultCursorX int = 0
    DefaultCursorY
    DefaultXBegin
    DefaultXEnd

    ColorDefault termbox.Attribute = termbox.ColorDefault
    ColorBlack termbox.Attribute = termbox.ColorBlack
    ColorRed termbox.Attribute = termbox.ColorRed
    ColorGreen termbox.Attribute = termbox.ColorGreen
    ColorYellow termbox.Attribute = termbox.ColorYellow
    ColorBlue termbox.Attribute = termbox.ColorBlue
    ColorMagenta termbox.Attribute = termbox.ColorMagenta
    ColorCyan termbox.Attribute = termbox.ColorCyan
    ColorWhite termbox.Attribute = termbox.ColorWhite
)
