//go:build !gui
// +build !gui

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	for {
		runOnce(in)

		fmt.Println()
		if !askYesNo(in, "Apakah ingin mencari solusi permainan Queen lainnya? (Ya/Tidak) ") {
			cprintln("\x1b[90m", "Bye") // gray
			return
		}
		fmt.Println()
	}
}

func runOnce(in *bufio.Reader) {
	printBox("Game Queen", []string{
		"1) Menerima input sebuah directory dari .txt yang memuat sebuah board",
		"2) Program akan mencari solusi peletakan queen sesuai dengan rules",
		"3) Hasil dan live report dapat di simpan",
	})

	// Input path
	inputPath := promptInputPath(in)

	// Parse
	b, err := parse_board_from_file(inputPath)
	if err != nil {
		cprintf("\x1b[31m", "ERR: %v\n", err)
		return
	}

	// Validasi input
	if err := validate_board(b); err != nil {
		cprintf("\x1b[31m", "INVALID BOARD: %v\n", err)
		return
	}

	// print board
	fmt.Println()
	PrintBoardGrid("BOARD INPUT", b, nil)
	fmt.Println()

	// tanya simpan
	wantTrace := askYesNo(in, "Apakah Anda ingin menyimpan proses pencarian (live report)? (Ya/Tidak) ")
	wantTXT := askYesNo(in, "Apakah Anda ingin menyimpan solusi sebagai .txt? (Ya/Tidak) ")
	wantPNG := askYesNo(in, "Apakah Anda ingin menyimpan solusi sebagai gambar PNG? (Ya/Tidak) ")

	needOutDir := wantTrace || wantTXT || wantPNG
	outDir := ""
	if needOutDir {
		outDir = promptOutputDirAndEnsure(in)
		cprintf("\x1b[90m", "Output folder: %s\n", outDir)
	}

	// live report
	var tw *liveReport
	var hook trace_live

	if wantTrace {
		tracePath := filepath.Join(outDir, "trace.txt")
		tw, err = live_writer(tracePath, 50*time.Millisecond)
		if err != nil {
			cprintf("\x1b[31m", "ERR: %v\n", err)
			return
		}
		defer tw.Close()

		cprintf("\x1b[90m", "Tracing aktif: %s (snapshot tiap 50ms)\n", tracePath)

		hook = func(b *Queen, solusi []int, iterasi int64) {
			tw.writer_snapshot(b, solusi, iterasi)
		}
	}

	// solve
	cprintln("\x1b[33m", "Solving...") // yellow
	start := time.Now()
	res := find_solusi(b, hook)
	execMs := time.Since(start).Milliseconds()

	if !res.ada {
		cprintln("\x1b[31m", "NO SOLUTION") // red
		fmt.Printf("Waktu pencarian: %d ms\n", execMs)
		fmt.Printf("Banyak kasus yang ditinjau: %d kasus\n", res.iterasi)
		return
	}

	// print hasil
	fmt.Println()
	cprintln("\x1b[32m", "SOLUTION FOUND ") // green
	PrintBoardGrid("BOARD RESULT", b, res.hasil)
	fmt.Println()

	fmt.Printf("Waktu pencarian: %d ms\n", execMs)
	fmt.Printf("Banyak kasus yang ditinjau: %d kasus\n", res.iterasi)

	// save
	if wantTXT {
		gridLines, err := build_grid_lines_plain(b, res.hasil)
		if err != nil {
			cprintf("\x1b[31m", "ERR: %v\n", err)
			return
		}

		outTXT := filepath.Join(outDir, "solution.txt")
		if err := write_lines_to_file(outTXT, gridLines, execMs, res.iterasi); err != nil {
			cprintf("\x1b[31m", "ERR: %v\n", err)
			return
		}
		cprintf("\x1b[36m", "Solusi disimpan ke %s\n", outTXT)
	}


	if wantPNG {
		crownPath := filepath.FromSlash("src/assets/crown.png")
		outPNG := filepath.Join(outDir, "solution.png")

		cellSize := 48
		gridThickness := 2
		outerBorder := 4

		if _, err := os.Stat(crownPath); err != nil {
			cprintf("\x1b[33m", "Warning: %s tidak ditemukan, PNG akan dibuat tanpa crown.\n", crownPath)
			crownPath = ""
		}

		if err := save_gambar(b, res.hasil, crownPath, outPNG, cellSize, gridThickness, outerBorder); err != nil {
			cprintf("\x1b[31m", "ERR: %v\n", err)
			return
		}
		cprintf("\x1b[36m", "Gambar disimpan ke %s\n", outPNG) // cyan
	}
}

func promptInputPath(in *bufio.Reader) string {
	for {
		fmt.Print("Masukkan path lengkap file input (.txt): ")
		s, _ := in.ReadString('\n')
		s = strings.TrimSpace(s)
		s = strings.Trim(s, `"'`)

		if s == "" {
			fmt.Println("Path kosong. Coba lagi.")
			continue
		}

		if strings.HasPrefix(s, "~") {
			if home, err := os.UserHomeDir(); err == nil {
				s = filepath.Join(home, strings.TrimPrefix(s, "~"))
			}
		}

		abs := s
		if !filepath.IsAbs(abs) {
			if a, err := filepath.Abs(abs); err == nil {
				abs = a
			}
		}
		abs = filepath.Clean(abs)

		info, err := os.Stat(abs)
		if err != nil {
			fmt.Println("File tidak ditemukan / tidak bisa diakses:", err)
			continue
		}
		if info.IsDir() {
			fmt.Println("Itu directory, bukan file. Masukkan path file .txt.")
			continue
		}

		return abs
	}
}

func promptOutputDirAndEnsure(in *bufio.Reader) string {
	for {
		fmt.Print("Masukkan directory/folder output untuk menyimpan file: ")
		s, _ := in.ReadString('\n')
		s = strings.TrimSpace(s)
		s = strings.Trim(s, `"'`)

		if s == "" {
			fmt.Println("Folder output kosong. Coba lagi.")
			continue
		}

		if strings.HasPrefix(s, "~") {
			if home, err := os.UserHomeDir(); err == nil {
				s = filepath.Join(home, strings.TrimPrefix(s, "~"))
			}
		}

		abs := s
		if !filepath.IsAbs(abs) {
			if a, err := filepath.Abs(abs); err == nil {
				abs = a
			}
		}
		abs = filepath.Clean(abs)

		if err := os.MkdirAll(abs, 0o755); err != nil {
			fmt.Println("Gagal membuat folder output:", err)
			continue
		}

		info, err := os.Stat(abs)
		if err != nil {
			fmt.Println("Folder output tidak bisa diakses:", err)
			continue
		}
		if !info.IsDir() {
			fmt.Println("Path output bukan folder (ternyata file). Masukkan folder yang benar.")
			continue
		}

		return abs
	}
}

func askYesNo(in *bufio.Reader, prompt string) bool {
	for {
		fmt.Print(prompt)
		s, _ := in.ReadString('\n')
		s = strings.ToLower(strings.TrimSpace(s))

		switch s {
		case "ya", "y", "yes":
			return true
		case "tidak", "t", "no", "n":
			return false
		default:
			fmt.Println("Jawaban tidak valid. Ketik Ya/Tidak.")
		}
	}
}