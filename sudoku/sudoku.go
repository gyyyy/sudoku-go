package sudoku

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var Verbose = false

type Cell struct {
	gong    *Gong
	index   [2]int
	isQuest bool
	value   int
	canVal  map[int]any
}

func (c *Cell) x() int {
	return (c.gong.index[0]-1)*3 + c.index[0]
}

func (c *Cell) y() int {
	return (c.gong.index[1]-1)*3 + c.index[1]
}

func (c *Cell) is(cell *Cell) bool {
	return cell.x() == c.x() && cell.y() == c.y()
}

func (c *Cell) hasValue() bool {
	return c.value >= 1 && c.value <= 9
}

func (c *Cell) exclude() error {
	for _, cell := range c.gong.cell {
		if cell.hasValue() || c.is(cell) {
			continue
		}
		if err := cell.Cannot(c.value); err != nil {
			return err
		}
	}
	for _, cell := range c.gong.sdk.Row(c.x()) {
		if cell.hasValue() || c.is(cell) {
			continue
		}
		if err := cell.Cannot(c.value); err != nil {
			return err
		}
	}
	for _, cell := range c.gong.sdk.Col(c.y()) {
		if cell.hasValue() || c.is(cell) {
			continue
		}
		if err := cell.Cannot(c.value); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cell) Set(n int) error {
	if c.hasValue() || n < 1 || n > 9 {
		return nil
	}
	if Verbose {
		log.Printf("setted value: G(%d,%d) - C(%d,%d) => %d[%v]\n", c.gong.index[0], c.gong.index[1], c.index[0], c.index[1], n, c.canVal)
	}
	if _, ok := c.canVal[n]; !ok {
		return fmt.Errorf("invalid value: G(%d,%d) - C(%d,%d) => %d\n", c.gong.index[0], c.gong.index[1], c.index[0], c.index[1], n)
	}
	if _, ok := c.gong.remain[n]; !ok {
		return fmt.Errorf("invalid value in gong: G(%d,%d) - C(%d,%d) => %d\n", c.gong.index[0], c.gong.index[1], c.index[0], c.index[1], n)
	}
	if _, ok := c.gong.sdk.remain[n]; !ok || c.gong.sdk.remain[n] <= 0 {
		return fmt.Errorf("invalid value in global: G(%d,%d) - C(%d,%d) => %d\n", c.gong.index[0], c.gong.index[1], c.index[0], c.index[1], n)
	}
	c.value = n
	c.canVal = map[int]any{n: nil}
	delete(c.gong.remain, n)
	if c.gong.sdk.remain[n]--; c.gong.sdk.remain[n] == 0 {
		delete(c.gong.sdk.remain, n)
	}
	err := c.exclude()
	return err
}

func (c *Cell) Cannot(n int) error {
	if _, ok := c.canVal[n]; !ok {
		return nil
	}
	if len(c.canVal) == 1 {
		return fmt.Errorf("invalid cannot-value: G(%d,%d) - C(%d,%d) => %d\n", c.gong.index[0], c.gong.index[1], c.index[0], c.index[1], n)
	}
	if delete(c.canVal, n); len(c.canVal) == 1 && c.value == 0 {
		for v := range c.canVal {
			return c.Set(v)
		}
	}
	return nil
}

func (c *Cell) Value(onlyQuest, hasDetail bool) string {
	txt := strconv.Itoa(c.value)
	if c.value == 0 || (onlyQuest && !c.isQuest) {
		txt = "#"
	}
	if c.isQuest {
		return color.BlueString(txt)
	}
	if txt == "#" {
		if hasDetail {
			var (
				can = make([]string, len(c.canVal))
				i   int
			)
			for k := range c.canVal {
				can[i] = strconv.Itoa(k)
				i++
			}
			slices.Sort(can)
			txt += fmt.Sprintf("[%s]", strings.Join(can, "|"))
		}
		return color.RedString(txt)
	}
	return color.HiWhiteString(txt)
}

func (c *Cell) String() string {
	return fmt.Sprintf("宫(%d,%d) - 格(%d,%d): %s\n", c.gong.index[0], c.gong.index[1], c.index[0], c.index[1], c.Value(false, true))
}

type Gong struct {
	sdk    *Sudoku
	index  [2]int
	cell   [9]*Cell
	remain map[int]any
}

func (g *Gong) is(gong *Gong) bool {
	return gong.index[0] == g.index[0] && gong.index[1] == g.index[1]
}

func (g *Gong) CellByIndex(i int) *Cell {
	if i -= 1; i >= 0 && i <= 8 {
		return g.cell[i]
	}
	return nil
}

func (g *Gong) Cell(x, y int) *Cell {
	return g.CellByIndex((x-1)*3 + y)
}

func (g *Gong) Flat() [9]*Cell {
	return g.cell
}

func (g *Gong) IsCompleted() bool {
	return len(g.remain) == 0
}

func (g *Gong) String() string {
	str := fmt.Sprintf("宫(%d,%d)\n", g.index[0], g.index[1])
	for i, cell := range g.cell {
		if str += fmt.Sprintf("  格(%d,%d): %-24s", cell.index[0], cell.index[1], cell.Value(false, true)); i%3 == 2 {
			str += "\n"
		}
	}
	return str
}

func newGong(sdk *Sudoku, x, y int) *Gong {
	gong := &Gong{
		sdk:    sdk,
		index:  [2]int{x, y},
		cell:   [9]*Cell{},
		remain: map[int]any{1: nil, 2: nil, 3: nil, 4: nil, 5: nil, 6: nil, 7: nil, 8: nil, 9: nil},
	}
	for i := range 9 {
		gong.cell[i] = &Cell{
			gong:   gong,
			index:  [2]int{i/3 + 1, i%3 + 1},
			value:  0,
			canVal: map[int]any{1: nil, 2: nil, 3: nil, 4: nil, 5: nil, 6: nil, 7: nil, 8: nil, 9: nil},
		}
	}
	return gong
}

type Sudoku struct {
	gong   [9]*Gong
	remain map[int]int
}

func (s *Sudoku) GongByIndex(i int) *Gong {
	if i -= 1; i >= 0 && i <= 8 {
		return s.gong[i]
	}
	return nil
}

func (s *Sudoku) Gong(x, y int) *Gong {
	return s.GongByIndex((x-1)*3 + y)
}

func (s *Sudoku) GongRow(index int) [3]*Gong {
	gRow := [3]*Gong{}
	if index -= 1; index < 0 || index > 2 {
		return gRow
	}
	copy(gRow[:], s.gong[index*3:index*3+3])
	return gRow
}

func (s *Sudoku) Row(index int) [9]*Cell {
	cRow := [9]*Cell{}
	if index -= 1; index < 0 || index > 8 {
		return cRow
	}
	for i, gong := range s.gong[index/3*3 : index/3*3+3] {
		for j, cell := range gong.cell[index%3*3 : index%3*3+3] {
			cRow[i*3+j] = cell
		}
	}
	return cRow
}

func (s *Sudoku) GongCol(index int) [3]*Gong {
	gCol := [3]*Gong{}
	if index -= 1; index < 0 || index > 2 {
		return gCol
	}
	var i int
	for _, gong := range s.gong {
		if gong.index[1]-1 != index {
			continue
		}
		gCol[i] = gong
		if i++; i == 3 {
			return gCol
		}
	}
	return gCol
}

func (s *Sudoku) Col(index int) [9]*Cell {
	cCol := [9]*Cell{}
	if index -= 1; index < 0 || index > 8 {
		return cCol
	}
	var i int
	for _, gong := range s.gong {
		if gong.index[1]-1 != index/3 {
			continue
		}
		for _, cell := range gong.cell {
			if cell.index[1]-1 != index%3 {
				continue
			}
			cCol[i] = cell
			if i++; i == 9 {
				return cCol
			}
		}
	}
	return cCol
}

func (s *Sudoku) Flat() [9]*Gong {
	return s.gong
}

func (s *Sudoku) IsCompleted() bool {
	if len(s.remain) > 0 {
		return false
	}
	for _, gong := range s.gong {
		if !gong.IsCompleted() {
			return false
		}
	}
	return true
}

func (s *Sudoku) Detail() string {
	var str string
	for _, gong := range s.gong {
		str += gong.String()
	}
	return str
}

func (s *Sudoku) Print(onlyQuest bool) string {
	str := "+-------+-------+-------+\n"
	for gx := range 3 {
		for cx := range 3 {
			for gy := range 3 {
				gong := s.Gong(gx+1, gy+1)
				if gong == nil {
					continue
				}
				str += "| "
				for cy := range 3 {
					cell := gong.Cell(cx+1, cy+1)
					if cell == nil {
						continue
					}
					str += fmt.Sprintf("%s ", cell.Value(onlyQuest, false))
				}
			}
			str = str[:len(str)-1] + " |\n"
		}
		str += "+-------+-------+-------+\n"
	}
	return str[:len(str)-1]
}

func (s *Sudoku) String() string {
	return s.Print(true)
}

func newSudoku() *Sudoku {
	sdk := &Sudoku{
		gong:   [9]*Gong{},
		remain: map[int]int{1: 9, 2: 9, 3: 9, 4: 9, 5: 9, 6: 9, 7: 9, 8: 9, 9: 9},
	}
	for i := range 9 {
		sdk.gong[i] = newGong(sdk, i/3+1, i%3+1)
	}
	return sdk
}

func Create(s string) (*Sudoku, error) {
	var (
		sdk  = newSudoku()
		lSep = "\n"
	)
	if strings.Contains(s, "|") {
		lSep = "|"
	}
	rows := strings.SplitN(strings.TrimSpace(s), lSep, 9)
	if len(rows) != 9 {
		return nil, fmt.Errorf("invalid rows: %d", len(rows))
	}
	for i, row := range rows {
		cSep := " "
		if row = strings.TrimSpace(row); !strings.Contains(row, cSep) {
			cSep = ""
		}
		cols := strings.SplitN(row, cSep, 9)
		if len(cols) != 9 {
			return nil, fmt.Errorf("invalid cols in row[%d]: %d", i+1, len(cols))
		}
		for j, cell := range cols {
			if cell = strings.TrimSpace(cell); cell == "" || cell == "#" || cell == "." || cell == "0" {
				continue
			}
			n, err := strconv.Atoi(cell)
			if err != nil || n < 0 || n > 9 {
				return nil, fmt.Errorf("invalid number at (%d,%d)", i+1, j+1)
			}
			var (
				gx, gy = i/3 + 1, j/3 + 1
				cx, cy = i%3 + 1, j%3 + 1
				g      = sdk.Gong(gx, gy)
			)
			if g == nil {
				return nil, fmt.Errorf("invalid gong at (%d,%d)", gx, gy)
			}
			c := g.Cell(i%3+1, j%3+1)
			if c == nil {
				return nil, fmt.Errorf("invalid cell at (%d,%d) in gong at (%d,%d)", cx, cy, gx, gy)
			}
			if err = c.Set(n); err != nil {
				return nil, err
			}
			c.isQuest = true
		}
	}
	return sdk, nil
}
