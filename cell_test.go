package terminal

import (
    "testing"
)

func TestGetNextCell(t *testing.T) {
    var cells = []Cell{
        newCell('w'),
        newCell('测'),
        newCell(0),
        newCell('试'),
        newCell(0),
        newCell('一'),
        newCell(0),
        newCell('下'),
        newCell(0),
        newCell('n'),
        newCell('哈'),
        newCell(0),
        newCell('哈'),
        newCell(0),
        newCell('a'),
        newCell('b'),
    }


    c := getNextCell(cells, 0, 1)
    if c.Ch != '测' {
        t.Errorf("0 的下一个不是 %c", c.Ch)
    }

    c = getNextCell(cells, 0, 6)
    if c.Ch != '哈' {
        t.Errorf("0 的下一个不是 %c", c.Ch)
    }

    c = getNextCell(cells, 0, 8)
    if c.Ch != 'a' {
        t.Errorf("0 8 的下一个不是 %c", c.Ch)
    }
    c = getNextCell(cells, 0, 9)
    if c.Ch != 'b' {
        t.Errorf("0 9 的下一个不是 %c", c.Ch)
    }

    c = getNextCell(cells, 2, 1)
    if c.Ch != '试' {
        t.Errorf("0 的下一个不是 %c", c.Ch)
    }

    c = getNextCell(cells, 6, 1)
    if c.Ch != '下' {
        t.Errorf("0 的下一个不是 %c", c.Ch)
    }

}

func TestGetPrevCell(t *testing.T) {
    var cells = []Cell{
        newCell('w'),
        newCell('测'),
        newCell(0),
        newCell('试'),
        newCell(0),
        newCell('一'),
        newCell(0),
        newCell('下'),
        newCell(0),
        newCell('n'),
        newCell('哈'),
        newCell(0),
        newCell('哈'),
        newCell(0),
        newCell('a'),
        newCell('b'),
    }


    c := getPrevCell(cells, 0, 1)
    if c.Ch != 0 {
        t.Errorf("0 1 的下一个不是 %c", c.Ch)
    }

    c = getPrevCell(cells, 1, 1)
    if c.Ch != 'w' {
        t.Errorf("1 1 的下一个不是 %c", c.Ch)
    }
    c = getPrevCell(cells, 3, 1)
    if c.Ch != '测' {
        t.Errorf("3 1 的下一个不是 %c", c.Ch)
    }

    c = getPrevCell(cells, 10, 1)
    if c.Ch != 'n' {
        t.Errorf("10 1 的下一个不是 %c", c.Ch)
    }

    c = getPrevCell(cells, 14, 1)
    if c.Ch != '哈' {
        t.Errorf("14 1 的下一个不是 %c", c.Ch)
    }

    c = getPrevCell(cells, 15, 3)
    if c.Ch != '哈' {
        t.Errorf("15 3 的下一个不是 %c", c.Ch)
    }
}
