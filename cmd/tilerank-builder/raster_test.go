// SPDX-License-Identifier: MIT

package main

import (
	"math"
	"testing"
)

func TestRaster_Paint(t *testing.T) {
	r := NewRaster(MakeTileKey(1, 1, 1), NewRaster(WorldTile, nil))
	r.Paint(MakeTileKey(2, 3, 3), 23)
	r.Paint(MakeTileKey(3, 6, 7), 42)
	wantPixels(t, r.pixels, [4][4]float32{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 23, 23},
		{0, 0, 65, 23},
	})
}

func TestRaster_Paint_SubPixel(t *testing.T) {
	tile := MakeTileKey(1, 0, 0)
	r := NewRaster(tile, NewRaster(WorldTile, nil))
	r.Paint(MakeTileKey(10, 256, 256), 100) // covers 1/4th of a pixel
	wantPixels(t, r.pixels, [4][4]float32{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 25, 0},
		{0, 0, 0, 0},
	})
}

func wantPixels(t *testing.T, got [256 * 256]float32, want [4][4]float32) {
	px := []int{0, 64, 128, 192}
	for j, vals := range want {
		for i, wantPix := range vals {
			x, y := px[i], px[j]
			gotPix := got[y<<8+x]
			if math.Abs(float64(gotPix-wantPix)) > 1e-4 {
				t.Errorf("pixel (%d, %d): got %f, want %f", x, y, gotPix, wantPix)
			}
		}
	}
}