package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gyyyy/sudoku-go/sudoku"
)

func main() {
	var (
		s string // 6######52|3#1#2####|###9###3#|###5#19##|#7##3##8#|##62#9###|#8###7###|####5#2#8|56######4
		v bool
	)
	flag.StringVar(&s, "s", "", "sudoku")
	flag.BoolVar(&v, "v", false, "verbose")
	flag.Parse()
	if s = strings.TrimSpace(s); s == "" {
		log.Fatalln("invalid arg [s]")
	}
	sudoku.Verbose = v
	sdk, err := sudoku.Create(s)
	if err != nil {
		log.Fatalln(err)
	}
	if fmt.Printf("题目：\n%s\n", sdk); sudoku.Resolve(sdk) {
		fmt.Printf("答案（完成）：\n%s\n", sdk.Print(false))
	} else {
		fmt.Printf("答案（未完成）：\n%s\n", sdk.Print(false))
		fmt.Println(sdk.Detail())
	}
}
