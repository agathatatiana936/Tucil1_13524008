package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func save_gambar(b *Queen, solusi []int, crownPath, outPath string, cellSize, tebalGrid, outerBorder int) error {
	N := b.n
	if N <= 0 {
		return fmt.Errorf("invalid board size: %d", N)
	}
	if len(solusi) != N {
		return fmt.Errorf("invalid solution length: got %d, expected %d", len(solusi), N)
	}

	w := outerBorder*2 + N*cellSize
	h := outerBorder*2 + N*cellSize

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	warna := generate_warna(b.idUnique)

	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			reg := b.id[r][c]
			if reg < 0 || reg >= b.idUnique {
				return fmt.Errorf("invalid region id at (%d,%d): %d", r+1, c+1, reg)
			}

			x0 := outerBorder + c*cellSize
			y0 := outerBorder + r*cellSize
			cellRect := image.Rect(x0, y0, x0+cellSize, y0+cellSize)

			draw.Draw(img, cellRect, &image.Uniform{warna[reg]}, image.Point{}, draw.Src)
		}
	}

	var crown image.Image
	if crownPath != "" {
		cr, err := loadPNG(crownPath)
		if err != nil {
			return fmt.Errorf("failed to load crown asset: %w", err)
		}
		crown = cr
	}

	if crown != nil {
		targetW := int(float64(cellSize) * 0.70)
		targetH := int(float64(cellSize) * 0.70)
		if targetW < 1 {
			targetW = 1
		}
		if targetH < 1 {
			targetH = 1
		}

		scaledCrown := resizeNearest(crown, targetW, targetH)

		for r := 0; r < N; r++ {
			c := solusi[r]
			if c < 0 || c >= N {
				return fmt.Errorf("invalid queen col at row %d: %d", r+1, c)
			}

			x0 := outerBorder + c*cellSize
			y0 := outerBorder + r*cellSize

			cx := x0 + (cellSize-targetW)/2
			cy := y0 + (cellSize-targetH)/2

			dstRect := image.Rect(cx, cy, cx+targetW, cy+targetH)
			draw.Draw(img, dstRect, scaledCrown, image.Point{}, draw.Over)
		}
	}

	gridColor := color.RGBA{0, 0, 0, 255}

	drawBorder(img, outerBorder, gridColor)

	if tebalGrid < 1 {
		tebalGrid = 1
	}

	for i := 1; i < N; i++ {
		x := outerBorder + i*cellSize
		drawVLine(img, x, outerBorder, outerBorder+N*cellSize, tebalGrid, gridColor)

		y := outerBorder + i*cellSize
		drawHLine(img, y, outerBorder, outerBorder+N*cellSize, tebalGrid, gridColor)
	}

	if err := savePNG(outPath, img); err != nil {
		return err
	}
	return nil
}


func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode png %q: %w", path, err)
	}
	return img, nil
}

func savePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("encode png %q: %w", path, err)
	}
	return nil
}


func drawBorder(img *image.RGBA, thickness int, col color.Color) {
	b := img.Bounds()
	// Top
	draw.Draw(img, image.Rect(b.Min.X, b.Min.Y, b.Max.X, b.Min.Y+thickness), &image.Uniform{col}, image.Point{}, draw.Src)
	// Bottom
	draw.Draw(img, image.Rect(b.Min.X, b.Max.Y-thickness, b.Max.X, b.Max.Y), &image.Uniform{col}, image.Point{}, draw.Src)
	// Left
	draw.Draw(img, image.Rect(b.Min.X, b.Min.Y, b.Min.X+thickness, b.Max.Y), &image.Uniform{col}, image.Point{}, draw.Src)
	// Right
	draw.Draw(img, image.Rect(b.Max.X-thickness, b.Min.Y, b.Max.X, b.Max.Y), &image.Uniform{col}, image.Point{}, draw.Src)
}

func drawVLine(img *image.RGBA, x, y0, y1, thickness int, col color.Color) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	draw.Draw(img, image.Rect(x-thickness/2, y0, x+(thickness+1)/2, y1), &image.Uniform{col}, image.Point{}, draw.Src)
}

func drawHLine(img *image.RGBA, y, x0, x1, thickness int, col color.Color) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	draw.Draw(img, image.Rect(x0, y-thickness/2, x1, y+(thickness+1)/2), &image.Uniform{col}, image.Point{}, draw.Src)
}


func generate_warna(count int) []color.RGBA {
	if count < 1 {
		return []color.RGBA{{200, 200, 200, 255}}
	}

	cols := make([]color.RGBA, count)
	for i := 0; i < count; i++ {
		h := float64(i) / float64(count) // 0..1
		cols[i] = hsv_to_rgb(h, 0.35, 0.95) // pastel-like
	}
	return cols
}

func hsv_to_rgb(h, s, v float64) color.RGBA {
	h = math.Mod(h, 1.0)
	if h < 0 {
		h += 1.0
	}

	c := v * s
	x := c * (1 - math.Abs(math.Mod(h*6, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case 0 <= h && h < 1.0/6.0:
		r, g, b = c, x, 0
	case 1.0/6.0 <= h && h < 2.0/6.0:
		r, g, b = x, c, 0
	case 2.0/6.0 <= h && h < 3.0/6.0:
		r, g, b = 0, c, x
	case 3.0/6.0 <= h && h < 4.0/6.0:
		r, g, b = 0, x, c
	case 4.0/6.0 <= h && h < 5.0/6.0:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	R := uint8(clamp01(r+m) * 255)
	G := uint8(clamp01(g+m) * 255)
	B := uint8(clamp01(b+m) * 255)
	return color.RGBA{R, G, B, 255}
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func resizeNearest(src image.Image, newW, newH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	sb := src.Bounds()

	sw := sb.Dx()
	sh := sb.Dy()
	if sw <= 0 || sh <= 0 {
		return dst
	}

	for y := 0; y < newH; y++ {
		sy := sb.Min.Y + (y*sh)/newH
		for x := 0; x < newW; x++ {
			sx := sb.Min.X + (x*sw)/newW
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}