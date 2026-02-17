package main

import "math"

type solve struct {
	ada     bool
	hasil   []int 
	iterasi int64 
}

type trace_live func(b *Queen, solusi []int, iterasi int64)

func find_solusi(b *Queen, live trace_live) solve {
	N := b.n

	// solusi = array untuk simpan percobaan kemungkinan
	solusi := make([]int, N)
	for i := 0; i < N; i++ {
		solusi[i] = 0
	}

	var iterasi int64 = 0
	found := false

	// final = hasil copy dari array solusi kalau kemungkinan sudah valid
	final := make([]int, N)
	for i := 0; i < N; i++ {
		final[i] = -1
	}

	// cari semua kemungkinan solusi, stop kalau solusi udah follow semua rules
	var algo func(row int)
	algo = func(row int) {
		if found {
			return
		}

		if row == N {
			iterasi++

			if live != nil {
				live(b, solusi, iterasi)
			}

			if cek_solusi(b, solusi) {
				found = true
				for i := 0; i < N; i++ {
					final[i] = solusi[i]
				}
			}
			return
		}

		for col := 0; col < N; col++ {
			solusi[row] = col
			algo(row + 1)
			if found {
				return
			}
		}
	}

	algo(0)

	return solve{
		ada:     found,
		hasil:   final,
		iterasi: iterasi,
	}
}

func cek_solusi(b *Queen, solusi []int) bool {
	N := b.n

	usedColumn := make([]bool, N)
	for i  := 0; i  < N; i++ {
		c := solusi[i]
		if c < 0 || c >= N {
			return false
		}
		if usedColumn[c] {
			return false
		}
		usedColumn[c] = true
	}

	usedColor := make([]bool, b.idUnique)
	for i  := 0; i  < N; i++ {
		c := solusi[i]
		color := b.id[i][c]
		if color < 0 || color >= b.idUnique {
			return false
		}
		if usedColor[color] {
			return false
		}
		usedColor[color] = true
	}

	used := make([][]bool, N)
	for i := 0; i < N; i++ {
		used[i] = make([]bool, N)
	}
	for j := 0; j < N; j++ {
		used[j][solusi[j]] = true
	}

	// cek supaya solusi tidak bersebelahan dan tidak diagonal
	for k := 0; k < N; k++ {
		c := solusi[k]
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				if dr == 0 && dc == 0 {
					continue
				}
				nr := k + dr
				nc := c + dc
				if nr < 0 || nr >= N || nc < 0 || nc >= N {
					continue
				}
				if used[nr][nc] {
					return false
				}
			}
		}
	}

	return true
}

func total_configurations(N int) int64 {
	return int64(math.Pow(float64(N), float64(N)))
}
