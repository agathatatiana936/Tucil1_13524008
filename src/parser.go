package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Data Structure
type Board struct {
	n int
	id [][]int
	idUnique int
	raw [][]rune
}

// Membaca file dan membuat sebuah Board dari data pada file
func parse_board_from_file(path string) (*Board, error) {
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

func parse_lines(lines []string) (*Board, error) {
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

	return &Board{
		n: N, 
		id: id, 
		idUnique: nextID,
		raw: raw,
	}, nil
}

