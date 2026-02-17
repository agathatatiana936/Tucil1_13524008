package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Data Structure for board game
type Queen struct {
	n int
	id [][]int
	idUnique int
	raw [][]rune
}

// Membaca file dan membuat sebuah Board dari data pada file
func parse_board_from_file(path string) (*Queen, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to Open FIle %q: %w", path, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	buff := make([]byte, 0, 64*1024)
	sc.Buffer(buff, 1024*1024)

	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err!=nil {
		return nil, fmt.Errorf("failed reading file %q: %w", path, err)
	}

	return parse_lines(lines)
}

// Membaca setiap line, hapus empty
func read_line(lines []string) []string {
	results := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		results = append(results, t)
	}
	return results
}

func parse_lines(lines []string) (*Queen, error) {
	lines = read_line(lines)
	N := len(lines)
	if N == 0 {
		return nil, fmt.Errorf("no lines found")
	}

	raw := make([][]rune, N)
	id := make([][]int, N)

	idMap := make(map[rune]int)
	nextID := 0

	for i := 0; i<N; i++ {
		temp := make([]rune, 0, N)
		runes := []rune(lines[i])
		for j := 0; j<len(runes); j++ {
			ch := runes[j]
			if unicode.IsSpace(ch){
				continue
			}
			temp = append(temp, ch)
		}

		if len(temp) != N {
			return nil, fmt.Errorf("different dimension: row %d has %d cells, expected %d (board must be NxN)", i+1, len(temp), N,)
		}

		raw[i] = make([]rune, N)
		id[i] = make([]int, N)

		for k := 0; k < N; k++ {
			a := temp[k]
			raw[i][k] = a

			b, ok := idMap[a]
			if !ok {
				idMap[a] = nextID
				b = nextID
				nextID++
			}
			id[i][k] = b
		}
	}

	return &Queen{
		n: N, 
		id: id, 
		idUnique: nextID,
		raw: raw,
	}, nil
}

func validate_board(b *Queen) error {
	if b == nil {
		return fmt.Errorf("board is nil")
	}
	if b.n <= 0 {
		return fmt.Errorf("invalid board size n=%d", b.n)
	}

	// cek dimensi
	if len(b.raw) != b.n {
		return fmt.Errorf("invalid raw rows: got %d, expected %d", len(b.raw), b.n)
	}
	if len(b.id) != b.n {
		return fmt.Errorf("invalid id rows: got %d, expected %d", len(b.id), b.n)
	}

	// cek jumlah column - panjang row
	for r := 0; r < b.n; r++ {
		if len(b.raw[r]) != b.n {
			return fmt.Errorf("invalid raw row %d length: got %d, expected %d", r+1, len(b.raw[r]), b.n)
		}
		if len(b.id[r]) != b.n {
			return fmt.Errorf("invalid id row %d length: got %d, expected %d", r+1, len(b.id[r]), b.n)
		}
	}

	if b.idUnique <= 0 {
		return fmt.Errorf("invalid region count: %d", b.idUnique)
	}
	if b.idUnique != b.n {
		return fmt.Errorf(
			"invalid board: number of regions (%d) must equal N (%d) because there must be exactly one queen per row/col/region",
			b.idUnique, b.n,
		)
	}

	for r := 0; r < b.n; r++ {
		for c := 0; c < b.n; c++ {
			v := b.id[r][c]
			if v < 0 || v >= b.idUnique {
				return fmt.Errorf("invalid region id at (%d,%d): got %d, expected 0..%d",
					r+1, c+1, v, b.idUnique-1)
			}
		}
	}

	// cek raw dan id, apakah konsisten?
	temp := make(map[rune]int)
	for i := 0; i < b.n; i++ {
		for j := 0; j < b.n; j++ {
			ch := b.raw[i][j]
			v := b.id[i][j]
			if old, ok := temp[ch]; ok {
				if old != v {
					return fmt.Errorf("inconsistent mapping for symbol %q: previously %d, now %d at (%d,%d)",
						ch, old, v, i+1, j+1)
				}
			} else {
				temp[ch] = v
			}
		}
	}

	return nil
}