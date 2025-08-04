package sudoku

import (
	"log"
)

func deepCopy(sdk *Sudoku) *Sudoku {
	cpy := newSudoku()
	for i, gong := range sdk.gong {
		for j, cell := range gong.cell {
			cpy.gong[i].cell[j].isQuest = cell.isQuest
			cpy.gong[i].cell[j].value = cell.value
			cpy.gong[i].cell[j].canVal = map[int]any{}
			for k := range cell.canVal {
				cpy.gong[i].cell[j].canVal[k] = nil
			}
		}
		cpy.gong[i].remain = map[int]any{}
		for k := range gong.remain {
			cpy.gong[i].remain[k] = nil
		}
	}
	cpy.remain = map[int]int{}
	for k, v := range sdk.remain {
		cpy.remain[k] = v
	}
	return cpy
}

func check(sdk *Sudoku) bool {
	if !sdk.IsCompleted() {
		return false
	}
	for _, gong := range sdk.gong {
		if !gong.IsCompleted() {
			return false
		}
		mGong := map[int]any{}
		for _, cell := range gong.cell {
			if mGong[cell.value] = nil; gong.index[0] == gong.index[1] && cell.index[0] == cell.index[1] {
				var (
					mRow = map[int]any{}
					mCol = map[int]any{}
				)
				for _, cell := range sdk.Row(cell.x()) {
					mRow[cell.value] = nil
				}
				for _, cell := range sdk.Col(cell.y()) {
					mCol[cell.value] = nil
				}
				if len(mRow) != 9 || len(mCol) != 9 {
					return false
				}
			}
		}
		if len(mGong) != 9 {
			return false
		}
	}
	return true
}

func handleOnlyCan(sdk *Sudoku) error {
	var changed bool
	for _, gong := range sdk.gong {
		if gong.IsCompleted() {
			continue
		}
		m := map[int][]*Cell{}
		for _, cell := range gong.cell {
			if cell.hasValue() {
				continue
			}
			for can := range cell.canVal {
				m[can] = append(m[can], cell)
			}
		}
		for k, v := range m {
			n := len(v)
			if n < 1 || n > 3 {
				continue
			}
			if n == 1 {
				if Verbose {
					log.Printf("found only-can value: G(%d,%d) - C(%d,%d) => %d\n", gong.index[0], gong.index[1], v[0].index[0], v[0].index[1], k)
				}
				if err := v[0].Set(k); err != nil {
					return err
				}
				changed = true
				continue
			}

			var cells []*Cell
			if x := v[0].x(); v[1].x() == x {
				if n == 3 && v[2].x() != x {
					continue
				}
				for _, cell := range sdk.Row(x) {
					if cell.gong.is(v[0].gong) || cell.hasValue() {
						continue
					}
					if _, ok := cell.canVal[k]; !ok {
						continue
					}
					cells = append(cells, cell)
				}
			} else if y := v[0].y(); v[1].y() == y {
				if n == 3 && v[2].y() != y {
					continue
				}
				for _, cell := range sdk.Col(y) {
					if cell.gong.is(v[0].gong) || cell.hasValue() {
						continue
					}
					if _, ok := cell.canVal[k]; !ok {
						continue
					}
					cells = append(cells, cell)
				}
			} else {
				continue
			}
			if len(cells) == 0 {
				continue
			}
			if Verbose {
				log.Printf("found only-can-group value: G(%d,%d) - C(%d,%d) => %d\n", gong.index[0], gong.index[1], v[0].index[0], v[0].index[1], k)
			}
			for _, cell := range cells {
				if err := cell.Cannot(k); err != nil {
					return err
				}
				changed = true
			}
		}
	}
	if changed {
		return handleOnlyCan(sdk)
	}
	return nil
}

func bruteForce(sdk *Sudoku) bool {
	total := 1
	for _, v := range sdk.remain {
		total *= v
	}
	if Verbose {
		log.Printf("starting brute force with %d round\n", total)
	}
LOOP:
	for i := 0; i <= total; {
		cpy := deepCopy(sdk)
		for _, gong := range cpy.gong {
			if gong.IsCompleted() {
				continue
			}
			for _, cell := range gong.cell {
				if cell.hasValue() {
					continue
				}
				for n := range cell.canVal {
					if Verbose {
						log.Printf("[R-%d] brute forcing: G(%d,%d) - C(%d,%d) => %d\n", i+1, gong.index[0], gong.index[1], cell.index[0], cell.index[1], n)
					}
					if i++; cell.Set(n) == nil && handleOnlyCan(cpy) == nil {
						break
					}
					if cell.Cannot(n) == nil && sdk.Gong(gong.index[0], gong.index[1]).Cell(cell.index[0], cell.index[1]).Cannot(n) != nil {
						return false
					}
					continue LOOP
				}
			}
		}
		if check(cpy) {
			if *sdk = *cpy; Verbose {
				log.Printf("succeed in round %d\n", i)
			}
			return true
		}
	}
	if Verbose {
		log.Println("brute force failed")
	}
	return false
}

func Resolve(sdk *Sudoku) bool {
	return handleOnlyCan(sdk) == nil && (sdk.IsCompleted() || bruteForce(sdk))
}
