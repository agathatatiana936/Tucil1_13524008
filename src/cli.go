package main

import (
	"fmt"
	"os"
	"strings"
)

type ansi string

const (
	ansiReset ansi = "\x1b[0m"
	ansiBold  ansi = "\x1b[1m"
)

var regionFg = []ansi{
	"\x1b[31m", // red
	"\x1b[32m", // green
	"\x1b[33m", // yellow
	"\x1b[34m", // blue
	"\x1b[35m", // magenta
	"\x1b[36m", // cyan
	"\x1b[91m", // bright red
	"\x1b[92m", // bright green
	"\x1b[93m", // bright yellow
	"\x1b[94m", // bright blue
	"\x1b[95m", // bright magenta
	"\x1b[96m", // bright cyan
}

func colorsEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return true
}

func colorWrap(s string, codes ...ansi) string {
	if !colorsEnabled() || len(codes) == 0 {
		return s
	}
	var b strings.Builder
	for _, c := range codes {
		b.WriteString(string(c))
	}
	b.WriteString(s)
	b.WriteString(string(ansiReset))
	return b.String()
}

func cprintln(c ansi, s string) {
	fmt.Println(colorWrap(s, c))
}

func cprintf(c ansi, format string, a ...any) {
	out := fmt.Sprintf(format, a...)
	fmt.Print(colorWrap(out, c))
}


// ASCII 

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s) * n)
	for i := 0; i < n; i++ {
		b.WriteString(s)
	}
	return b.String()
}

func printBox(title string, lines []string) {
	maxLen := len(title)
	for _, ln := range lines {
		if len(ln) > maxLen {
			maxLen = len(ln)
		}
	}
	pad := 2
	width := maxLen + pad*2

	top := "┌" + repeat("─", width) + "┐"
	mid := "│" + repeat(" ", pad) + title + repeat(" ", width-pad-len(title)) + "│"

	fmt.Println(top)
	fmt.Println(mid)
	fmt.Println("├" + repeat("─", width) + "┤")
	for _, ln := range lines {
		row := "│" + repeat(" ", pad) + ln + repeat(" ", width-pad-len(ln)) + "│"
		fmt.Println(row)
	}
	fmt.Println("└" + repeat("─", width) + "┘")
}


// Board Grid Rendering

func queenLookup(solusi []int, n int) []int {
	q := make([]int, n)
	for i := 0; i < n; i++ {
		q[i] = -1
	}
	if solusi == nil {
		return q
	}
	for r := 0; r < n && r < len(solusi); r++ {
		q[r] = solusi[r]
	}
	return q
}

func regionColor(regionID int) ansi {
	if regionID < 0 {
		return "\x1b[90m" // gray
	}
	return regionFg[regionID%len(regionFg)]
}

func formatCell(ch rune, regionID int, isQueen bool) string {
	if isQueen {
		return colorWrap("#", ansiBold, "\x1b[97m")
	}
	return colorWrap(string(ch), regionColor(regionID))
}

func RenderBoardGrid(b *Queen, solusi []int) []string {
	N := b.n
	q := queenLookup(solusi, N)

	lines := make([]string, 0, N+2)

	cellW := 3
	top := "┌" + repeat("─", cellW)
	for i := 1; i < N; i++ {
		top += "┬" + repeat("─", cellW)
	}
	top += "┐"
	lines = append(lines, top)

	for r := 0; r < N; r++ {
		row := "│"
		for c := 0; c < N; c++ {
			ch := b.raw[r][c]
			reg := b.id[r][c]
			isQ := (q[r] == c)

			cell := formatCell(ch, reg, isQ)

			row += " " + cell + " "
			if c != N-1 {
				row += "│"
			}
		}
		row += "│"
		lines = append(lines, row)

		if r != N-1 {
			mid := "├" + repeat("─", cellW)
			for i := 1; i < N; i++ {
				mid += "┼" + repeat("─", cellW)
			}
			mid += "┤"
			lines = append(lines, mid)
		}
	}

	bot := "└" + repeat("─", cellW)
	for i := 1; i < N; i++ {
		bot += "┴" + repeat("─", cellW)
	}
	bot += "┘"
	lines = append(lines, bot)

	return lines
}

func PrintBoardGrid(title string, b *Queen, solusi []int) {
	if title != "" {
		fmt.Println(colorWrap(title, ansiBold))
	}
	for _, ln := range RenderBoardGrid(b, solusi) {
		fmt.Println(ln)
	}
}