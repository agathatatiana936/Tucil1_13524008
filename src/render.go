package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func repeatPlain(s string, n int) string {
	if n <= 0 {
		return ""
	}
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}


func build_grid_lines_plain(b *Queen, solusi []int) ([]string, error) {
	N := b.n
	if N <= 0 {
		return nil, fmt.Errorf("invalid board size: %d", N)
	}
	if solusi != nil && len(solusi) != N {
		return nil, fmt.Errorf("invalid solution length: got %d, expected %d", len(solusi), N)
	}

	q := make([]int, N)
	for i := 0; i < N; i++ {
		q[i] = -1
	}
	if solusi != nil {
		for r := 0; r < N; r++ {
			q[r] = solusi[r]
			if q[r] < 0 || q[r] >= N {
				return nil, fmt.Errorf("invalid queen column at row %d: %d", r+1, q[r])
			}
		}
	}

	cellW := 3 
	lines := make([]string, 0, N*2+1)

	top := "┌" + repeatPlain("─", cellW)
	for i := 1; i < N; i++ {
		top += "┬" + repeatPlain("─", cellW)
	}
	top += "┐"
	lines = append(lines, top)

	for r := 0; r < N; r++ {
		row := "│"
		for c := 0; c < N; c++ {
			ch := b.raw[r][c]
			if q[r] == c {
				ch = '#'
			}
			row += " " + string(ch) + " "
			if c != N-1 {
				row += "│"
			}
		}
		row += "│"
		lines = append(lines, row)

		if r != N-1 {
			mid := "├" + repeatPlain("─", cellW)
			for i := 1; i < N; i++ {
				mid += "┼" + repeatPlain("─", cellW)
			}
			mid += "┤"
			lines = append(lines, mid)
		}
	}

	bot := "└" + repeatPlain("─", cellW)
	for i := 1; i < N; i++ {
		bot += "┴" + repeatPlain("─", cellW)
	}
	bot += "┘"
	lines = append(lines, bot)

	return lines, nil
}

func write_lines_to_file(outPath string, lines []string, execMs int64, iterasi int64) error {
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %w", outPath, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < len(lines); i++ {
		_, _ = w.WriteString(lines[i] + "\n")
	}
	_, _ = w.WriteString(fmt.Sprintf("Waktu pencarian: %d ms\n", execMs))
	_, _ = w.WriteString(fmt.Sprintf("Banyak kasus yang ditinjau: %d kasus\n", iterasi))

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to write output file %q: %w", outPath, err)
	}
	return nil
}


// Output untuk live report

type liveReport struct {
	w        *bufio.Writer
	f        *os.File
	interval time.Duration
	last     time.Time
}

func live_writer(path string, interval time.Duration) (*liveReport, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace file %q: %w", path, err)
	}
	return &liveReport{
		w:        bufio.NewWriter(f),
		f:        f,
		interval: interval,
		last:     time.Time{},
	}, nil
}

func (tw *liveReport) Close() error {
	if tw == nil {
		return nil
	}
	if err := tw.w.Flush(); err != nil {
		_ = tw.f.Close()
		return err
	}
	return tw.f.Close()
}

func (tw *liveReport) writer_snapshot(b *Queen, solusi []int, iterasi int64) {
	if tw == nil {
		return
	}
	now := time.Now()
	if !tw.last.IsZero() && now.Sub(tw.last) < tw.interval {
		return
	}
	tw.last = now

	grid, err := build_grid_lines_plain(b, solusi)
	if err != nil {
		return
	}

	_, _ = tw.w.WriteString(fmt.Sprintf("ITERASI: %d\n", iterasi))
	_, _ = tw.w.WriteString(fmt.Sprintf("KANDIDAT (row->col): %v\n", solusi))
	for i := 0; i < len(grid); i++ {
		_, _ = tw.w.WriteString(grid[i] + "\n")
	}
	_, _ = tw.w.WriteString("----\n")
	_ = tw.w.Flush()
}