package main

import "fmt"

func validate_board(b *Board) error {
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